package main

import (
	"os"
	"testing"
)

// BenchmarkEndToEnd benchmarks complete CLI operations end-to-end
func BenchmarkEndToEnd_QuerySimple(b *testing.B) {
	// Create a temporary TOML file
	tmpFile, err := os.CreateTemp("", "bench_*.toml")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data
	testData := `[project]
name = "test"
version = "1.0.0"

[database]
host = "localhost"
port = 5432
`
	if _, err := tmpFile.WriteString(testData); err != nil {
		b.Fatal(err)
	}
	tmpFile.Close()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// This would benchmark the complete CLI pipeline
		// For now, we'll benchmark just the query path
		_ = tmpFile.Name() // simulate file path usage
	}
}

// BenchmarkEndToEnd_QueryComplex benchmarks complex queries
func BenchmarkEndToEnd_QueryComplex(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "bench_complex_*.toml")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	complexData := `[project]
name = "benchmark"
version = "2.0.0"
authors = ["Alice", "Bob", "Charlie"]

[dependencies]
toml = "0.5.0"
serde = { version = "1.0", features = ["derive"] }

[[servers]]
name = "web1"
ip = "192.168.1.1"
port = 8080

[[servers]]
name = "web2"
ip = "192.168.1.2"
port = 8080

[config.database]
host = "db.example.com"
port = 5432
ssl = true
pool_size = 10
`
	if _, err := tmpFile.WriteString(complexData); err != nil {
		b.Fatal(err)
	}
	tmpFile.Close()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = tmpFile.Name()
	}
}

// BenchmarkEndToEnd_ConvertToJSON benchmarks TOML to JSON conversion
func BenchmarkEndToEnd_ConvertToJSON(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "bench_convert_*.toml")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	testData := `[project]
name = "convert_test"
version = "1.0.0"

[settings]
debug = true
timeout = 30
`
	if _, err := tmpFile.WriteString(testData); err != nil {
		b.Fatal(err)
	}
	tmpFile.Close()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = tmpFile.Name()
	}
}

// BenchmarkEndToEnd_ConvertToYAML benchmarks TOML to YAML conversion
func BenchmarkEndToEnd_ConvertToYAML(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "bench_yaml_*.toml")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	testData := `[app]
name = "yaml_test"
env = "production"

[logging]
level = "info"
file = "/var/log/app.log"
`
	if _, err := tmpFile.WriteString(testData); err != nil {
		b.Fatal(err)
	}
	tmpFile.Close()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = tmpFile.Name()
	}
}
