package utils

import (
	"amardini/findDupFiles/log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

func WalkFiles(root string, onFile func(path string, info os.FileInfo) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || IsExcluded(path) {
			return nil
		}
		return onFile(path, info)
	})
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
			log.Debug("Excluded files...", "path", path)
			return true
		}
	}
	return false
}

func CountFiles(root string) (int64, error) {
	var total int64
	err := WalkFiles(root, func(path string, info os.FileInfo) error {
		log.Debug("Counting file:", "path", path)
		atomic.AddInt64(&total, 1)
		return nil
	})
	return total, err
}
