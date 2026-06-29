# dirwalker — parallel directory scanner

Parallel, depth-aware directory walker. Used as a library and as a standalone CLI binary by [pvdu](https://github.com/NeutryFD/pvdu) for Kubernetes PVC usage scanning.

## Install

```bash
go install github.com/NeutryFD/dirwalker/cmd/dirwalker@latest
```

## CLI

```bash
dirwalker [path] [flags]
dirwalker /mnt/data -d 3 -w 8 --files --output=table
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--max-depth` | `-d` | Max depth (0 = unlimited) | `0` |
| `--exclude` | `-e` | Comma-separated paths to exclude | `""` |
| `--workers` | `-w` | Parallel workers (0 = auto) | `0` |
| `--files` | `-f` | Report individual file sizes | `false` |
| `--output` | `-o` | Output format: table, json, yaml, json-lines | `json-lines` |

### Output formats

```
dirwalker /data -o table     # human-readable table
dirwalker /data -o json      # JSON object with total + entries
dirwalker /data -o yaml      # YAML output
dirwalker /data              # JSON Lines (streaming, default)
```

## Library

```go
import "github.com/NeutryFD/dirwalker"

total, err := dirwalker.ScanDirectory("/path", 0, nil, nil, 4, false)
```

- `ScanDirectory(root, maxDepth, excludes, progress, workers, reportFiles)` — scans a directory tree in parallel
- `ProgressFn func(path string, size int64, isDir bool)` — called for each directory/file
- `FormatBytes(b int64) string` — `"1.5 GiB"`, `"500 B"`
- `FormatBytesShort(b int64) string` — `"1.5Gi"`, `"500 B"`
- `RenderTable(summary ScanSummary) string` — table output
- `RenderJSON(summary ScanSummary) string` — JSON output
- `RenderYAML(summary ScanSummary) string` — YAML output
