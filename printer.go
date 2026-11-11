package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"regexp"
)

var pageRangeRe = regexp.MustCompile(`^[1-9][0-9]*(-[1-9][0-9]*)?(,[1-9][0-9]*(-[1-9][0-9]*)?)*$`)
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


func printPDF(filePath string, printerName string, copies int, baseDir string) error {
	sumatraPath := filepath.Join(baseDir, "SumatraPDF.exe")

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
