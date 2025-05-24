package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	dirPtr := flag.String("dir", ".", "Directory to scan for duplicate files")
	flag.Parse()
	root := *dirPtr

	fmt.Printf("Scanning directory: %s\n", root)
	fmt.Println("Counting files...")
	totalFiles, err := countFiles(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting files: %v\n", err)
		return
	}
	fmt.Printf("Found %d files\n", totalFiles)

	findDuplicates(root, totalFiles)
}
