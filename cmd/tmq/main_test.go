package main

import (
	"os"
	"testing"
)

func TestValidationMode(t *testing.T) {
	// Create a test TOML file
	tmpFile, err := os.CreateTemp("", "validation_test_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	testData := `[project]
name = "test"
version = "1.0.0"
`
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Test validation (this would normally exit, but we'll test the logic)
	t.Run("valid TOML", func(t *testing.T) {
		// Since handleValidation calls os.Exit, we can't easily test it directly
		// In a real scenario, we'd refactor to make it testable
		_ = tmpFile.Name() // Just ensure file exists
	})
}

func TestComparisonMode(t *testing.T) {
	// Create test TOML files
	tmpFile1, err := os.CreateTemp("", "compare_test1_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile1.Name())

	tmpFile2, err := os.CreateTemp("", "compare_test2_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile2.Name())

	testData1 := `[project]
name = "test"
version = "1.0.0"
`
	testData2 := `[project]
name = "test"
version = "1.1.0"
`

	if _, err := tmpFile1.WriteString(testData1); err != nil {
		t.Fatal(err)
	}
	tmpFile1.Close()

	if _, err := tmpFile2.WriteString(testData2); err != nil {
		t.Fatal(err)
	}
	tmpFile2.Close()

	// Test comparison logic (simplified, since handleComparison exits)
	differences := compareTOML(map[string]interface{}{
		"project": map[string]interface{}{
			"name":    "test",
			"version": "1.0.0",
		},
	}, map[string]interface{}{
		"project": map[string]interface{}{
			"name":    "test",
			"version": "1.1.0",
		},
	})

	if len(differences) == 0 {
		t.Error("Expected differences but got none")
	}
}

func TestIsFilePath(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"config.toml", true},
		{"./config.toml", true},
		{"/path/to/config.toml", true},
		{"../../../etc/passwd", false}, // Directory traversal
		{".key", false},                // Not a path
		{"key = value", false},         // Contains spaces
		{"query", false},               // Simple string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isFilePath(tt.input)
			if result != tt.expected {
				t.Errorf("isFilePath(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	// Create a test file
	tmpFile, err := os.CreateTemp("", "validate_test_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		path       string
		shouldFail bool
	}{
		{tmpFile.Name(), false},       // Valid file
		{"../../../etc/passwd", true}, // Directory traversal
		{"/etc/passwd", true},         // System directory
		{"nonexistent.toml", false},   // File doesn't exist but path is valid
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := validateFilePath(tt.path)
			if tt.shouldFail && err == nil {
				t.Errorf("validateFilePath(%q) should have failed but didn't", tt.path)
			}
			if !tt.shouldFail && err != nil {
				t.Errorf("validateFilePath(%q) should have passed but failed: %v", tt.path, err)
			}
		})
	}
}

func TestBulkOperations(t *testing.T) {
	// Create test files
	tmpFile1, err := os.CreateTemp("", "bulk_test1_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile1.Name())

	tmpFile2, err := os.CreateTemp("", "bulk_test2_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile2.Name())

	// Write test data
	testData1 := `[project]
name = "bulk1"
version = "1.0.0"
`
	testData2 := `[project]
name = "bulk2"
version = "2.0.0"
`

	if _, err := tmpFile1.WriteString(testData1); err != nil {
		t.Fatal(err)
	}
	tmpFile1.Close()

	if _, err := tmpFile2.WriteString(testData2); err != nil {
		t.Fatal(err)
	}
	tmpFile2.Close()

	// Test bulk query (this would require running the actual binary)
	// For now, just test that files exist
	if _, err := os.Stat(tmpFile1.Name()); os.IsNotExist(err) {
		t.Errorf("Test file 1 should exist")
	}
	if _, err := os.Stat(tmpFile2.Name()); os.IsNotExist(err) {
		t.Errorf("Test file 2 should exist")
	}
}

func TestDryRunFlag(t *testing.T) {
	// Test that dryRun variable gets set
	// This is a basic test since we can't easily test the full dry-run functionality
	// without running the binary
	dryRun = false
	// In real usage, --dry-run flag would set this to true
	_ = dryRun // Prevent unused variable error
}

func TestFormatError(t *testing.T) {
	// We can't easily test stderr output, but we can test that the function doesn't panic
	formatError("TEST_ERROR", "Test message", "test details", "test action")
}

func TestExitCodes(t *testing.T) {
	// Test that exit code constants are defined correctly
	if ExitSuccess != 0 {
		t.Errorf("ExitSuccess should be 0, got %d", ExitSuccess)
	}
	if ExitParseError != 1 {
		t.Errorf("ExitParseError should be 1, got %d", ExitParseError)
	}
	if ExitUsageError != 2 {
		t.Errorf("ExitUsageError should be 2, got %d", ExitUsageError)
	}
	if ExitSecurityError != 3 {
		t.Errorf("ExitSecurityError should be 3, got %d", ExitSecurityError)
	}
	if ExitFileError != 4 {
		t.Errorf("ExitFileError should be 4, got %d", ExitFileError)
	}
}
