package ui

import (
    "fmt"
    "os"
    "strconv"

    "github.com/olekukonko/tablewriter"
    "shuru-hoja/pkg/types"
)

type Summary struct {
    TotalScannedBytes   int64
    TotalScannedFiles   int64
    TotalScannedDirs    int64
    PotentialCleanup    int64
    CriticalRiskCount   int64
    CautionRiskCount    int64
}

func CalculateSummary(results []types.ScanResult) Summary {
    var summary Summary
    
    for _, r := range results {
        summary.TotalScannedBytes += r.Info.Size
        if r.Info.IsDir {
            summary.TotalScannedDirs++
        } else {
            summary.TotalScannedFiles++
        }
        
        if r.Recommendation == types.RecDelete {
            summary.PotentialCleanup += r.Info.Size
        }
        
        switch r.RiskLevel {
        case types.RiskCritical:
            summary.CriticalRiskCount++
        case types.RiskCaution:
            summary.CautionRiskCount++
        }
    }
    
    return summary
}

func ShowTopFindings(results []types.ScanResult, maxResults int) {
    // Filter only items with recommendations
    var filtered []types.ScanResult
    for _, r := range results {
        if r.Recommendation != types.RecKeep && r.RiskLevel != types.RiskSafe {
            filtered = append(filtered, r)
        }
        if len(filtered) >= maxResults {
            break
        }
    }
    
    if len(filtered) == 0 {
        fmt.Println(ColorGreen + "✓ No cleanup recommendations found. System is clean!" + ColorReset)
        return
    }
    
    fmt.Println(ColorCyan + "══════════════════════════════════════════════════════════" + ColorReset)
    fmt.Println(ColorWhite + "                TOP CLEANUP RECOMMENDATIONS" + ColorReset)
    fmt.Println(ColorCyan + "══════════════════════════════════════════════════════════" + ColorReset)
    
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Size", "Type", "Risk", "Recommendation", "Path"})
    table.SetBorder(true)
    table.SetAutoWrapText(false)
    table.SetAutoFormatHeaders(true)
    
    for _, r := range filtered {
        size := FormatSize(r.Info.Size)
        riskColor := GetRiskColor(r.RiskLevel)
        recColor := GetRecommendationColor(r.Recommendation)
        
        row := []string{
            size,
            string(r.Type),
            riskColor + string(r.RiskLevel) + ColorReset,
            recColor + string(r.Recommendation) + ColorReset,
            TruncatePath(r.Info.Path, 50),
        }
        table.Append(row)
    }
    
    table.Render()
}

func FormatSize(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return strconv.FormatInt(bytes, 10) + " B"
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return strconv.FormatFloat(float64(bytes)/float64(div), 'f', 1, 64) + 
        " " + string("KMGTPE"[exp]) + "B"
}

func TruncatePath(path string, maxLength int) string {
    if len(path) <= maxLength {
        return path
    }
    
    // Keep beginning and end
    keep := (maxLength - 3) / 2
    return path[:keep] + "..." + path[len(path)-keep:]
}

func GetRiskColor(risk types.RiskLevel) string {
    switch risk {
    case types.RiskCritical:
        return ColorRed
    case types.RiskCaution:
        return ColorYellow
    default:
        return ColorGreen
    }
}

func GetRecommendationColor(rec types.Recommendation) string {
    switch rec {
    case types.RecDelete:
        return ColorRed
    case types.RecReview:
        return ColorYellow
    default:
        return ColorGreen
    }
}

func ShowRecommendations(results []types.ScanResult) {
    var critical, caution []types.ScanResult
    
    for _, r := range results {
        if r.RiskLevel == types.RiskCritical && r.Recommendation == types.RecDelete {
            critical = append(critical, r)
        } else if r.RiskLevel == types.RiskCaution && r.Recommendation == types.RecReview {
            caution = append(caution, r)
        }
    }
    
    if len(critical) > 0 {
        fmt.Println()
        fmt.Println(ColorRed + "══════════════════════════════════════════════════════════" + ColorReset)
        fmt.Println(ColorWhite + "                 CRITICAL RECOMMENDATIONS" + ColorReset)
        fmt.Println(ColorRed + "══════════════════════════════════════════════════════════" + ColorReset)
        
        for i, r := range critical {
            if i >= 10 {
                break
            }
            fmt.Printf("%s• %s%s - %s (%s)%s\n", 
                ColorRed, FormatSize(r.Info.Size), ColorReset,
                TruncatePath(r.Info.Path, 60),
                r.Reason, ColorReset)
        }
    }
    
    if len(caution) > 0 {
        fmt.Println()
        fmt.Println(ColorYellow + "══════════════════════════════════════════════════════════" + ColorReset)
        fmt.Println(ColorWhite + "                  CAUTION RECOMMENDATIONS" + ColorReset)
        fmt.Println(ColorYellow + "══════════════════════════════════════════════════════════" + ColorReset)
        
        for i, r := range caution {
            if i >= 10 {
                break
            }
            fmt.Printf("%s• %s%s - %s (%s)%s\n", 
                ColorYellow, FormatSize(r.Info.Size), ColorReset,
                TruncatePath(r.Info.Path, 60),
                r.Reason, ColorReset)
        }
    }
}
