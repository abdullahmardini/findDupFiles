package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type fileHash struct {
	path string
	hash string
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func countFiles(path string) (int64, error) {
	fmt.Printf("Scanning directory: %s/n", path)
	fmt.Println("Counting files...")
	var totalFiles int64
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			atomic.AddInt64(&totalFiles, 1)
		}
		return nil
	})

	return totalFiles, err
}

func main() {
	dirPtr := flag.String("dir", ".", "Directory to scan for duplicate files")
	flag.Parse()
	root := *dirPtr

	totalFiles, err := countFiles(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the directory: %v\n", err)
		return
	}
	fmt.Printf("Found %d files\n", totalFiles)

	fileChan := make(chan string, 100)
	resultChan := make(chan fileHash, 100)
	var wg sync.WaitGroup
	var processed int64

	numWorkers := 8
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				hash, err := hashFile(path)
				if err == nil {
					resultChan <- fileHash{path, hash}
				}
				atomic.AddInt64(&processed, 1)
			}
		}()
	}

	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				fileChan <- path
			}
			return nil
		})
		close(fileChan)
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			p := atomic.LoadInt64(&processed)
			fmt.Printf("Processed %d / %d files\r", p, totalFiles)
			if p >= totalFiles {
				break
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	hashes := make(map[string][]string)
	for result := range resultChan {
		hashes[result.hash] = append(hashes[result.hash], result.path)
	}

	// Print duplicates
	for hash, files := range hashes {
		if len(files) > 1 {
			fmt.Printf("Duplicate files (hash: %s):\n", hash)
			for _, file := range files {
				fmt.Printf("  %s\n", file)
			}
		}
	}
}
