package dirwalker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type ProgressFn func(path string, size int64, isDir bool)

func ScanDirectory(root string, maxDepth int, excludes []string, progress ProgressFn, workers int, reportFiles bool) (int64, error) {
	rootInfo, err := os.Stat(root)
	if err != nil {
		return 0, fmt.Errorf("stat %s: %w", root, err)
	}
	if !rootInfo.IsDir() {
		return rootInfo.Size(), nil
	}

	excludeSet := buildExcludeSet(excludes)

	if workers <= 0 {
		workers = min(runtime.NumCPU()*2, 8)
	}

	type job struct {
		path  string
		depth int
	}

	work := make(chan job, 100000)
	var pending sync.WaitGroup
	pending.Add(1)
	work <- job{path: root, depth: 0}

	go func() {
		pending.Wait()
		close(work)
	}()

	var mu sync.Mutex
	var total int64
	var workerWG sync.WaitGroup

	for range workers {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()
			for j := range work {
				children, size := processDir(j.path, j.depth, maxDepth, excludeSet, reportFiles, progress)

				if progress != nil {
					progress(j.path, size, true)
				}

				mu.Lock()
				total += size
				mu.Unlock()

				for _, child := range children {
					pending.Add(1)
					work <- job{path: child, depth: j.depth + 1}
				}

				pending.Done()
			}
		}()
	}

	workerWG.Wait()
	return total, nil
}

func processDir(path string, depth, maxDepth int, excludes map[string]bool, reportFiles bool, progress ProgressFn) (children []string, size int64) {
	if maxDepth > 0 && depth > maxDepth {
		addSize := func(p string, s int64) {
			size += s
			if progress != nil && reportFiles {
				progress(p, s, false)
			}
		}
		filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			info, err := d.Info()
			if err == nil && info.Mode().IsRegular() {
				addSize(p, info.Size())
			}
			return nil
		})
		return nil, size
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, 0
	}

	for _, entry := range entries {
		p := filepath.Join(path, entry.Name())
		if excludes[p] || excludes[entry.Name()] {
			continue
		}

		if entry.IsDir() {
			children = append(children, p)
		} else {
			info, err := entry.Info()
			if err == nil && info.Mode().IsRegular() {
				if progress != nil && reportFiles {
					progress(p, info.Size(), false)
				}
				size += info.Size()
			}
		}
	}
	return children, size
}

func buildExcludeSet(excludes []string) map[string]bool {
	m := make(map[string]bool, len(excludes))
	for _, e := range excludes {
		m[e] = true
	}
	return m
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
