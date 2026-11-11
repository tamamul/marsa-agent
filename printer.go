package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

var pageRangeRe = regexp.MustCompile(`^[1-9][0-9]*(-[1-9][0-9]*)?(,[1-9][0-9]*(-[1-9][0-9]*)?)*$`)

// downloadFile ambil file job dari server, simpan ke folder temp.
// Endpoint file butuh token (header Authorization).
func downloadFile(url string, destDir string, jobID int, token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gagal download file, status %d", resp.StatusCode)
	}

	destPath := filepath.Join(destDir, fmt.Sprintf("job_%d.pdf", jobID))
	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}
	return destPath, nil
}

// printPDF kirim PDF ke printer lewat SumatraPDF (silent).
// pageRange kosong = semua halaman, contoh isi: "1-3,5".
func printPDF(filePath string, printerName string, copies int, pageRange string, baseDir string) error {
	sumatraPath := filepath.Join(baseDir, "SumatraPDF.exe")

	args := []string{"-print-to", printerName, "-silent"}

	if pageRange != "" {
		if len(pageRange) > 50 || !pageRangeRe.MatchString(pageRange) {
			return fmt.Errorf("format halaman tidak valid: %q", pageRange)
		}
		args = append(args, "-print-settings", pageRange)
	}
	args = append(args, filePath)

	for i := 0; i < copies; i++ {
		cmd := exec.Command(sumatraPath, args...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("gagal print copy ke-%d: %w", i+1, err)
		}
	}
	return nil
}

// cleanupFile hapus file job setelah selesai dicetak (jaga privasi + hemat disk)
func cleanupFile(path string) {
	_ = os.Remove(path)
}
