package dupes

import (
	"amardini/findDupFiles/log"
	"amardini/findDupFiles/utils"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
)

type fileInfo struct {
	path string
}

func FindDuplicates(rootPath string, totalFiles int64) error {
	fileHashes := make(map[string][]fileInfo)
	var processed int64

	err := utils.WalkFiles(rootPath, func(path string, info os.FileInfo) error {
		curr := atomic.AddInt64(&processed, 1)

		if curr%100 == 0 || curr == totalFiles {
			percent := float64(curr) / float64(totalFiles) * 100
			log.Info("Progress: %.0f%% - %d / %d", percent, curr, totalFiles)
		}
		file, err := os.Open(path)
		if err != nil {
			log.Error("Error opening %q: %v", path, err)
			return nil
		}
		defer file.Close()

		hashString, err := utils.HashFile(path)
		if err != nil {
			log.Error("Could not hash %q: %v", path, err)
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
			var paths []string

			for _, file := range files {
				paths = append(paths, file.path)
			}

			log.Info(
				fmt.Sprintf(
					"\nDuplicate hash found:\n  Hash: %s\n  Files (%d):\n    - %s\n",
					hash,
					len(paths),
					strings.Join(paths, "\n    - "),
				),
			)
		}
	}

	if !foundDuplicates {
		log.Info("No duplicate files found.")
	}

	return nil
}
