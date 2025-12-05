package analyzer

import (
    "context"
    "fmt"
    "os"
    "sort"
    "time"

    "shuru-hoja/internal/config"
    "shuru-hoja/internal/scanner"
    "shuru-hoja/internal/analyzer/detectors"
    "shuru-hoja/pkg/types"
)

type Analyzer struct {
    scanner     *scanner.ConcurrentScanner
    config      *config.Config
    detectors   []detectors.Detector
}

func NewAnalyzer(s *scanner.ConcurrentScanner, cfg *config.Config) *Analyzer {
    analyzer := &Analyzer{
        scanner: s,
        config:  cfg,
    }
    
    // Initialize detectors
    analyzer.initDetectors()
    
    return analyzer
}

func (a *Analyzer) initDetectors() {
    // Create all detectors
    a.detectors = []detectors.Detector{
        detectors.NewLogFileDetector(a.config.Detection.LogFileAgeDays),
    }
}

func (a *Analyzer) Analyze(ctx context.Context, root string) ([]types.ScanResult, error) {
    results := []types.ScanResult{}
    
    // Start scanning
    fileChan, errChan := a.scanner.Scan(ctx, root)
    
    // Process files
    for {
        select {
        case <-ctx.Done():
            return results, ctx.Err()
            
        case fileInfo, ok := <-fileChan:
            if !ok {
                fileChan = nil
            } else {
                result := a.analyzeFile(fileInfo)
                if result != nil {
                    results = append(results, *result)
                }
            }
            
        case err, ok := <-errChan:
            if !ok {
                errChan = nil
            } else {
                // Log permission errors but continue
                fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
            }
        }
        
        if fileChan == nil && errChan == nil {
            break
        }
    }
    
    // Sort by size (largest first)
    sort.Slice(results, func(i, j int) bool {
        return results[i].Info.Size > results[j].Info.Size
    })
    
    return results, nil
}

func (a *Analyzer) analyzeFile(info types.FileInfo) *types.ScanResult {
    // Run through all detectors
    for _, detector := range a.detectors {
        if result := detector.Detect(info); result != nil {
            return result
        }
    }
    
    // Default result for files not caught by detectors
    result := &types.ScanResult{
        Info: info,
        Type: types.TypeFile,
        RiskLevel: types.RiskSafe,
        Recommendation: types.RecKeep,
    }
    
    if info.IsDir {
        result.Type = types.TypeDirectory
    }
    
    return result
}
