package dirwalker

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type ScanEntry struct {
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"is_dir"`
	Human string `json:"human"`
}

type ScanSummary struct {
	TotalSize int64       `json:"total_size"`
	TotalStr  string      `json:"total_str"`
	Entries   []ScanEntry `json:"entries,omitempty"`
}

func RenderTable(summary ScanSummary) string {
	var b strings.Builder

	colW := 60
	for _, e := range summary.Entries {
		if len(e.Path) > colW {
			colW = len(e.Path)
		}
	}

	sep := func() {
		fmt.Fprintf(&b, "+%s+%s+\n", strings.Repeat("-", colW+2), strings.Repeat("-", 14))
	}

	row := func(path, size string) {
		fmt.Fprintf(&b, "| %-*s | %*s |\n", colW, path, 12, size)
	}

	if len(summary.Entries) > 0 {
		sep()
		row("Path", "Size")
		sep()
		for _, e := range summary.Entries {
			row(e.Path, e.Human)
		}
		sep()
	}

	fmt.Fprintf(&b, "Total: %s\n", summary.TotalStr)
	return b.String()
}

func RenderJSON(summary ScanSummary) string {
	b, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshal error: %v", err)
	}
	return string(b)
}

func RenderYAML(summary ScanSummary) string {
	b, err := yaml.Marshal(summary)
	if err != nil {
		return fmt.Sprintf("yaml marshal error: %v", err)
	}
	return string(b)
}
