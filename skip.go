package main

import (
	"path/filepath"
	"strings"
)

// skipExts contains binary/non-text file extensions that should be
// excluded from scanning to avoid false positives.
var skipExts = map[string]bool{
	// Compiled binaries
	".exe":   true,
	".dll":   true,
	".so":    true,
	".bin":   true,
	".o":     true,
	".a":     true,
	".dylib": true,

	// JVM / Python bytecode
	".pyc":   true,
	".class": true,
	".jar":   true,

	// Archives
	".zip": true,
	".tar": true,
	".gz":  true,

	// Images / media
	".png": true,
	".jpg": true,
	".gif": true,
	".pdf": true,
	".ico": true,

	// Fonts
	".woff":  true,
	".woff2": true,
	".ttf":   true,
}

// shouldSkipFile returns true if the file has a binary extension
func shouldSkipFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return skipExts[ext]
}
