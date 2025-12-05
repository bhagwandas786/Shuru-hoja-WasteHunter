package scanner

import (
    "fmt"
    "os"
    "shuru-hoja/pkg/types"
)

type PermissionError struct {
    Path string
    Err  error
}

func (e *PermissionError) Error() string {
    return fmt.Sprintf("permission denied: %s: %v", e.Path, e.Err)
}

type Scanner interface {
    Scan(root string) (<-chan types.FileInfo, <-chan error)
    GetStats() Stats
}

type Stats struct {
    FilesScanned int64
    DirsScanned  int64
    TotalSize    int64
    Errors       int64
}

func isSymlink(info os.FileInfo) bool {
    return info.Mode()&os.ModeSymlink != 0
}

func shouldSkipSystemPath(path string) bool {
    systemPaths := []string{
        "/proc", "/sys", "/dev", "/run",
        "/var/lib/docker", "/snap",
    }
    
    for _, sysPath := range systemPaths {
        if path == sysPath {
            return true
        }
    }
    return false
}
