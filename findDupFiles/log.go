package main

import (
	"flag"
	"fmt"
)

var debugEnabled bool

func init() {
	flag.BoolVar(&debugEnabled, "debug", false, "Enable debug logging")
}

func debugLog(format string, args ...interface{}) {
	if debugEnabled {
		fmt.Printf("[DEBUG] "+format+" \n", args...)
	}
}
