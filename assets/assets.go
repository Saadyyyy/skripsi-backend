package assets

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func PilihGambarAcak() (string, error) {
	gambarList := []string{"foto1.png", "foto2.png", "foto3.png", "foto4.png"}
	rand.Seed(time.Now().UnixNano())
	gambarTerpilih := gambarList[rand.Intn(len(gambarList))]

	// Path ke file gambar di folder assets
	path := filepath.Join("assets", gambarTerpilih) // Adjust path to match server structure

	fmt.Printf("Checking path: %s\n", path) // Debug path

	// Pastikan file gambar dapat diakses
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("file gambar tidak ditemukan: %s", path)
	}

	return "/" + gambarTerpilih, nil // Return path with leading "/"
}
