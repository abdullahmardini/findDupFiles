package main

import (
	"amardini/findDupFiles/dupes"
	"amardini/findDupFiles/log"
	"amardini/findDupFiles/utils"
	"flag"
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
	log.Info("Found %d files\n", totalFiles)

	dupes.FindDuplicates(root, totalFiles)
}
