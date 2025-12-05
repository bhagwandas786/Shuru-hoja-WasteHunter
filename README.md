# SHURU HOJA - Production-Safe Filesystem Analyzer

**Version:** 1.0.0  
**Developed by:** Bunny 🐯  
**Tagline:** Start Now, Analyze Safely!

---

## *🚀 Overview*

**SHURU HOJA** is an enterprise-grade, production-safe filesystem analysis tool designed for system administrators and DevOps engineers. It performs comprehensive filesystem analysis in **100% read-only mode** - meaning it **never modifies or deletes any files**. 

Perfect for production servers where safety is paramount, this tool helps identify space-consuming files and provides actionable recommendations for cleanup.

---

## *Project Structure*

shuru-hoja/
├── shuru-hoja*              # Main executable
├── cmd/shuru-hoja/main.go   # Entry point
├── internal/                # Core logic
│   ├── analyzer/           # Analysis engine
│   ├── scanner/            # Filesystem scanner
│   ├── ui/                 # User interface
│   ├── config/             # Configuration
│   └── safety/             # Safety features
├── pkg/types/              # Data types
├── scripts/                # Installation scripts
├── etc/                    # Configuration files
├── go.mod                  # Go module
└── README.md               # This file

## *✨ Features*

### 🔒 **Safety First**
- **100% Read-Only Operations** - No files are ever modified or deleted
- **Production-Safe** - Designed for enterprise environments
- **Permission Respect** - Gracefully handles permission errors
- **No Sensitive Data Collection** - Never stores or transmits file contents

### 📊 **Smart Analysis**
- **Large File Detection** - Identifies files consuming significant space
- **Old Log File Detection** - Finds log files older than 30 days
- **Cache Directory Scanning** - Detects APT, NPM, Docker caches
- **Temporary File Identification** - Locates temp files in `/tmp`, `/var/tmp`
- **Duplicate File Detection** - Finds identical files via content hashing
- **Orphan Directory Detection** - Identifies large, unused directories

### 🎨 **Beautiful Interface**
- **Color-Coded Output** - Risk levels visually represented
- **Table Format Display** - Clean, organized results
- **Progress Indicators** - Real-time scanning progress
- **Summary Statistics** - Quick overview of findings

### ⚡ **Performance**
- **Concurrent Scanning** - Multi-threaded directory walking
- **Low Memory Footprint** - Efficient memory usage
- **Fast Processing** - Handles millions of files efficiently
- **Graceful Interrupt Handling** - Safe Ctrl+C termination

### 🛡️ **Enterprise Ready**
- **Configurable Rules** - Customize detection thresholds
- **Logging Support** - Optional debug logging
- **Resource Limits** - Configurable CPU/memory limits
- **Cross-Platform** - Works on Linux/Unix systems

---
## **Basic Commands**

# Show version
shuru-hoja --version

# Start full system analysis
sudo shuru-hoja

# Scan specific directory
shuru-hoja --path /home/user

# Quick scan mode
shuru-hoja --quick

# Show help
shuru-hoja --help

## **Example Output**
┌─────────────────────────────────────────────┐
│        SHURU HOJA - Filesystem Analyzer     │
│    Production-Safe • Read-Only • Enterprise │
└─────────────────────────────────────────────┘

════════════════════════════════════════════════
               SCAN SUMMARY
════════════════════════════════════════════════
Total Scanned:           245.67 GB
Total Files:             1,234,567
Total Directories:       89,123
Potential Cleanup:       45.23 GB
Critical Risk Items:     12
Caution Risk Items:      89
Scan Duration:           42.3 seconds

════════════════════════════════════════════════
          TOP CLEANUP RECOMMENDATIONS
════════════════════════════════════════════════
┌──────────┬───────────┬──────────┬──────────────┬────────────────────────────┐
│ Size     │ Type      │ Risk     │ Recommendation│ Path                      │
├──────────┼───────────┼──────────┼──────────────┼────────────────────────────┤
│ 12.4 GB  │ cache     │ Critical │ Delete       │ /var/cache/apt/archives   │
│ 8.2 GB   │ log       │ Caution  │ Review       │ /var/log/journal/*.journal│
│ 5.1 GB   │ duplicate │ Critical │ Delete       │ /home/user/backup.tar.gz  │
│ 4.7 GB   │ temp      │ Caution  │ Review       │ /tmp/large_temp_file      │
│ 3.2 GB   │ node_modules │ Critical │ Delete    │ /app/node_modules         │
└──────────┴───────────┴──────────┴──────────────┴────────────────────────────┘

## **Risk Levels:**

🟢 Safe: Normal system files, recently accessed

🟡 Caution: Old logs, moderate caches, potential cleanup

🔴 Critical: Large duplicates, massive caches, immediate action needed

## 📦 Installation

### Quick Install (Linux)
```
sudo ./scripts/install.sh
```
# Download the binary (if available)
# Or build from source as shown below
## *📄 License*
This tool is developed by Bunny for production use. Modify and distribute as needed.
