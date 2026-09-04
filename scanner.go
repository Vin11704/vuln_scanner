package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// // OLD: single rule, gi param rebuilt inside scan, os.Exit in scan
// func scan(root string, scanDir string, rule Rule, gi *Gitignore, numworkers int) []Finding {
// 	if !gi.CheckGitignore(root) {
// 		os.Exit(0)
// 	}
// 	gi, err := NewGitignore(root)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "failed to read .gitignore: %v\n", err)
// 		os.Exit(1)
// 	}
// 	...
// }

// scan walks scanDir, runs all rules against every file, and returns
// findings that survive gitignore suppression.
//
// Caller is responsible for CheckGitignore / NewGitignore / os.Exit
// — this function is a pure "do the work, return results" function.
func scan(scanDir string, rules []Rule, gi *Gitignore, numWorkers int) []Finding {
	pathsChan := make(chan string)
	findingsChan := make(chan Finding)
	var wg sync.WaitGroup

	// Launch workers
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(pathsChan, findingsChan, rules, gi, &wg)
	}

	// Walker goroutine — feeds file paths into pathsChan
	go func() {
		filepath.WalkDir(scanDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			pathsChan <- path
			return nil
		})
		close(pathsChan) // tells workers "no more files"
	}()

	// Closer goroutine — waits for all workers to finish, then closes
	// findingsChan so the collector loop below can exit.
	go func() {
		wg.Wait()
		close(findingsChan)
	}()

	// Collect results
	var results []Finding
	for f := range findingsChan {
		results = append(results, f)
	}

	return results
}

// worker reads file paths off the channel, runs every rule against each
// file, and forwards findings that aren't suppressed by gitignore.
func worker(paths <-chan string, findings chan<- Finding, rules []Rule, gi *Gitignore, wg *sync.WaitGroup) {
	defer wg.Done()

	for path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		for _, rule := range rules {
			hits := rule.Check(path, content)

			for _, f := range hits {
				if gi.Match(path) {
					continue // suppress *reporting* only
				}
				findings <- f
			}
		}
	}
}

// // OLD worker: single rule, chan []Finding, all-or-nothing gitignore per file
// func worker(paths <-chan string, results chan<- []Finding, rule Rule, gi *Gitignore, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for path := range paths {
// 		content, err := os.ReadFile(path)
// 		if err != nil {
// 			continue
// 		}
// 		findings := rule.Check(path, content)
// 		if len(findings) == 0 {
// 			continue
// 		}
// 		if gi.Match(path) {
// 			continue
// 		}
// 		results <- findings
// 	}
// }
