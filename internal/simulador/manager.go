package simulador

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/kyriosdata/assinatura/internal/jdk"
)

const (
	defaultPort        = 8090
	defaultHealthPath  = "/actuator/health"
	defaultCacheFolder = "simulador"
	defaultJarName     = "simulador.jar"
	defaultTimeout     = 30 * time.Second
	defaultReleaseAPI  = "https://api.github.com/repos/kyriosdata/runner/releases/latest"
)

type Options struct {
	Port       int
	SourceURL  string
	Checksum   string
	HealthPath string
	Timeout    time.Duration
	ReleaseAPI string
}

type ProcessInfo struct {
	PID       int       `json:"pid"`
	Port      int       `json:"port"`
	JarPath   string    `json:"jarPath"`
	StartedAt time.Time `json:"startedAt"`
}

type Status struct {
	Running bool
	Ready   bool
	Info    *ProcessInfo
	Message string
}

func DefaultOptions() Options {
	return Options{
		Port:       defaultPort,
		HealthPath: defaultHealthPath,
		Timeout:    defaultTimeout,
	}
}

func Start(ctx context.Context, opts Options) (Status, error) {
	opts = normalizeOptions(opts)

	if status, _ := CurrentStatus(ctx, opts); status.Running && status.Ready {
		return status, nil
	}

	if err := ensurePortAvailable(opts.Port); err != nil {
		return Status{}, err
	}

	jarPath, err := EnsureJar(ctx, opts)
	if err != nil {
		return Status{}, err
	}

	javaBin, err := jdk.EnsureJava()
	if err != nil {
		return Status{}, fmt.Errorf("nao foi possivel localizar ou provisionar Java 21: %w", err)
	}

	args := []string{
		fmt.Sprintf("-Dserver.port=%d", opts.Port),
		"-jar",
		jarPath,
	}
	cmd := exec.CommandContext(ctx, javaBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return Status{}, fmt.Errorf("falha ao iniciar simulador.jar: %w", err)
	}

	info := ProcessInfo{
		PID:       cmd.Process.Pid,
		Port:      opts.Port,
		JarPath:   jarPath,
		StartedAt: time.Now().UTC(),
	}
	if err := saveProcessInfo(info); err != nil {
		return Status{}, fmt.Errorf("simulador iniciado, mas falhou ao gravar registro: %w", err)
	}

	ready := waitUntilReady(ctx, opts)
	return Status{
		Running: true,
		Ready:   ready,
		Info:    &info,
		Message: readinessMessage(ready, opts.Port),
	}, nil
}

func Stop() (Status, error) {
	info, err := loadProcessInfo()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Status{Running: false, Ready: false, Message: "nenhum simulador registrado em execucao"}, nil
		}
		return Status{}, err
	}

	process, err := os.FindProcess(info.PID)
	if err != nil {
		_ = removeProcessInfo()
		return Status{Running: false, Ready: false, Info: &info, Message: "processo registrado nao foi encontrado"}, nil
	}

	if runtime.GOOS == "windows" {
		err = process.Kill()
	} else {
		err = process.Signal(os.Interrupt)
		if err == nil {
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		_ = removeProcessInfo()
		return Status{}, fmt.Errorf("falha ao encerrar processo %d: %w", info.PID, err)
	}

	_ = removeProcessInfo()
	return Status{Running: false, Ready: false, Info: &info, Message: "simulador encerrado com sucesso"}, nil
}

func CurrentStatus(ctx context.Context, opts Options) (Status, error) {
	opts = normalizeOptions(opts)
	info, err := loadProcessInfo()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Status{Running: false, Ready: false, Message: "simulador nao esta em execucao"}, nil
		}
		return Status{}, err
	}

	if info.Port != opts.Port && opts.Port != 0 {
		return Status{Running: false, Ready: false, Info: &info, Message: fmt.Sprintf("ha registro para a porta %d, nao para %d", info.Port, opts.Port)}, nil
	}

	ready := healthCheck(ctx, opts)
	if !ready {
		return Status{Running: false, Ready: false, Info: &info, Message: "processo registrado nao respondeu ao health check"}, nil
	}

	return Status{Running: true, Ready: true, Info: &info, Message: "simulador em execucao e pronto"}, nil
}

