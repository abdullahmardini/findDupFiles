package main

import (
	"crypto/sha256"
	"encoding/hex"
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

type fileInfo struct {
	path string
}

func findDuplicates(rootPath string, totalFiles int64) error {
	fileHashes := make(map[string][]fileInfo)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading %q: %v\n", path, err)
			return err
		}

		if !info.Mode().IsRegular() || isExcluded(path) {
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

func findDuplicatesConcurrent(root string, totalFiles int64) {
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
			if err == nil && !info.IsDir() && !isExcluded(path) {
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

	fmt.Println("\n\nDuplicate files found:")
	for hash, files := range hashes {
		if len(files) > 1 {
			fmt.Printf("Hash: %s\n", hash)
			for _, file := range files {
				fmt.Printf("  %s\n", file)
			}
		}
	}
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
