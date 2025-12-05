package config

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type Config struct {
    General      GeneralConfig
    Detection    DetectionConfig
    Risk         RiskConfig
    Output       OutputConfig
    Safety       SafetyConfig
}

type GeneralConfig struct {
    MaxWorkers int
    SkipPaths  []string
    MaxDepth   int
}

type DetectionConfig struct {
    LogFileAgeDays      int
    OrphanDirAgeDays    int
    OrphanDirMinSize    int64
    CacheMinSize        int64
    DuplicateMinSize    int64
    NodeModulesMaxSize  int64
    PythonVenvMaxSize   int64
    DockerCacheMaxSize  int64
    JournalLogMaxSize   int64
}

type RiskConfig struct {
    CriticalSizeGB int64
    CautionSizeGB  int64
    CriticalAgeDays int
    CautionAgeDays  int
}

type OutputConfig struct {
    Format            string
    Color             bool
    MaxResults        int
    TruncatePathLength int
}

type SafetyConfig struct {
    ReadOnly         bool
    MaxFileSize      int64
    SkipSystemDirs   bool
    MaxMemoryMB      int
    MaxOpenFiles     int
}

func Load() (*Config, error) {
    cfg := defaultConfig()
    
    // Try to load from /etc/shuruhoja.conf
    if err := loadFromFile("/etc/shuruhoja.conf", cfg); err != nil {
        // Try local config
        home, _ := os.UserHomeDir()
        localConfig := filepath.Join(home, ".config", "shuruhoja.conf")
        if err := loadFromFile(localConfig, cfg); err != nil {
            // Use defaults
            fmt.Println("Using default configuration")
        }
    }
    
    return cfg, nil
}

func defaultConfig() *Config {
    return &Config{
        General: GeneralConfig{
            MaxWorkers: 100,
            SkipPaths:  []string{"/proc", "/sys", "/dev", "/run"},
            MaxDepth:   0, // Unlimited
        },
        Detection: DetectionConfig{
            LogFileAgeDays:      30,
            OrphanDirAgeDays:    90,
            OrphanDirMinSize:    1 * 1024 * 1024 * 1024, // 1GB
            CacheMinSize:        100 * 1024 * 1024,      // 100MB
            DuplicateMinSize:    10 * 1024 * 1024,       // 10MB
            NodeModulesMaxSize:  500 * 1024 * 1024,      // 500MB
            PythonVenvMaxSize:   1 * 1024 * 1024 * 1024, // 1GB
            DockerCacheMaxSize:  5 * 1024 * 1024 * 1024, // 5GB
            JournalLogMaxSize:   2 * 1024 * 1024 * 1024, // 2GB
        },
        Risk: RiskConfig{
            CriticalSizeGB: 10,
            CautionSizeGB:  1,
            CriticalAgeDays: 365,
            CautionAgeDays:  180,
        },
        Output: OutputConfig{
            Format:            "table",
            Color:             true,
            MaxResults:        50,
            TruncatePathLength: 80,
        },
        Safety: SafetyConfig{
            ReadOnly:       true,
            MaxFileSize:    10 * 1024 * 1024 * 1024, // 10GB
            SkipSystemDirs: true,
            MaxMemoryMB:    1024,
            MaxOpenFiles:   10000,
        },
    }
}

func loadFromFile(path string, cfg *Config) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    
    lines := strings.Split(string(data), "\n")
    currentSection := ""
    
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
            currentSection = line[1 : len(line)-1]
            continue
        }
        
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        
        parseConfigValue(currentSection, key, value, cfg)
    }
    
    return nil
}

func parseConfigValue(section, key, value string, cfg *Config) {
    switch section {
    case "general":
        switch key {
        case "max_workers":
            if v, err := strconv.Atoi(value); err == nil {
                cfg.General.MaxWorkers = v
            }
        case "skip_paths":
            cfg.General.SkipPaths = strings.Split(value, ",")
        case "max_depth":
            if v, err := strconv.Atoi(value); err == nil {
                cfg.General.MaxDepth = v
            }
        }
    case "detection":
        switch key {
        case "log_file_age_days":
            if v, err := strconv.Atoi(value); err == nil {
                cfg.Detection.LogFileAgeDays = v
            }
        case "orphan_dir_age_days":
            if v, err := strconv.Atoi(value); err == nil {
                cfg.Detection.OrphanDirAgeDays = v
            }
        }
    // Add more cases for other sections
    }
}