func EnsureJar(ctx context.Context, opts Options) (string, error) {
	opts = normalizeOptions(opts)

	localPath, err := localJarPath()
	if err != nil {
		return "", err
	}
	if fileExists(localPath) {
		if opts.Checksum != "" {
			if err := verifySHA256(localPath, opts.Checksum); err != nil {
				return "", err
			}
		}
		return localPath, nil
	}

	if fileExists(defaultJarName) {
		return defaultJarName, nil
	}

	if opts.SourceURL == "" {
		artifact, err := latestReleaseArtifact(ctx, opts.ReleaseAPI)
		if err != nil {
			return "", fmt.Errorf("simulador.jar nao encontrado localmente e falha ao consultar GitHub Releases: %w", err)
		}
		opts.SourceURL = artifact.DownloadURL
		if opts.Checksum == "" {
			opts.Checksum = artifact.Checksum
		}
	}

	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return "", err
	}

	if err := downloadFile(ctx, opts.SourceURL, localPath); err != nil {
		return "", err
	}
	if opts.Checksum != "" {
		if err := verifySHA256(localPath, opts.Checksum); err != nil {
			_ = os.Remove(localPath)
			return "", err
		}
	}

	return localPath, nil
}

func normalizeOptions(opts Options) Options {
	if opts.Port == 0 {
		opts.Port = defaultPort
	}
	if opts.HealthPath == "" {
		opts.HealthPath = defaultHealthPath
	}
	if opts.Timeout == 0 {
		opts.Timeout = defaultTimeout
	}
	if opts.ReleaseAPI == "" {
		opts.ReleaseAPI = defaultReleaseAPI
	}
	return opts
}

type releaseArtifact struct {
	DownloadURL string
	Checksum    string
}

type githubRelease struct {
	Assets []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func latestReleaseArtifact(ctx context.Context, apiURL string) (releaseArtifact, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return releaseArtifact{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return releaseArtifact{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return releaseArtifact{}, fmt.Errorf("GitHub Releases retornou HTTP %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return releaseArtifact{}, err
	}

	var artifact releaseArtifact
	for _, asset := range release.Assets {
		switch strings.ToLower(asset.Name) {
		case defaultJarName:
			artifact.DownloadURL = asset.BrowserDownloadURL
		case defaultJarName + ".sha256":
			checksum, err := readRemoteChecksum(ctx, asset.BrowserDownloadURL)
			if err != nil {
				return releaseArtifact{}, err
			}
			artifact.Checksum = checksum
		}
	}
	if artifact.DownloadURL == "" {
		return releaseArtifact{}, errors.New("asset simulador.jar nao encontrado no ultimo release")
	}
	return artifact, nil
}

func readRemoteChecksum(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checksum remoto retornou HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return "", errors.New("arquivo de checksum vazio")
	}
	return fields[0], nil
}

func ensurePortAvailable(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("porta %d indisponivel: %w", port, err)
	}
	return ln.Close()
}

func waitUntilReady(ctx context.Context, opts Options) bool {
	deadline := time.Now().Add(opts.Timeout)
	for time.Now().Before(deadline) {
		if healthCheck(ctx, opts) {
			return true
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(1 * time.Second):
		}
	}
	return false
}

func healthCheck(ctx context.Context, opts Options) bool {
	url := fmt.Sprintf("http://127.0.0.1:%d%s", opts.Port, opts.HealthPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func downloadFile(ctx context.Context, sourceURL, dest string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao baixar simulador.jar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download retornou HTTP %d", resp.StatusCode)
	}

	tmp := dest + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		_ = out.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return os.Rename(tmp, dest)
}

func verifySHA256(path, expected string) error {
	expected = strings.ToLower(strings.TrimSpace(expected))
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	actual := hex.EncodeToString(hash.Sum(nil))
	if actual != expected {
		return fmt.Errorf("checksum invalido para %s: esperado %s, obtido %s", path, expected, actual)
	}
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func hubSaudeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hubsaude"), nil
}

func localJarPath() (string, error) {
	base, err := hubSaudeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, defaultCacheFolder, defaultJarName), nil
}

func processInfoPath() (string, error) {
	base, err := hubSaudeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "simulador.json"), nil
}

func saveProcessInfo(info ProcessInfo) error {
	path, err := processInfoPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadProcessInfo() (ProcessInfo, error) {
	path, err := processInfoPath()
	if err != nil {
		return ProcessInfo{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ProcessInfo{}, err
	}
	var info ProcessInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return ProcessInfo{}, fmt.Errorf("registro do simulador invalido: %w", err)
	}
	return info, nil
}

func removeProcessInfo() error {
	path, err := processInfoPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func readinessMessage(ready bool, port int) string {
	if ready {
		return "simulador iniciado e pronto"
	}
	return "simulador iniciado, mas ainda nao respondeu ao health check na porta " + strconv.Itoa(port)
}
