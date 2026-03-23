#!/bin/bash
# config/validate.sh

set -e

echo "Validating CUE configuration..."

# Check if cue is installed
if ! command -v cue &> /dev/null; then
    echo "Error: cue is not installed"
    echo "Install with: go install cuelang.org/go/cmd/cue@latest"
    exit 1
fi

# Validate basic syntax
echo "1. Validating syntax..."
cue vet ./config.cue 2>&1
echo "✅ Syntax valid"

# Test each profile
echo ""
echo "2. Testing configuration profiles..."

test_profile() {
    local profile=$1
    local test_file="test_${profile}.cue"
    
    cat > "$test_file" <<EOF
package config

#${profile} & {
    input_path: "/test/input.dat"
    output_path: "/test/output.dat"
}
EOF
    
    echo -n "Testing $profile... "
    if cue vet "$test_file" 2>/dev/null; then
        echo "✅"
        rm "$test_file"
        return 0
    else
        echo "❌"
        cue vet "$test_file" 2>&1
        rm "$test_file"
        return 1
    fi
}

# Test all profiles
test_profile "Config" || exit 1
test_profile "ProductionConfig" || exit 1
test_profile "DevelopmentConfig" || exit 1
test_profile "NVMeProfile" || exit 1
test_profile "HDDProfile" || exit 1
test_profile "NetworkProfile" || exit 1
test_profile "LowMemoryProfile" || exit 1

echo ""
echo "✅ All CUE validation tests passed!"