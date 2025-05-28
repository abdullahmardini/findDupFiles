package main

import (
	"amardini/findDupFiles/dupes"
	"amardini/findDupFiles/log"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
)

func countFiles(root string) (int64, error) {
	var total int64
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() || !dupes.IsExcluded(root) {
			atomic.AddInt64(&total, 1)
		}
		return nil
	})
	return total, err
}

func main() {
	dirPtr := flag.String("dir", ".", "Directory to scan for duplicate files")
	debugPtr := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()
	root := *dirPtr

	log.Init(*debugPtr)

	log.Logger.Info("Scanning directory:", "", root)
	fmt.Println("Counting files...")
	totalFiles, err := countFiles(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting files: %v\n", err)
		return
	}
	fmt.Printf("Found %d files\n", totalFiles)

	dupes.FindDuplicates(root, totalFiles)
}
