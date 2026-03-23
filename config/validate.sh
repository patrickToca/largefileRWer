#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Validating CUE Configuration${NC}"

# Check if cue is installed
if ! command -v cue &> /dev/null; then
    echo -e "${RED}Error: cue is not installed${NC}"
    exit 1
fi

cd "$(dirname "$0")/.."

echo -n "1. Validating config.cue syntax... "
cue vet ./config/config.cue && echo -e "${GREEN}✅${NC}" || { echo -e "${RED}❌${NC}"; exit 1; }

echo -n "2. Testing Config profile... "
cue eval ./config/config.cue -e '#Config' > /dev/null 2>&1 && echo -e "${GREEN}✅${NC}" || { echo -e "${RED}❌${NC}"; exit 1; }

echo -n "3. Testing ProductionConfig... "
cue eval ./config/config.cue -e '#ProductionConfig' > /dev/null 2>&1 && echo -e "${GREEN}✅${NC}" || { echo -e "${RED}❌${NC}"; exit 1; }

echo -n "4. Testing DevelopmentConfig... "
cue eval ./config/config.cue -e '#DevelopmentConfig' > /dev/null 2>&1 && echo -e "${GREEN}✅${NC}" || { echo -e "${RED}❌${NC}"; exit 1; }

echo -n "5. Testing NVMeProfile... "
cue eval ./config/config.cue -e '#NVMeProfile' > /dev/null 2>&1 && echo -e "${GREEN}✅${NC}" || { echo -e "${RED}❌${NC}"; exit 1; }

echo -n "6. Testing HDDProfile... "
cue eval ./config/config.cue -e '#HDDProfile' > /dev/null 2>&1 && echo -e "${GREEN}✅${NC}" || { echo -e "${RED}❌${NC}"; exit 1; }

echo -n "7. Testing NetworkProfile... "
cue eval ./config/config.cue -e '#NetworkProfile' > /dev/null 2>&1 && echo -e "${GREEN}✅${NC}" || { echo -e "${RED}❌${NC}"; exit 1; }

echo -n "8. Testing LowMemoryProfile... "
cue eval ./config/config.cue -e '#LowMemoryProfile' > /dev/null 2>&1 && echo -e "${GREEN}✅${NC}" || { echo -e "${RED}❌${NC}"; exit 1; }

echo -e "\n${GREEN}✅ All CUE validations passed!${NC}"