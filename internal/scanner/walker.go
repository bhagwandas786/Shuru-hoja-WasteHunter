package scanner

import (
    "os"
    "path/filepath"
)

// Walker provides directory walking functionality
type Walker struct {
    skipDirs map[string]bool
}

// NewWalker creates a new Walker instance
func NewWalker(skipDirs []string) *Walker {
    w := &Walker{
        skipDirs: make(map[string]bool),
    }
    for _, dir := range skipDirs {
        w.skipDirs[dir] = true
    }
    return w
}

// ShouldSkip checks if a directory should be skipped
func (w *Walker) ShouldSkip(path string) bool {
    // Check if this exact path should be skipped
    if w.skipDirs[path] {
        return true
    }
    
    // Check parent directories
    dir := path
    for dir != "/" && dir != "." {
        if w.skipDirs[dir] {
            return true
        }
        dir = filepath.Dir(dir)
    }
    
    return false
}

// IsAccessible checks if we can access a path
func IsAccessible(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
