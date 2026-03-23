// Package config defines the configuration schema for large file operations
package config

// Config defines the complete configuration schema
#Config: {
	// File paths - required fields
	input_path:  string & !=""
	output_path: string & !=""
	
	// Performance tuning
	workers:         int | *8
	buffer_size_kb:  int | *256
	chunk_size_mb:   int | *16
	max_workers:     int | *16
	page_size:       int | *4096
	
	// I/O strategies
	use_mmap:      bool | *false
	use_direct_io: bool | *true
	use_async_io:  bool | *false
	
	// Data integrity
	verify_checksum:     bool | *true
	checksum_algorithm:  string | *"sha256"
	
	// Memory management
	gc_percent:      int | *500
	memory_limit_gb: int | *8
	
	// Logging
	log_level:         string | *"info"
	log_format:        string | *"json"
	log_file_path:     string | *""
	progress_interval: int | *30  // seconds
	enable_verbose:    bool | *false
	
	// Recovery
	enable_checkpoint:   bool | *true
	checkpoint_dir:      string | *"."
	checkpoint_interval: int | *100
	
	// Advanced
	direct_io_alignment: int | *4096
	queue_size:          int | *100
	timeout:             int | *0  // seconds
}

// Production-optimized configuration
#ProductionConfig: #Config & {
	workers: 16
	buffer_size_kb: 1024
	chunk_size_mb: 64
	memory_limit_gb: 32
	enable_verbose: true
	use_direct_io: true
	verify_checksum: true
	queue_size: 200
	log_level: "info"
	progress_interval: 10
}

// Development configuration
#DevelopmentConfig: #Config & {
	workers: 4
	buffer_size_kb: 128
	chunk_size_mb: 8
	memory_limit_gb: 4
	enable_verbose: true
	verify_checksum: false
	enable_checkpoint: false
	gc_percent: 200
	log_level: "debug"
	log_format: "text"
	progress_interval: 30
}

// NVMe optimized configuration
#NVMeProfile: #Config & {
	workers: 16
	buffer_size_kb: 1024
	chunk_size_mb: 64
	use_direct_io: true
	queue_size: 200
	direct_io_alignment: 4096
	use_async_io: true
	progress_interval: 30
}

// HDD optimized configuration
#HDDProfile: #Config & {
	workers: 4
	buffer_size_kb: 256
	chunk_size_mb: 8
	use_direct_io: false
	queue_size: 50
	progress_interval: 30
}

// Network storage optimized configuration
#NetworkProfile: #Config & {
	workers: 8
	buffer_size_kb: 512
	chunk_size_mb: 32
	use_mmap: true
	use_direct_io: false
	queue_size: 100
	verify_checksum: true
	progress_interval: 30
}

// Low-memory configuration
#LowMemoryProfile: #Config & {
	workers: 2
	buffer_size_kb: 64
	chunk_size_mb: 4
	memory_limit_gb: 2
	gc_percent: 300
	queue_size: 25
	enable_checkpoint: true
	checkpoint_interval: 50
	progress_interval: 30
}