package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/NeutryFD/dirwalker"
)

type progressLine struct {
	Type  string `json:"type"`
	Path  string `json:"path,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Human string `json:"human,omitempty"`
}

type doneLine struct {
	Type  string `json:"type"`
	Size  int64  `json:"size"`
	Human string `json:"human"`
}

func main() {
	var (
		maxDepth     int
		excludeStr   string
		workers      int
		reportFiles  bool
		outputFormat string
	)

	flag.IntVar(&maxDepth, "max-depth", 0, "max directory depth (0 = unlimited)")
	flag.IntVar(&maxDepth, "d", 0, "shorthand for --max-depth")
	flag.StringVar(&excludeStr, "exclude", "", "comma-separated paths to exclude")
	flag.StringVar(&excludeStr, "e", "", "shorthand for --exclude")
	flag.IntVar(&workers, "workers", 0, "number of parallel workers (0 = auto)")
	flag.IntVar(&workers, "w", 0, "shorthand for --workers")
	flag.BoolVar(&reportFiles, "files", false, "report individual file sizes")
	flag.BoolVar(&reportFiles, "f", false, "shorthand for --files")
	flag.StringVar(&outputFormat, "output", "json-lines", "output format (table, json, yaml, json-lines)")
	flag.StringVar(&outputFormat, "o", "json-lines", "shorthand for --output")
	flag.Parse()

	root := ""
	args := flag.Args()
	if len(args) > 0 {
		root = args[0]
		if len(args) > 1 {
			flag.CommandLine.Parse(args[1:])
		}
	}
	if root == "" {
		root = "."
	}

	var excludeList []string
	if excludeStr != "" {
		excludeList = strings.Split(excludeStr, ",")
	}

	var entries []dirwalker.ScanEntry

	enc := json.NewEncoder(os.Stdout)
	streaming := outputFormat == "json-lines"

	progress := func(p string, size int64, isDir bool) {
		if streaming {
			typ := "dir"
			if !isDir {
				typ = "file"
			}
			_ = enc.Encode(progressLine{
				Type:  typ,
				Path:  p,
				Size:  size,
				Human: dirwalker.FormatBytesShort(size),
			})
		} else {
			entries = append(entries, dirwalker.ScanEntry{
				Path:  p,
				Size:  size,
				IsDir: isDir,
				Human: dirwalker.FormatBytesShort(size),
			})
		}
	}

	total, err := dirwalker.ScanDirectory(root, maxDepth, excludeList, progress, workers, reportFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	summary := dirwalker.ScanSummary{
		TotalSize: total,
		TotalStr:  dirwalker.FormatBytesShort(total),
		Entries:   entries,
	}

	switch outputFormat {
	case "json":
		fmt.Println(dirwalker.RenderJSON(summary))
	case "yaml":
		fmt.Println(dirwalker.RenderYAML(summary))
	case "table":
		fmt.Print(dirwalker.RenderTable(summary))
	default:
		_ = enc.Encode(doneLine{
			Type:  "done",
			Size:  total,
			Human: dirwalker.FormatBytesShort(total),
		})
	}
}
