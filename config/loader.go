package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

// Loader handles loading and merging CUE configurations
type Loader struct {
	ctx  *cue.Context
	base cue.Value
}

// NewLoader creates a new CUE configuration loader
func NewLoader() *Loader {
	ctx := cuecontext.New()

	// Load the base schema
	instances := load.Instances([]string{"."}, &load.Config{
		Dir:       ".",
		Package:   "config",
		Tags:      nil,
		Overlay:   nil,
		Stdin:     nil,
		DataFiles: true,
		Tests:     false,
		Tools:     false,
	})

	if len(instances) == 0 {
		return &Loader{ctx: ctx}
	}

	val := ctx.BuildInstance(instances[0])
	if val.Err() != nil {
		slog.Warn("Failed to load CUE schema, using defaults", "error", val.Err())
		return &Loader{ctx: ctx}
	}

	return &Loader{
		ctx:  ctx,
		base: val,
	}
}

// LoadFromFile loads configuration from a CUE file
func (l *Loader) LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return l.LoadFromBytes(data)
}

// LoadFromBytes loads configuration from CUE bytes
func (l *Loader) LoadFromBytes(data []byte) (*Config, error) {
	val := l.ctx.CompileBytes(data)
	if val.Err() != nil {
		return nil, fmt.Errorf("failed to parse CUE: %w", val.Err())
	}

	if l.base.Exists() {
		val = l.base.Unify(val)
		if err := val.Validate(cue.Concrete(true)); err != nil {
			return nil, fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	var config Config
	if err := val.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	if config.ProgressInterval == 0 {
		config.ProgressInterval = 30
	}

	return &config, nil
}

// LoadFromJSON loads configuration from a JSON file with CUE validation
func (l *Loader) LoadFromJSON(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON config: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	val := l.ctx.Encode(raw)
	if val.Err() != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", val.Err())
	}

	if l.base.Exists() {
		val = l.base.Unify(val)
		if err := val.Validate(cue.Concrete(true)); err != nil {
			return nil, fmt.Errorf("JSON validation failed: %w", err)
		}
	}

	var config Config
	if err := val.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	if config.ProgressInterval == 0 {
		config.ProgressInterval = 30
	}

	return &config, nil
}

// LoadFromEnvironment loads configuration from environment variables
func (l *Loader) LoadFromEnvironment() (*Config, error) {
	envMap := make(map[string]interface{})

	if val := os.Getenv("LF_INPUT_PATH"); val != "" {
		envMap["input_path"] = val
	}
	if val := os.Getenv("LF_OUTPUT_PATH"); val != "" {
		envMap["output_path"] = val
	}
	if val := os.Getenv("LF_WORKERS"); val != "" {
		var workers int
		if _, err := fmt.Sscanf(val, "%d", &workers); err == nil {
			envMap["workers"] = workers
		}
	}
	if val := os.Getenv("LF_BUFFER_SIZE_KB"); val != "" {
		var size int
		if _, err := fmt.Sscanf(val, "%d", &size); err == nil {
			envMap["buffer_size_kb"] = size
		}
	}
	if val := os.Getenv("LF_CHUNK_SIZE_MB"); val != "" {
		var size int
		if _, err := fmt.Sscanf(val, "%d", &size); err == nil {
			envMap["chunk_size_mb"] = size
		}
	}
	if val := os.Getenv("LF_USE_MMAP"); val != "" {
		envMap["use_mmap"] = strings.ToLower(val) == "true" || val == "1"
	}
	if val := os.Getenv("LF_USE_DIRECT_IO"); val != "" {
		envMap["use_direct_io"] = strings.ToLower(val) == "true" || val == "1"
	}
	if val := os.Getenv("LF_VERIFY_CHECKSUM"); val != "" {
		envMap["verify_checksum"] = strings.ToLower(val) == "true" || val == "1"
	}
	if val := os.Getenv("LF_MEMORY_LIMIT_GB"); val != "" {
		var limit int
		if _, err := fmt.Sscanf(val, "%d", &limit); err == nil {
			envMap["memory_limit_gb"] = limit
		}
	}
	if val := os.Getenv("LF_LOG_LEVEL"); val != "" {
		envMap["log_level"] = val
	}
	if val := os.Getenv("LF_LOG_FORMAT"); val != "" {
		envMap["log_format"] = val
	}
	if val := os.Getenv("LF_VERBOSE"); val != "" {
		envMap["enable_verbose"] = strings.ToLower(val) == "true" || val == "1"
	}
	if val := os.Getenv("LF_CHECKPOINT_DIR"); val != "" {
		envMap["checkpoint_dir"] = val
	}
	if val := os.Getenv("LF_PROGRESS_INTERVAL"); val != "" {
		var interval int
		if _, err := fmt.Sscanf(val, "%d", &interval); err == nil {
			envMap["progress_interval"] = interval
		}
	}

	if len(envMap) == 0 {
		return nil, nil
	}

	val := l.ctx.Encode(envMap)
	if val.Err() != nil {
		return nil, fmt.Errorf("failed to encode environment: %w", val.Err())
	}

	if l.base.Exists() {
		val = l.base.Unify(val)
		if err := val.Validate(cue.Concrete(true)); err != nil {
			return nil, fmt.Errorf("environment validation failed: %w", err)
		}
	}

	var config Config
	if err := val.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode environment: %w", err)
	}

	if config.ProgressInterval == 0 {
		config.ProgressInterval = 30
	}

	return &config, nil
}

