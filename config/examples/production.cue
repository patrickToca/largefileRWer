package config

#ProductionConfig & {
	input_path:  "/data/largefile.dat"
	output_path: "/data/output.dat"
	workers:     16
	buffer_size_kb: 1024
	chunk_size_mb:  64
	verify_checksum: true
	enable_checkpoint: true
	checkpoint_dir: "/data/checkpoints"
	enable_verbose: true
	log_level: "info"
	log_format: "json"
}