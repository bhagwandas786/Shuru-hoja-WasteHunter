package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "shuru-hoja/internal/analyzer"
    "shuru-hoja/internal/scanner"
    "shuru-hoja/internal/ui"
    "shuru-hoja/internal/config"
    "shuru-hoja/pkg/types"
)

func main() {
    // Setup signal handling for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
        fmt.Println("\n\nScan interrupted by user. Exiting safely...")
        os.Exit(0)
    }()

    // Start analysis
    if err := runAnalysis(ctx); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}

func runAnalysis(ctx context.Context) error {
    startTime := time.Now()
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    // Initialize scanner
    scanner := scanner.NewConcurrentScanner(cfg.MaxWorkers)
    
    // Initialize analyzer with detection rules
    analyzer := analyzer.NewAnalyzer(scanner, cfg)
    
    // Start scanning
    ui.ShowWelcome()
    
    results, err := analyzer.Analyze(ctx, "/")
    if err != nil {
        return fmt.Errorf("analysis failed: %w", err)
    }
    
    // Render results
    duration := time.Since(startTime)
    ui.RenderResults(results, duration, cfg)
    
    return nil
}
