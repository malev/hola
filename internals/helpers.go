package internals

import (
	"fmt"
	"log/slog"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		slog.Debug(fmt.Sprintf("Error checking file: %v\n", err))
		return false
	}

	if info.IsDir() {
		slog.Debug(fmt.Sprintf("%s is a directory\n", filename))
		return false
	}

	return true
}
