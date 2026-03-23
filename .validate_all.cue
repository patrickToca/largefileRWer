// Import the config definitions
import "pmjtoca/largefileRWer/config"

// Test all valid profiles
_valid_config: config.#Config & {input_path: "a", output_path: "b"}
_valid_production: config.#ProductionConfig & {input_path: "a", output_path: "b"}
_valid_development: config.#DevelopmentConfig & {input_path: "a", output_path: "b"}
_valid_nvme: config.#NVMeProfile & {input_path: "a", output_path: "b"}
_valid_hdd: config.#HDDProfile & {input_path: "a", output_path: "b"}
_valid_network: config.#NetworkProfile & {input_path: "a", output_path: "b"}
_valid_lowmem: config.#LowMemoryProfile & {input_path: "a", output_path: "b"}

// Test custom configuration
_custom: config.#Config & {
    input_path: "/custom/input.dat"
    output_path: "/custom/output.dat"
    workers: 12
    buffer_size_kb: 512
    chunk_size_mb: 32
    verify_checksum: true
    enable_verbose: true
    log_level: "debug"
}
