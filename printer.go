package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// downloadFile ambil file job dari server, simpan ke folder temp
func downloadFile(url string, destDir string, jobID int) (string, error) {
	resp, err := http.Get(url)
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

// printPDF kirim file PDF ke printer tertentu secara silent (tanpa dialog)
// pakai SumatraPDF.exe (portable, gratis, bagus untuk automation cetak).
// Taruh SumatraPDF.exe satu folder dengan agent.exe, atau set path lengkap.
func printPDF(filePath string, printerName string, copies int) error {
	sumatraPath := "SumatraPDF.exe" // bisa diganti path absolut kalau perlu

	for i := 0; i < copies; i++ {
		cmd := exec.Command(sumatraPath,
			"-print-to", printerName,
			"-silent",
			filePath,
		)
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
