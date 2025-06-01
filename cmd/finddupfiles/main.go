package main

import (
	"flag"

	"github.com/abdullahmardini/findDupFiles/dupes"
	"github.com/abdullahmardini/findDupFiles/log"
	"github.com/abdullahmardini/findDupFiles/utils"
)

func main() {
	dirPtr := flag.String("dir", ".", "Directory to scan for duplicate files")
	debugPtr := flag.Bool("debug", false, "Enable debug logging")

	flag.Parse()
	root := *dirPtr
	log.Init(*debugPtr)

	log.Info("Scanning directory: %s", root)
	log.Info("Counting files...")
	totalFiles, err := utils.CountFiles(root)
	if err != nil {
		log.Error("Error counting files: %v", err)
	}
	log.Info("Found %d files", totalFiles)

	dupes.FindDuplicates(root, totalFiles)
}
