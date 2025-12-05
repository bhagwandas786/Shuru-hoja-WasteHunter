package scanner

import (
    "context"
    "os"
    "path/filepath"
    "sync"
    "sync/atomic"
    "syscall"

    "shuru-hoja/pkg/types"
)

type ConcurrentScanner struct {
    maxWorkers   int
    results      chan types.FileInfo
    errors       chan error
    wg           sync.WaitGroup
    scannedFiles int64
    scannedDirs  int64
    totalSize    int64
}

func NewConcurrentScanner(maxWorkers int) *ConcurrentScanner {
    if maxWorkers <= 0 {
        maxWorkers = 100
    }
    return &ConcurrentScanner{
        maxWorkers: maxWorkers,
        results:    make(chan types.FileInfo, 10000),
        errors:     make(chan error, 100),
    }
}

func (s *ConcurrentScanner) Scan(ctx context.Context, root string) (<-chan types.FileInfo, <-chan error) {
    go func() {
        defer close(s.results)
        defer close(s.errors)
        
        semaphore := make(chan struct{}, s.maxWorkers)
        var scanWg sync.WaitGroup
        
        s.walkDir(ctx, root, semaphore, &scanWg)
        scanWg.Wait()
    }()
    
    return s.results, s.errors
}

func (s *ConcurrentScanner) walkDir(ctx context.Context, path string, sem chan struct{}, wg *sync.WaitGroup) {
    select {
    case <-ctx.Done():
        return
    default:
    }
    
    sem <- struct{}{}
    wg.Add(1)
    
    go func() {
        defer func() {
            <-sem
            wg.Done()
        }()
        
        entries, err := os.ReadDir(path)
        if err != nil {
            if os.IsPermission(err) {
                s.errors <- &PermissionError{Path: path, Err: err}
            }
            return
        }
        
        for _, entry := range entries {
            select {
            case <-ctx.Done():
                return
            default:
            }
            
            fullPath := filepath.Join(path, entry.Name())
            
            if s.shouldSkip(fullPath) {
                continue
            }
            
            info, err := entry.Info()
            if err != nil {
                s.errors <- err
                continue
            }
            
            fileInfo := s.createFileInfo(fullPath, info)
            
            if info.IsDir() {
                atomic.AddInt64(&s.scannedDirs, 1)
                s.walkDir(ctx, fullPath, sem, wg)
            } else {
                atomic.AddInt64(&s.scannedFiles, 1)
                atomic.AddInt64(&s.totalSize, info.Size())
            }
            
            select {
            case s.results <- fileInfo:
            case <-ctx.Done():
                return
            }
        }
    }()
}

func (s *ConcurrentScanner) shouldSkip(path string) bool {
    skipPaths := []string{
        "/proc", "/sys", "/dev", "/run",
        ".snapshot", ".zfs",
    }
    
    for _, skip := range skipPaths {
        if path == skip {
            return true
        }
    }
    return false
}

func (s *ConcurrentScanner) createFileInfo(path string, info os.FileInfo) types.FileInfo {
    sys := info.Sys()
    
    var uid, gid uint32
    var inode uint64
    
    if stat, ok := sys.(*syscall.Stat_t); ok {
        uid = stat.Uid
        gid = stat.Gid
        inode = stat.Ino
    }
    
    return types.FileInfo{
        Path:       path,
        Size:       info.Size(),
        IsDir:      info.IsDir(),
        Mode:       info.Mode(),
        ModTime:    info.ModTime(),
        AccessTime: info.ModTime(),
        UID:        uid,
        GID:        gid,
        Inode:      inode,
    }
}
