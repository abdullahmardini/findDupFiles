package dupes

import (
	"amardini/findDupFiles/log"
	"amardini/findDupFiles/utils"
	"fmt"
	"os"
)

type fileInfo struct {
	path string
}

func FindDuplicates(rootPath string, totalFiles int64) error {
	fileHashes := make(map[string][]fileInfo)

	err := utils.WalkFiles(rootPath, func(path string, info os.FileInfo) error {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening %q: %v\n", path, err)
			return nil
		}
		defer file.Close()

		hashString, err := utils.HashFile(path)
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
		log.Logger.Info("No duplicate files found.")
	}

	return nil
}
