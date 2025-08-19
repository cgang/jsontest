#!/bin/bash

echo "JSON Encoding/Decoding Benchmark"
echo "================================="

# Generate test data if it doesn't exist
if [ ! -f "test_data.json" ]; then
    echo "Generating test data..."
    go run main.go
    echo ""
fi

# Run benchmarks
echo "Running benchmarks..."
go test -bench=.
