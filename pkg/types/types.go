package types

import (
    "os"
    "time"
)

type FileType string

const (
    TypeFile      FileType = "file"
    TypeDirectory FileType = "directory"
    TypeLog       FileType = "log"
    TypeCache     FileType = "cache"
    TypeTemp      FileType = "temp"
    TypeBackup    FileType = "backup"
    TypeDuplicate FileType = "duplicate"
    TypeOrphan    FileType = "orphan"
)

type RiskLevel string

const (
    RiskSafe     RiskLevel = "Safe"
    RiskCaution  RiskLevel = "Caution"
    RiskCritical RiskLevel = "Critical"
)

type Recommendation string

const (
    RecKeep    Recommendation = "Keep"
    RecReview  Recommendation = "Review"
    RecDelete  Recommendation = "Delete"
)

type FileInfo struct {
    Path          string
    Size          int64
    IsDir         bool
    Mode          os.FileMode
    ModTime       time.Time
    AccessTime    time.Time
    UID           uint32
    GID           uint32
    HardLinks     uint64
    Inode         uint64
}

type ScanResult struct {
    Info           FileInfo
    Type           FileType
    RiskLevel      RiskLevel
    Recommendation Recommendation
    Reason         string
    DuplicateGroup string
    AgeDays        int
}

type Summary struct {
    TotalScannedBytes   int64
    TotalScannedFiles   int64
    TotalScannedDirs    int64
    PotentialCleanup    int64
    CriticalRiskCount   int64
    CautionRiskCount    int64
    ScanDuration       time.Duration
}