// Merge merges multiple configurations (later takes precedence)
// Only non-zero values from later configs override earlier ones
func (l *Loader) Merge(configs ...*Config) (*Config, error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("no configurations to merge")
	}

	// Start with a copy of the first config
	merged := &Config{}
	*merged = *configs[0]

	// Merge each subsequent config
	for i := 1; i < len(configs); i++ {
		cfg := configs[i]

		// Only override if the value is non-zero (or specifically set)
		if cfg.InputPath != "" {
			merged.InputPath = cfg.InputPath
		}
		if cfg.OutputPath != "" {
			merged.OutputPath = cfg.OutputPath
		}
		if cfg.Workers > 0 {
			merged.Workers = cfg.Workers
		}
		if cfg.BufferSizeKB > 0 {
			merged.BufferSizeKB = cfg.BufferSizeKB
		}
		if cfg.ChunkSizeMB > 0 {
			merged.ChunkSizeMB = cfg.ChunkSizeMB
		}
		if cfg.MaxWorkers > 0 {
			merged.MaxWorkers = cfg.MaxWorkers
		}
		if cfg.PageSize > 0 {
			merged.PageSize = cfg.PageSize
		}
		if cfg.UseMmap {
			merged.UseMmap = cfg.UseMmap
		}
		if cfg.UseDirectIO {
			merged.UseDirectIO = cfg.UseDirectIO
		}
		if cfg.UseAsyncIO {
			merged.UseAsyncIO = cfg.UseAsyncIO
		}
		if cfg.VerifyChecksum {
			merged.VerifyChecksum = cfg.VerifyChecksum
		}
		if cfg.ChecksumAlgorithm != "" && cfg.ChecksumAlgorithm != "sha256" {
			merged.ChecksumAlgorithm = cfg.ChecksumAlgorithm
		}
		if cfg.GCPercent > 0 {
			merged.GCPercent = cfg.GCPercent
		}
		if cfg.MemoryLimitGB > 0 {
			merged.MemoryLimitGB = cfg.MemoryLimitGB
		}
		if cfg.LogLevel != "" {
			merged.LogLevel = cfg.LogLevel
		}
		if cfg.LogFormat != "" {
			merged.LogFormat = cfg.LogFormat
		}
		if cfg.LogFilePath != "" {
			merged.LogFilePath = cfg.LogFilePath
		}
		if cfg.ProgressInterval > 0 {
			merged.ProgressInterval = cfg.ProgressInterval
		}
		if cfg.EnableVerbose {
			merged.EnableVerbose = cfg.EnableVerbose
		}
		if cfg.EnableCheckpoint {
			merged.EnableCheckpoint = cfg.EnableCheckpoint
		}
		if cfg.CheckpointDir != "" {
			merged.CheckpointDir = cfg.CheckpointDir
		}
		if cfg.CheckpointInterval > 0 {
			merged.CheckpointInterval = cfg.CheckpointInterval
		}
		if cfg.DirectIOAlignment > 0 {
			merged.DirectIOAlignment = cfg.DirectIOAlignment
		}
		if cfg.QueueSize > 0 {
			merged.QueueSize = cfg.QueueSize
		}
		if cfg.Timeout > 0 {
			merged.Timeout = cfg.Timeout
		}
	}

	return merged, nil
}

// Validate validates a configuration against the CUE schema
func (l *Loader) Validate(cfg *Config) error {
	val := l.ctx.Encode(cfg)
	if val.Err() != nil {
		return fmt.Errorf("failed to encode config: %w", val.Err())
	}

	if l.base.Exists() {
		val = l.base.Unify(val)
		if err := val.Validate(cue.Concrete(true)); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}

// GenerateProfile generates a configuration profile based on detected hardware
func (l *Loader) GenerateProfile() *Config {
	storageType := detectStorageType()
	switch storageType {
	case "nvme":
		return NVMeProfile()
	case "hdd":
		return HDDProfile()
	case "network":
		return NetworkProfile()
	default:
		return DefaultConfig()
	}
}

// ToJSON exports configuration to JSON
func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

// Save saves configuration to a file
func (c *Config) Save(path string) error {
	data, err := c.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// detectStorageType attempts to detect the storage type
func detectStorageType() string {
	if _, err := os.Stat("/proc/mounts"); err == nil {
		data, _ := os.ReadFile("/proc/mounts")
		if strings.Contains(string(data), "nfs") || strings.Contains(string(data), "cifs") {
			return "network"
		}
	}

	if _, err := os.Stat("/dev/nvme0"); err == nil {
		return "nvme"
	}

	return "hdd"
}

// Helper to check if a value is zero
func isZero(v interface{}) bool {
	return reflect.ValueOf(v).IsZero()
}
