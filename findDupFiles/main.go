package main

import (
	"amardini/findDupFiles/dupes"
	"amardini/findDupFiles/log"
	"amardini/findDupFiles/utils"
	"flag"
	"fmt"
	"os"
)

func main() {
	dirPtr := flag.String("dir", ".", "Directory to scan for duplicate files")
	debugPtr := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()
	root := *dirPtr

	log.Init(*debugPtr)

	log.Logger.Info("Scanning directory:", "", root)
	fmt.Println("Counting files...")
	totalFiles, err := utils.CountFiles(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting files: %v\n", err)
		return
	}
	fmt.Printf("Found %d files\n", totalFiles)

	dupes.FindDuplicates(root, totalFiles)
}
