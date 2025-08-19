# Sonic Performance Issue with Streaming Compressed JSON Data

## Issue Summary

This benchmark demonstrates a significant performance degradation in Sonic when processing large compressed JSON data using streaming APIs. The issue is particularly pronounced when decompressing and parsing JSON data simultaneously through streaming interfaces.

## Technical Details

### Problem Description
When using streaming APIs to process gzip-compressed JSON data:
- Sonic's performance degrades by ~90%+ compared to raw JSON processing
- Standard library shows moderate degradation (~30-40%) under the same conditions
- This creates a scenario where standard library outperforms Sonic for compressed data

### Test Configuration
- **Data Structure**: Tree with 10,000 nodes (~9.2 MB JSON)
- **Compression**: gzip default level
- **Streaming Method**: `gzip.Reader` directly with `sonic.ConfigDefault.NewDecoder()`
- **Sonic Config**: `sonic.ConfigDefault.NewDecoder()`

## Relevant Code Example

The problematic pattern in `benchmark_test.go` (BenchmarkSonicDecodingCompressed):
```go
var tree TreeNode
// Use a bytes.Reader instead of bytes.Buffer for better performance
reader := bytes.NewReader(compressedData.Bytes())
gzipReader, err := gzip.NewReader(reader)
if err != nil {
    b.Fatal(err)
}
decoder := sonic.ConfigDefault.NewDecoder(gzipReader)
err = decoder.Decode(&tree)
if err != nil {
    b.Fatal(err)
}
// Close the gzip reader
err = gzipReader.Close()
if err != nil {
    b.Fatal(err)
}
```

## Expected Behavior vs Actual Performance

**Expected**: Sonic should maintain reasonable performance advantage even with compressed data streaming.

**Actual**: Sonic performance drops below standard library for compressed data streaming scenarios.

## Reproduction Steps

1. Generate test data: `go run main.go`
2. Run specific benchmark: `go test -bench=BenchmarkSonicDecodingCompressed -benchmem`
3. Compare with standard library: `go test -bench=BenchmarkStdlibDecodingCompressed -benchmem`

## Files

- `benchmark_test.go` - Contains the streaming compressed JSON benchmarks that demonstrate the issue
- `test_data.json` - Pre-generated 9.2MB JSON test data
- `main.go` - Data generation utility (not needed for benchmark reproduction)

## Credits

AI-generated benchmark suite designed to help identify and resolve Sonic performance issues with streaming compressed JSON data.
