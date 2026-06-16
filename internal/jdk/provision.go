package jdk

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// EnsureJava verifica se o JRE existe localmente; se não, baixa e retorna o caminho do executável
func EnsureJava() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("falha ao obter diretório home: %w", err)
	}

	baseDir := filepath.Join(home, ".hubsaude", "jdk")
	javaExe := "java"
	if runtime.GOOS == "windows" {
		javaExe = "java.exe"
	}

	// 1. Tenta encontrar o executável do Java já baixado
	existingPath := findExecutable(baseDir, javaExe)
	if existingPath != "" {
		return existingPath, nil
	}

	// 2. Se não encontrou, realiza o download e extração
	fmt.Println("⏳ JRE não encontrado. Iniciando download do Eclipse Temurin (Java 21)...")
	err = provision(baseDir)
	if err != nil {
		return "", fmt.Errorf("erro ao provisionar JRE: %w", err)
	}

	// 3. Busca novamente após extrair
	existingPath = findExecutable(baseDir, javaExe)
	if existingPath == "" {
		return "", fmt.Errorf("falha fatal: java não encontrado mesmo após o download")
	}

	fmt.Println("✔ JRE provisionado com sucesso!")
	return existingPath, nil
}

func getDownloadConfig() (string, string) {
	switch runtime.GOOS {
	case "windows":
		return "https://api.adoptium.net/v3/binary/latest/21/ga/windows/x64/jre/hotspot/normal/eclipse", ".zip"
	case "darwin":
		return "https://api.adoptium.net/v3/binary/latest/21/ga/mac/x64/jre/hotspot/normal/eclipse", ".tar.gz"
	default:
		return "https://api.adoptium.net/v3/binary/latest/21/ga/linux/x64/jre/hotspot/normal/eclipse", ".tar.gz"
	}
}

func provision(baseDir string) error {
	os.MkdirAll(baseDir, 0755)
	url, ext := getDownloadConfig()
	tempFile := filepath.Join(baseDir, "jre_download"+ext)

	// Download
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return err
	}

	// Extração
	if ext == ".zip" {
		err = extractZip(tempFile, baseDir)
	} else {
		err = extractTarGz(tempFile, baseDir)
	}

	os.Remove(tempFile) // Limpa o arquivo compactado após extrair
	return err
}

func findExecutable(root, exeName string) string {
	foundPath := ""
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.Name() == exeName {
			foundPath = path
			return filepath.SkipDir // Para a busca ao encontrar
		}
		return nil
	})
	return foundPath
}

// Funções utilitárias de extração
func extractZip(zipFile, dest string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue // Proteção contra ZipSlip
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}
	return nil
}

func extractTarGz(tarFile, dest string) error {
	f, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			io.Copy(outFile, tr)
			outFile.Close()
		}
	}
	return nil
}
