# **🚀SHURU HOJA - Production-Safe Filesystem Analyzer🚀**

**Version:** 1.0.0

**Developed by:** Bunny 🐯

**Tagline:** Start Now, Analyze Safely!

**Repository:** https://github.com/bhagwandas786/Shuru-hoja-WasteHunter

---

## **🚀 Overview**

**SHURU HOJA** is an enterprise-grade, production-safe filesystem analysis tool designed for system administrators and DevOps engineers. It performs comprehensive filesystem analysis in **100% read-only mode** - meaning it **never modifies or deletes any files**.

Perfect for production servers where safety is paramount, this tool helps identify space-consuming files and provides actionable recommendations for cleanup.

---

## **Project Structure**

<img width="1567" height="380" alt="image" src="https://github.com/user-attachments/assets/906bc4cc-5f30-4133-ae02-d4d7c2a1813b" />

## **✨ Features**

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

## 📦 Installation

### Quick Install & Setup (Prerequisites Installation)

```bash
# Update system packages
sudo apt update && sudo apt upgrade -y

# Install Go compiler and Git
sudo apt install -y golang-go git

# Verify installations
go version && git --version
```
## **Clone Repository**
```bash
# Clone the repository
git clone https://github.com/bhagwandas786/Shuru-hoja-WasteHunter.git

# Navigate to project directory
cd Shuru-hoja-WasteHunter
```
## **Basic Commands**
```bash
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
```
## **Scan Specific Locations**
```bash
# Scan home directory
shuru-hoja --path /home

# Scan log files
sudo shuru-hoja --path /var/log

# Scan temporary files
shuru-hoja --path /tmp
```
## **Example Output**

<img width="1288" height="476" alt="image" src="https://github.com/user-attachments/assets/6d28fcb7-edcb-4c88-bc22-aa07ad5c3ca4" />

<img width="1060" height="240" alt="image" src="https://github.com/user-attachments/assets/875adf00-449b-4aee-8fe9-ce8331af88d1" />

## **Risk Levels:**

🟢 Safe: Normal system files, recently accessed

🟡 Caution: Old logs, moderate caches, potential cleanup

🔴 Critical: Large duplicates, massive caches, immediate action needed

## *📄 License*
This tool is developed by Bunny for production use. Modify and distribute as needed.
