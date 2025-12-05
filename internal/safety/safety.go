package safety

import (
    "os"
    "path/filepath"
    "syscall"
)

// IsSafeToScan checks if we should scan a path
func IsSafeToScan(path string) bool {
    // Check if it's a system directory
    if isSystemDirectory(path) {
        return false
    }
    
    // Check permissions
    if !hasReadPermission(path) {
        return false
    }
    
    // Check if it's a dangerous path
    if isDangerousPath(path) {
        return false
    }
    
    return true
}

func isSystemDirectory(path string) bool {
    systemDirs := []string{
        "/proc", "/sys", "/dev", "/run",
        "/boot", "/snap", "/var/lib/docker",
    }
    
    for _, dir := range systemDirs {
        if path == dir || strings.HasPrefix(path, dir+"/") {
            return true
        }
    }
    
    return false
}

func hasReadPermission(path string) bool {
    // Check if we can stat the file/directory
    if _, err := os.Stat(path); err != nil {
        return false
    }
    
    // Check if we can open for reading
    if filepath.Base(path) != "" {
        file, err := os.Open(path)
        if err != nil {
            return false
        }
        file.Close()
    }
    
    return true
}

func isDangerousPath(path string) bool {
    dangerousPatterns := []string{
        "/etc/shadow", "/etc/passwd", "/etc/gshadow",
        "/root/", "/var/lib/", "/usr/lib/",
    }
    
    for _, pattern := range dangerousPatterns {
        if strings.Contains(path, pattern) {
            return true
        }
    }
    
    return false
}

// SetResourceLimits applies resource limits for safety
func SetResourceLimits(maxMemoryMB, maxOpenFiles int) error {
    // Set memory limit
    if maxMemoryMB > 0 {
        var limit syscall.Rlimit
        limit.Cur = uint64(maxMemoryMB) * 1024 * 1024
        limit.Max = uint64(maxMemoryMB) * 1024 * 1024
        if err := syscall.Setrlimit(syscall.RLIMIT_AS, &limit); err != nil {
            return err
        }
    }
    
    // Set file descriptor limit
    if maxOpenFiles > 0 {
        var limit syscall.Rlimit
        limit.Cur = uint64(maxOpenFiles)
        limit.Max = uint64(maxOpenFiles)
        if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
            return err
        }
    }
    
    return nil
}
