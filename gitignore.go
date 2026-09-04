package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Gitignore struct {
	patterns []string //entries in gitignore file
	root     string
}

// // OLD: was a method — required a *Gitignore to exist before checking
// func (g *Gitignore) CheckGitignore(dir string) bool {
// 	_, err := os.Stat(filepath.Join(dir, ".gitignore"))
// 	return err == nil
// }

// CheckGitignore checks if a .gitignore file exists in the given directory.
// Standalone function so it can be called before NewGitignore.
func CheckGitignore(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".gitignore"))
	return err == nil
}

func NewGitignore(dir string) (*Gitignore, error) {
	// content, err := os.ReadFile(dir) // OLD: dir is a directory, not the .gitignore file
	content, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var patterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == " " || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	return &Gitignore{
		patterns: patterns,
		root:     dir,
	}, nil
}

// check if the file(s) containing secrets is ignored in gitignore
func (g *Gitignore) Match(path string) bool {
	rel, err := filepath.Rel(g.root, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)

	for _, pat := range g.patterns {
		pat = filepath.ToSlash(pat)
		// check if path contains directory
		if strings.Contains(pat, "/") {
			dirName := strings.TrimSuffix(pat, "/")
			segments := strings.Split(rel, string(filepath.Separator))

			for _, i := range segments {
				if i == dirName {
					return true
				}
			}

		} else {
			matched, _ := filepath.Match(pat, rel)
			if matched {
				return true
			}
		}
	}
	return false
}
