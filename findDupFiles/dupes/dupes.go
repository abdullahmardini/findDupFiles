package dupes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type fileInfo struct {
	path string
}

func FindDuplicates(rootPath string, totalFiles int64) error {
	fileHashes := make(map[string][]fileInfo)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading %q: %v\n", path, err)
			return err
		}

		if !info.Mode().IsRegular() || IsExcluded(path) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening %q: %v\n", path, err)
			return nil
		}
		defer file.Close()

		hashString, err := hashFile(path)
		if err != nil {
			fmt.Printf("Could not hash %q: %v\n", path, err)
		}

		fileHashes[hashString] = append(fileHashes[hashString], fileInfo{path: path})

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path %q: %v", rootPath, err)
	}

	foundDuplicates := false
	for hash, files := range fileHashes {
		if len(files) >= 2 {
			foundDuplicates = true
			fmt.Printf("\n--- Duplicate Hash found ---\n")
			fmt.Printf("Hash: %s\n", hash)
			fmt.Printf("Files:\n")
			for _, file := range files {
				fmt.Printf("  - %s\n", file.path)
			}
		}
	}

	if !foundDuplicates {
		fmt.Println("\nNo duplicate files found.")
	}

	return nil
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func IsExcluded(path string) bool {
	excludeDirs := []string{
		".git",
		"@snapshots",
		".zfs/snapshot",
		".cache",
	}

	for _, excl := range excludeDirs {
		if strings.Contains(path, excl) {
			return true
		}
	}
	return false
}
