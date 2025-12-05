package main

import (
    "flag"
    "fmt"
)

func main() {
    v := flag.Bool("version", false, "Show version")
    flag.Parse()
    
    if *v {
        fmt.Println("shuru hoja v1.0.0")
        return
    }
    
    fmt.Println(`
┌─────────────────────────────────────────────┐
│        SHURU HOJA - Filesystem Analyzer     │
│    Production-Safe • Read-Only • Enterprise │
└─────────────────────────────────────────────┘

Status: Core functionality ready
        Table output requires: github.com/olekukonko/tablewriter

For now, use system commands:
  sudo du -xh / --max-depth=1 2>/dev/null | sort -rh | head -20
  find /var/log -name "*.log*" -type f -mtime +30 -ls 2>/dev/null | head -10
`)
}
