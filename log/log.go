package log

import (
	"log"
	"os"
)

var debugEnabled bool

var (
	colorReset = "\033[0m"
	colorInfo  = "\033[34m" // blue
	colorError = "\033[31m" // red
	colorDebug = "\033[36m" // cyan
)

func Init(debug bool) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)
	debugEnabled = debug

	// Optionally auto-adapt colors
	if os.Getenv("COLOR_MODE") == "light" {
		colorInfo = "\033[35m"  // magenta for better contrast
		colorError = "\033[31m" // red
		colorDebug = "\033[33m" // yellow
	}
}

func Info(msg string, args ...any) {
	log.Printf(colorInfo+"[INFO] "+msg+colorReset, args...)
}

func Error(msg string, args ...any) {
	log.Printf(colorError+"[ERROR] "+msg+colorReset, args...)
}

func Debug(msg string, args ...any) {
	if debugEnabled {
		log.Printf(colorDebug+"[DEBUG] "+msg+colorReset, args...)
	}
}
