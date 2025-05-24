package main

import (
	"os"
	"path/filepath"
	"sync/atomic"
)

func countFiles(root string) (int64, error) {
	var total int64
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() || !isExcluded(path) {
			atomic.AddInt64(&total, 1)
		}
		return nil
	})
	return total, err
}
