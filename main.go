package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"pmjtoca/largefileRWer/config"
	"pmjtoca/largefileRWer/internal/logger"
	"pmjtoca/largefileRWer/read"
	"pmjtoca/largefileRWer/write"
)

var (
	// Command-line flags
	inputFile  = flag.String("input", "", "Input file path (required)")
	outputFile = flag.String("output", "", "Output file path (required)")
	configFile = flag.String("config", "", "Configuration file (CUE or JSON)")
	profile    = flag.String("profile", "", "Configuration profile (production, development, nvme, hdd, network, lowmem, auto)")

	// Override flags
	workers     = flag.Int("workers", 0, "Override: number of workers")
	bufferSize  = flag.Int("buffer", 0, "Override: buffer size in KB")
	chunkSize   = flag.Int("chunk", 0, "Override: chunk size in MB")
	useMmap     = flag.Bool("mmap", false, "Override: use memory-mapped I/O")
	useDirectIO = flag.Bool("directio", true, "Override: use direct I/O")
	verify      = flag.Bool("verify", true, "Override: verify checksums")
	verbose     = flag.Bool("verbose", false, "Override: verbose logging")
	logLevel    = flag.String("log-level", "", "Override: log level (debug, info, warn, error)")
	checkpoint  = flag.Bool("checkpoint", true, "Override: enable checkpoint")
)

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	flag.Parse()

	// Validate required flags
	if *inputFile == "" || *outputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: input and output paths are required\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -input <file> -output <file> [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Setup logging
	if err := logger.Setup(cfg.LogLevel, cfg.LogFormat, cfg.LogFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	slog.Info("Starting Large File Processor",
		"version", Version,
		"build_time", BuildTime,
		"git_commit", GitCommit,
		"go_version", runtime.Version())

	// Apply runtime configuration
	debug.SetGCPercent(cfg.GCPercent)
	debug.SetMemoryLimit(cfg.GetMemoryLimitBytes())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Log configuration
	if cfg.EnableVerbose {
		jsonCfg, _ := cfg.ToJSON()
		slog.Debug("Configuration", "config", string(jsonCfg))
	}

	// Create checkpoint directory if needed
	if cfg.EnableCheckpoint {
		if err := os.MkdirAll(cfg.CheckpointDir, 0755); err != nil {
			slog.Warn("Failed to create checkpoint directory", "error", err)
		}
	}

	// Create reader and writer
	reader, err := read.NewReader(cfg)
	if err != nil {
		slog.Error("Failed to create reader", "error", err)
		os.Exit(1)
	}
	defer reader.Close()

	writer, err := write.NewWriter(cfg)
	if err != nil {
		slog.Error("Failed to create writer", "error", err)
		os.Exit(1)
	}
	defer writer.Close()

	// Process the file
	ctx := context.Background()
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.GetTimeoutDuration())
		defer cancel()
	}

	startTime := time.Now()
	slog.Info("Starting file processing",
		"input", cfg.InputPath,
		"output", cfg.OutputPath)

	// Read chunks
	chunkChan, readErrChan := reader.ReadAll(ctx)

	// Convert to write chunks
	writeChunkChan := make(chan *write.WriteChunk, cfg.QueueSize)
	go func() {
		defer close(writeChunkChan)
		for chunk := range chunkChan {
			if chunk.Error != nil {
				slog.Error("Read error", "error", chunk.Error)
				return
			}

			writeChunkChan <- &write.WriteChunk{
				ID:       chunk.ID,
				Offset:   chunk.Offset,
				Data:     chunk.Data,
				Checksum: chunk.Checksum,
			}
		}
	}()

	// Write chunks in order
	if err := writer.WriteChunksOrdered(ctx, writeChunkChan); err != nil {
		slog.Error("Write failed", "error", err)
		os.Exit(1)
	}

	// Check for read errors
	select {
	case readErr := <-readErrChan:
		if readErr != nil {
			slog.Error("Read failed", "error", readErr)
			os.Exit(1)
		}
	default:
	}

	// Flush to ensure all data is on disk
	if err := writer.Flush(); err != nil {
		slog.Warn("Flush failed", "error", err)
	}

	elapsed := time.Since(startTime)
	readMetrics := reader.GetMetrics()
	writeMetrics := writer.GetMetrics()

	slog.Info("Processing completed successfully",
		"duration", elapsed.Round(time.Second),
		"bytes_read_gb", float64(readMetrics.GetBytesRead())/(1024*1024*1024),
		"bytes_written_gb", float64(writeMetrics.GetBytesWritten())/(1024*1024*1024),
		"read_speed_mbps", readMetrics.GetSpeedMBps(true),
		"write_speed_mbps", writeMetrics.GetSpeedMBps(false),
		"avg_speed_mbps", float64(readMetrics.GetBytesRead())/elapsed.Seconds()/(1024*1024))
}

