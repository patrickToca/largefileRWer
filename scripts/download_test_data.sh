#!/bin/bash
# scripts/download_test_data.sh

set -e

echo "Downloading test data..."

# Option 1: NYC Taxi Data (small sample - 100MB)
if [ ! -f "test_data/taxi_sample.csv" ]; then
    mkdir -p test_data
    echo "Downloading NYC Taxi sample data..."
    curl -L -o test_data/taxi_sample.csv \
        "https://data.cityofnewyork.us/api/views/5rqd-h5ci/rows.csv?accessType=DOWNLOAD&bom=true"
    echo "✅ Downloaded $(du -h test_data/taxi_sample.csv | cut -f1) file"
fi

# Option 2: Create synthetic data if needed
if [ ! -f "test_data/test_1gb.csv" ]; then
    echo "Creating 1GB synthetic test file..."
    ./scripts/create_test_file.sh test_data/test_1gb.csv 1
fi

echo ""
echo "Test files ready:"
ls -lh test_data/