package main

import (
	"strings"
)

func isExcluded(path string) bool {
	excludeDirs := []string{
		".git",
		"@snapshots",
		".zfs/snapshot",
		".cache",
	}

	for _, excl := range excludeDirs {
		if strings.Contains(path, excl) {
			return true
		}
	}
	return false
}
