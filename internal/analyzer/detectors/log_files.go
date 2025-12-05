package detectors

import (
    "fmt"
    "path/filepath"
    "strings"
    "time"
    
    "shuru-hoja/pkg/types"
)

type LogFileDetector struct {
    MaxAgeDays int
    Patterns   []string
}

func NewLogFileDetector(maxAgeDays int) *LogFileDetector {
    return &LogFileDetector{
        MaxAgeDays: maxAgeDays,
        Patterns: []string{".log", ".log.", ".journal", ".gz", ".bz2"},
    }
}

func (d *LogFileDetector) Detect(info types.FileInfo) *types.ScanResult {
    if info.IsDir {
        return nil
    }
    
    filename := filepath.Base(info.Path)
    lowerName := strings.ToLower(filename)
    
    // Check if it's a log file
    isLogFile := false
    for _, pattern := range d.Patterns {
        if strings.Contains(lowerName, pattern) || 
           strings.Contains(info.Path, "/var/log/") ||
           strings.Contains(info.Path, "/var/logs/") {
            isLogFile = true
            break
        }
    }
    
    if !isLogFile {
        return nil
    }
    
    ageDays := int(time.Since(info.ModTime).Hours() / 24)
    
    result := &types.ScanResult{
        Info: info,
        Type: types.TypeLog,
        AgeDays: ageDays,
    }
    
    // Set risk level based on age and size
    if ageDays > d.MaxAgeDays {
        if info.Size > 100*1024*1024 { // 100MB
            result.RiskLevel = types.RiskCritical
            result.Recommendation = types.RecDelete
        } else {
            result.RiskLevel = types.RiskCaution
            result.Recommendation = types.RecReview
        }
        result.Reason = fmt.Sprintf("Old log file (%d days, %s)", 
            ageDays, formatSize(info.Size))
    } else {
        result.RiskLevel = types.RiskSafe
        result.Recommendation = types.RecKeep
    }
    
    return result
}

func formatSize(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
