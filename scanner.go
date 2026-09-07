package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// this function is a pure "do the work, return results" function.
func scan(scanDir string, rules []Rule, gi *Gitignore, numWorkers int) []Finding {
	pathsChan := make(chan string)
	findingsChan := make(chan Finding)
	var wg sync.WaitGroup

	// Launch workers
	wg.Add(numWorkers)
	// for i := 0; i < numWorkers; i++ {
	for range numWorkers {
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
			if shouldSkipFile(path) {
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
				// blocked as this prevents SAFE statement in main.go from executing
				// if gi.Match(path) {
				// 	continue // suppress *reporting* only
				// }
				findings <- f
			}
		}
	}
}
