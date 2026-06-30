package simulador

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureJarDownloadsAndValidatesChecksum(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)

	payload := []byte("jar fake para teste")
	sum := sha256.Sum256(payload)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer server.Close()

	path, err := EnsureJar(context.Background(), Options{
		SourceURL: server.URL,
		Checksum:  hex.EncodeToString(sum[:]),
	})
	if err != nil {
		t.Fatalf("EnsureJar retornou erro: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("falha ao ler jar baixado: %v", err)
	}
	if string(data) != string(payload) {
		t.Fatalf("conteudo baixado inesperado: %q", string(data))
	}
}

func TestEnsureJarRejectsInvalidChecksum(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("jar fake"))
	}))
	defer server.Close()

	_, err := EnsureJar(context.Background(), Options{
		SourceURL: server.URL,
		Checksum:  "0000000000000000000000000000000000000000000000000000000000000000",
	})
	if err == nil {
		t.Fatal("esperava erro de checksum invalido")
	}
}

func TestCurrentStatusWithoutRegistry(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)

	status, err := CurrentStatus(context.Background(), DefaultOptions())
	if err != nil {
		t.Fatalf("CurrentStatus retornou erro: %v", err)
	}
	if status.Running || status.Ready {
		t.Fatalf("status inesperado: %+v", status)
	}
}

func TestEnsureJarUsesCachedFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)

	cached := filepath.Join(home, ".hubsaude", defaultCacheFolder, defaultJarName)
	if err := os.MkdirAll(filepath.Dir(cached), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cached, []byte("cached"), 0644); err != nil {
		t.Fatal(err)
	}

	path, err := EnsureJar(context.Background(), DefaultOptions())
	if err != nil {
		t.Fatalf("EnsureJar retornou erro: %v", err)
	}
	if path != cached {
		t.Fatalf("esperava caminho em cache %q, obteve %q", cached, path)
	}
}

func TestEnsureJarResolvesLatestReleaseWhenSourceIsEmpty(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)

	payload := []byte("jar vindo do release")
	sum := sha256.Sum256(payload)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/latest":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"assets":[{"name":"simulador.jar","browser_download_url":"` + server.URL + `/simulador.jar"},{"name":"simulador.jar.sha256","browser_download_url":"` + server.URL + `/simulador.jar.sha256"}]}`))
		case "/simulador.jar":
			_, _ = w.Write(payload)
		case "/simulador.jar.sha256":
			_, _ = w.Write([]byte(hex.EncodeToString(sum[:]) + "  simulador.jar\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	path, err := EnsureJar(context.Background(), Options{ReleaseAPI: server.URL + "/latest"})
	if err != nil {
		t.Fatalf("EnsureJar retornou erro: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("falha ao ler jar: %v", err)
	}
	if string(data) != string(payload) {
		t.Fatalf("conteudo inesperado: %q", string(data))
	}
}
