#!/bin/bash
# scripts/create_test_file.sh - Generate large test CSV files

set -e

FILE=${1:-test.csv}
SIZE_GB=${2:-1}
USE_RANDOM=${3:-false}

echo "Creating ${SIZE_GB}GB test file: $FILE"

# Calculate rows needed for ~1GB
# Rough estimate: ~100 bytes per row
ROWS=$((SIZE_GB * 1024 * 1024 * 1024 / 100))

# Create header
echo "id,timestamp,value,category,description" > "$FILE"

# Create data in chunks
CHUNK_SIZE=1000000
CHUNKS=$((ROWS / CHUNK_SIZE))

echo "Generating $ROWS rows in $CHUNKS chunks..."

for ((i=0; i<$CHUNKS; i++)); do
    echo -n "Chunk $((i+1))/$CHUNKS... "
    
    # Generate data chunk
    for ((j=0; j<$CHUNK_SIZE; j++)); do
        ID=$((i * CHUNK_SIZE + j + 1))
        TIMESTAMP=$(date -u +"%Y-%m-%d %H:%M:%S")
        VALUE=$((RANDOM % 10000))
        
        if [ "$USE_RANDOM" = "true" ]; then
            CATEGORY=$((RANDOM % 10))
            DESC="random_text_$RANDOM"
        else
            CATEGORY=$((ID % 10))
            DESC="data_point_$ID"
        fi
        
        echo "$ID,$TIMESTAMP,$VALUE,$CATEGORY,$DESC" >> "$FILE"
    done
    
    echo "done"
done

echo "✅ Created $(du -h $FILE | cut -f1) file: $FILE"