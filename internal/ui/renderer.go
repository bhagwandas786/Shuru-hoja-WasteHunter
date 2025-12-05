package ui

import (
    "fmt"
    "os"
    "time"

    "github.com/olekukonko/tablewriter"
    "shuru-hoja/pkg/types"
)

var (
    ColorReset  = "\033[0m"
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorBlue   = "\033[34m"
    ColorPurple = "\033[35m"
    ColorCyan   = "\033[36m"
    ColorWhite  = "\033[37m"
)

func ShowWelcome() {
    fmt.Println(ColorCyan + "┌─────────────────────────────────────────────┐" + ColorReset)
    fmt.Println(ColorCyan + "│        " + ColorWhite + "SHURU HOJA - Filesystem Analyzer" + ColorCyan + "        │" + ColorReset)
    fmt.Println(ColorCyan + "│    " + ColorYellow + "Production-Safe • Read-Only • Enterprise" + ColorCyan + "   │" + ColorReset)
    fmt.Println(ColorCyan + "└─────────────────────────────────────────────┘" + ColorReset)
    fmt.Println()
    fmt.Println(ColorYellow + "⚠  Scanning filesystem... (This may take a while)" + ColorReset)
    fmt.Println(ColorYellow + "⚠  Running in read-only mode - No files will be modified" + ColorReset)
    fmt.Println()
}

func RenderResults(results []types.ScanResult, duration time.Duration, maxResults int) {
    // Calculate summary
    summary := CalculateSummary(results)
    
    // Show summary first
    showSummary(summary, duration)
    
    // Show table of top findings
    ShowTopFindings(results, maxResults)
    
    // Show recommendations
    ShowRecommendations(results)
}

func showSummary(summary Summary, duration time.Duration) {
    fmt.Println()
    fmt.Println(ColorCyan + "══════════════════════════════════════════════════════════" + ColorReset)
    fmt.Println(ColorWhite + "                     SCAN SUMMARY" + ColorReset)
    fmt.Println(ColorCyan + "══════════════════════════════════════════════════════════" + ColorReset)
    
    table := tablewriter.NewWriter(os.Stdout)
    table.SetBorder(false)
    table.SetColumnSeparator("  ")
    table.SetAlignment(tablewriter.ALIGN_LEFT)
    table.SetAutoWrapText(false)
    
    data := [][]string{
        {"Total Scanned:", fmt.Sprintf("%.2f GB", float64(summary.TotalScannedBytes)/(1024*1024*1024))},
        {"Total Files:", fmt.Sprintf("%d", summary.TotalScannedFiles)},
        {"Total Directories:", fmt.Sprintf("%d", summary.TotalScannedDirs)},
        {"Potential Cleanup:", fmt.Sprintf("%.2f GB", float64(summary.PotentialCleanup)/(1024*1024*1024))},
        {"Critical Risk Items:", fmt.Sprintf("%d", summary.CriticalRiskCount)},
        {"Caution Risk Items:", fmt.Sprintf("%d", summary.CautionRiskCount)},
        {"Scan Duration:", fmt.Sprintf("%.2f seconds", duration.Seconds())},
    }
    
    table.AppendBulk(data)
    table.Render()
    fmt.Println()
}
