# dirwalker — parallel directory scanner

Parallel, depth-aware directory walker. Used as a library and as a standalone CLI binary (embedded in [pvdu](https://github.com/neutry/pvdu) for remote pod execution).

## Library

```go
import "github.com/neutry/dirwalker"

total, err := dirwalker.ScanDirectory("/path", 0, nil, nil, 4, false)
```

- `ScanDirectory(root, maxDepth, excludes, progress, workers, reportFiles)`
- `ProgressFn func(path string, size int64, isDir bool)`

## CLI

```bash
dirwalker [path] [flags]
dirwalker /mnt/data -d 3 -w 8 --files
```

| Flag | Short | Description |
|------|-------|-------------|
| `--max-depth` | `-d` | Max depth (0 = unlimited) |
| `--exclude` | `-e` | Comma-separated paths to exclude |
| `--workers` | `-w` | Parallel workers (0 = auto) |
| `--files` | `-f` | Report individual file sizes |

Outputs JSON lines with `"type":"dir"`, `"type":"file"`, and `"type":"done"` (total).