func loadConfiguration() (*config.Config, error) {
	loader := config.NewLoader()

	// Start with default config
	baseCfg := config.DefaultConfig()

	// Apply profile if specified
	switch *profile {
	case "production":
		baseCfg = config.ProductionConfig()
		slog.Info("Using production profile")
	case "development":
		baseCfg = config.DevelopmentConfig()
		slog.Info("Using development profile")
	case "nvme":
		baseCfg = config.NVMeProfile()
		slog.Info("Using NVMe profile")
	case "hdd":
		baseCfg = config.HDDProfile()
		slog.Info("Using HDD profile")
	case "network":
		baseCfg = config.NetworkProfile()
		slog.Info("Using network profile")
	case "lowmem":
		baseCfg = config.LowMemoryProfile()
		slog.Info("Using low memory profile")
	case "auto":
		baseCfg = loader.GenerateProfile()
		slog.Info("Using auto-detected profile")
	case "":
		// No profile, use default
	default:
		return nil, fmt.Errorf("unknown profile: %s", *profile)
	}

	// Load from config file if provided
	if *configFile != "" {
		var fileCfg *config.Config
		var err error

		if strings.HasSuffix(*configFile, ".cue") {
			fileCfg, err = loader.LoadFromFile(*configFile)
		} else {
			fileCfg, err = loader.LoadFromJSON(*configFile)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}

		// Merge file config with base (file takes precedence)
		baseCfg, err = loader.Merge(baseCfg, fileCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to merge config file: %w", err)
		}
		slog.Info("Loaded configuration from file", "path", *configFile)
	}

	// Load from environment variables (env takes precedence over file)
	envCfg, err := loader.LoadFromEnvironment()
	if err != nil {
		slog.Warn("Failed to load environment config", "error", err)
	} else if envCfg != nil {
		baseCfg, err = loader.Merge(baseCfg, envCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to merge environment config: %w", err)
		}
		slog.Info("Loaded configuration from environment")
	}

	// Apply command-line overrides (highest priority)
	// Only override if the flag was explicitly set (non-zero values)
	overrideCfg := &config.Config{
		InputPath:        *inputFile,
		OutputPath:       *outputFile,
		EnableVerbose:    *verbose,
		VerifyChecksum:   *verify,
		UseMmap:          *useMmap,
		UseDirectIO:      *useDirectIO,
		EnableCheckpoint: *checkpoint,
	}

	if *workers > 0 {
		overrideCfg.Workers = *workers
	}
	if *bufferSize > 0 {
		overrideCfg.BufferSizeKB = *bufferSize
	}
	if *chunkSize > 0 {
		overrideCfg.ChunkSizeMB = *chunkSize
	}
	if *logLevel != "" {
		overrideCfg.LogLevel = *logLevel
	}

	// Merge overrides (they take highest priority)
	finalCfg, err := loader.Merge(baseCfg, overrideCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to apply overrides: %w", err)
	}

	// Validate final configuration
	if err := finalCfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return finalCfg, nil
}
