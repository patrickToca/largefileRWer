package config

valid_Config: #Config & {
    input_path: "/test/input.dat"
    output_path: "/test/output.dat"
}

valid_Production: #ProductionConfig & {
    input_path: "/test/input.dat"
    output_path: "/test/output.dat"
}

valid_Development: #DevelopmentConfig & {
    input_path: "/test/input.dat"
    output_path: "/test/output.dat"
}

valid_NVMe: #NVMeProfile & {
    input_path: "/test/input.dat"
    output_path: "/test/output.dat"
}

valid_HDD: #HDDProfile & {
    input_path: "/test/input.dat"
    output_path: "/test/output.dat"
}

valid_Network: #NetworkProfile & {
    input_path: "/test/input.dat"
    output_path: "/test/output.dat"
}

valid_LowMemory: #LowMemoryProfile & {
    input_path: "/test/input.dat"
    output_path: "/test/output.dat"
}
