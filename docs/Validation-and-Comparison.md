# Validation & Comparison

tmq provides tools to validate TOML files and compare differences between them.

## Validation

### Basic Syntax Validation
```bash
# Check if TOML file is valid
tmq --validate config.toml

# Returns exit code 0 if valid, 1 if invalid
echo $?  # Check exit code
```

### Batch Validation
```bash
# Validate multiple files
tmq --validate config/*.toml

# Use with find
find . -name "*.toml" -exec tmq --validate {} \;
```

### Validation in Scripts
```bash
#!/bin/bash
validate_toml() {
    local file="$1"
    if tmq --validate "$file" >/dev/null 2>&1; then
        echo "✓ $file is valid"
        return 0
    else
        echo "✗ $file is invalid"
        return 1
    fi
}

# Validate all TOML files in directory
for file in *.toml; do
    validate_toml "$file" || exit 1
done
```

## Comparison

### Basic File Comparison
```bash
# Compare two TOML files
tmq --compare config1.toml config2.toml

# Exit codes:
# 0 = files are identical
# 1 = files are different
```

### Detailed Comparison Output
```bash
# See detailed differences
tmq --compare old-config.toml new-config.toml
```

### Comparison in CI/CD
```bash
# Fail build if configuration changed unexpectedly
if ! tmq --compare expected.toml actual.toml >/dev/null; then
    echo "Configuration mismatch!"
    tmq --compare expected.toml actual.toml
    exit 1
fi
```

## Advanced Usage

### Validation with Error Details
```bash
# tmq shows detailed error messages
tmq --validate invalid.toml
# Error: parse error: expected newline at line 5, column 10
# DETAILS: Check your TOML syntax
# ACTION: Fix the syntax error and try again
```

### Comparison with Output
```bash
# Redirect comparison output to file
tmq --compare file1.toml file2.toml > differences.txt

# Use in scripts
if tmq --compare "$EXPECTED" "$ACTUAL" > diff.log; then
    echo "Files match"
else
    echo "Differences found:"
    cat diff.log
fi
```

## Integration Examples

### Pre-commit Hook
```bash
#!/bin/bash
# .git/hooks/pre-commit

# Validate all TOML files
echo "Validating TOML files..."
if ! find . -name "*.toml" -exec tmq --validate {} \;; then
    echo "TOML validation failed"
    exit 1
fi

echo "All TOML files are valid"
```

### Configuration Drift Detection
```bash
#!/bin/bash
# Check if production config matches expected config

PROD_CONFIG="prod-config.toml"
EXPECTED_CONFIG="expected-config.toml"

if tmq --compare "$EXPECTED_CONFIG" "$PROD_CONFIG" >/dev/null; then
    echo "✓ Production config matches expected configuration"
    exit 0
else
    echo "✗ Configuration drift detected!"
    echo "Differences:"
    tmq --compare "$EXPECTED_CONFIG" "$PROD_CONFIG"
    exit 1
fi
```

### CI Pipeline Integration
```yaml
# .github/workflows/validate.yml
name: Validate Configuration

on:
  pull_request:
    paths:
      - 'config/*.toml'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: checkout@v4

      - name: Setup tmq
        run: |
          wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
          chmod +x tmq-linux-amd64
          sudo mv tmq-linux-amd64 /usr/local/bin/tmq

      - name: Validate TOML files
        run: |
          for file in config/*.toml; do
            echo "Validating $file..."
            tmq --validate "$file"
          done

      - name: Check configuration consistency
        run: |
          tmq --compare config/base.toml config/production.toml
```

## Error Messages

### Validation Errors
```bash
# Syntax errors
tmq --validate malformed.toml
# ERROR: TOML parsing failed
# DETAILS: Expected newline at line 5, column 10
# ACTION: Check your TOML syntax and fix any formatting issues

# File not found
tmq --validate nonexistent.toml
# ERROR: File operation error
# DETAILS: Cannot read file 'nonexistent.toml'
# ACTION: Ensure the file exists and is readable
```

### Comparison Errors
```bash
# File not found
tmq --compare missing1.toml missing2.toml
# ERROR: File operation error
# DETAILS: Cannot read file 'missing1.toml'
# ACTION: Ensure both files exist and are readable

# Parse error in one file
tmq --compare valid.toml invalid.toml
# ERROR: TOML parsing failed
# DETAILS: Invalid TOML in 'invalid.toml'
# ACTION: Fix the TOML syntax in invalid.toml
```

## Output Formats

### Validation Output
```bash
# Successful validation (no output, exit code 0)
tmq --validate valid.toml

# Failed validation (error message, exit code 1)
tmq --validate invalid.toml
```

### Comparison Output
```bash
# Identical files (no output, exit code 0)
tmq --compare same1.toml same1.toml

# Different files (shows differences, exit code 1)
tmq --compare old.toml new.toml
# Files are different:
# - old.toml: version = "1.0.0"
# + new.toml: version = "2.0.0"
```

## Best Practices

### Validation in Development
```bash
# Add to your build process
make validate-config:
    @for file in config/*.toml; do \
        echo "Validating $$file..."; \
        tmq --validate "$$file" || exit 1; \
    done
    @echo "All configurations valid"
```

### Configuration Testing
```bash
# Test configuration changes
#!/bin/bash

CONFIG_FILE="app-config.toml"
BACKUP_FILE="${CONFIG_FILE}.backup"

# Backup current config
cp "$CONFIG_FILE" "$BACKUP_FILE"

# Make test changes
tmq '.debug = true' -i "$CONFIG_FILE"
tmq '.test_mode = true' -i "$CONFIG_FILE"

# Validate changes
if tmq --validate "$CONFIG_FILE" >/dev/null; then
    echo "✓ Configuration changes are valid"
    # Test your application here
    # ./run-tests.sh
else
    echo "✗ Invalid configuration after changes"
fi

# Restore backup
mv "$BACKUP_FILE" "$CONFIG_FILE"
```

### Version Control Integration
```bash
# Check for configuration changes in commits
git diff --name-only HEAD~1 | grep '\.toml$' | while read -r file; do
    echo "Checking $file..."
    if ! tmq --validate "$file" >/dev/null; then
        echo "Invalid TOML in $file"
        exit 1
    fi
done
```

### Monitoring Configuration Changes
```bash
#!/bin/bash
# Monitor configuration for unauthorized changes

BASELINE="config-baseline.toml"
CURRENT="config.toml"

if [ ! -f "$BASELINE" ]; then
    echo "Creating baseline configuration..."
    cp "$CURRENT" "$BASELINE"
    exit 0
fi

if ! tmq --compare "$BASELINE" "$CURRENT" >/dev/null; then
    echo "WARNING: Configuration has changed!"
    echo "Differences:"
    tmq --compare "$BASELINE" "$CURRENT"

    read -p "Accept changes and update baseline? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cp "$CURRENT" "$BASELINE"
        echo "Baseline updated"
    fi
fi
```

## Performance

- Validation is fast - typically < 100ms for most files
- Comparison speed depends on file size and differences
- Memory usage is minimal for both operations

## Exit Codes Summary

| Operation | Success | Error |
|-----------|---------|-------|
| Validation | 0 | 1 (parse error) |
| Comparison | 0 (identical) | 1 (different) |
| File errors | - | 4 (file not found/readable) |