package util

import (
	"os"
	"strings"
)

func IsCygwin() bool {
	if IsMSWindows() {
		shell := os.Getenv("SHELL")
		if len(shell) > 0 {
			return true
		}
	}
	return false
}

func IsMSWindows() bool {
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		return true
	}
	return false
}
