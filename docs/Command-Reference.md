# Command Reference

Complete reference for all tmq commands, flags, and options.

## Synopsis

```bash
tmq [OPTIONS] [QUERY] [FILE...]
tmq [OPTIONS] --validate [FILE...]
tmq [OPTIONS] --compare FILE1 FILE2
```

## Global Options

### Output Options
- `-o, --output FORMAT`: Output format (`toml`, `json`, `yaml`)
  - Default: `toml`
  - Example: `tmq '.data' config.toml -o json`

### Modification Options
- `-i, --inplace`: Modify files in-place
  - Must be used with set/delete operations
  - Example: `tmq '.version = "2.0"' -i config.toml`

### Dry Run
- `--dry-run`: Preview changes without modifying files
  - Shows what would be done
  - Exit code 0 for success, 1 for errors
  - Example: `tmq '.version = "2.0"' --dry-run config.toml`

### Validation
- `--validate`: Validate TOML syntax
  - Exit code 0 if valid, 1 if invalid
  - Can process multiple files
  - Example: `tmq --validate config.toml`

### Comparison
- `--compare FILE1 FILE2`: Compare two TOML files
  - Exit code 0 if identical, 1 if different
  - Shows detailed differences
  - Example: `tmq --compare old.toml new.toml`

### Information
- `-v, --version`: Show version information
- `-h, --help`: Show help text

## Query Syntax

### Basic Queries
```bash
# Root key
.key

# Nested table
.table.key

# Deep nesting
.app.database.host

# Array element
.array[0]

# Array element with nesting
.servers[1].name
```

### Set Operations
```bash
# Simple assignment
.key = "value"

# Number assignment
.port = 8080

# Boolean assignment
.enabled = true

# Array assignment
.tags = ["web", "api"]

# Object assignment
.config = { host = "localhost", port = 5432 }
```

### Delete Operations
```bash
# Delete root key
del(.key)

# Delete nested key
del(.table.key)

# Delete array element
del(.array[0])
```

## Examples

### Query Examples
```bash
# Simple query
tmq '.version' config.toml

# Nested query
tmq '.database.host' config.toml

# Array query
tmq '.servers[0].name' config.toml

# Multiple files
tmq '.version' config/*.toml

# JSON output
tmq '.config' file.toml -o json

# YAML output
tmq '.config' file.toml -o yaml
```

### Modification Examples
```bash
# Set string value
tmq '.version = "2.0.0"' -i config.toml

# Set nested value
tmq '.database.host = "prod-db"' -i config.toml

# Set array
tmq '.ports = [8080, 8443]' -i config.toml

# Delete key
tmq 'del(.obsolete)' -i config.toml

# Dry run
tmq '.version = "test"' --dry-run config.toml
```

### Validation Examples
```bash
# Single file
tmq --validate config.toml

# Multiple files
tmq --validate *.toml

# In script
if tmq --validate config.toml; then
    echo "Valid TOML"
fi
```

### Comparison Examples
```bash
# Compare files
tmq --compare config1.toml config2.toml

# In script
if tmq --compare expected.toml actual.toml; then
    echo "Files match"
fi
```

## Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | Operation completed successfully |
| 1 | Parse Error | TOML parsing failed or query error |
| 2 | Usage Error | Invalid command-line arguments |
| 3 | Security Error | Path traversal or security violation |
| 4 | File Error | File not found, permission denied, etc. |

## Error Messages

tmq provides structured error messages:

```
ERROR: <error_type>
DETAILS: <detailed_description>
ACTION: <suggested_fix>
```

### Common Errors

#### Parse Errors
```bash
tmq '.invalid..' config.toml
# ERROR: Invalid query path
# DETAILS: Query path cannot be empty
# ACTION: Check your query syntax
```

#### File Errors
```bash
tmq '.' nonexistent.toml
# ERROR: File operation error
# DETAILS: Cannot read file 'nonexistent.toml'
# ACTION: Ensure the file exists and is readable
```

#### Validation Errors
```bash
tmq --validate malformed.toml
# ERROR: TOML parsing failed
# DETAILS: Expected newline at line 5, column 10
# ACTION: Fix the TOML syntax error
```

## Environment Variables

tmq does not use environment variables for configuration. All options are specified via command-line flags.

## File Formats

### TOML Input
tmq accepts standard TOML 1.0.0 format:

```toml
# Comments are preserved in output
title = "Example"

[database]
host = "localhost"
port = 5432

[[servers]]
name = "web1"
ip = "192.168.1.1"
```

### JSON Output
```json
{
  "title": "Example",
  "database": {
    "host": "localhost",
    "port": 5432
  },
  "servers": [
    {
      "name": "web1",
      "ip": "192.168.1.1"
    }
  ]
}
```

### YAML Output
```yaml
title: Example
database:
  host: localhost
  port: 5432
servers:
- name: web1
  ip: 192.168.1.1
```

## Performance

### Benchmarks
- Query operations: O(1) - constant time
- File parsing: O(n) where n is file size
- Memory usage: < 10MB for typical files
- Large files: Scales linearly with size

### Optimization Tips
- Use specific queries instead of full file output
- Prefer JSON/YAML output for programmatic use
- Use bulk operations for multiple files

## Platform Support

### Supported Platforms
- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

### Binary Names
- `tmq-linux-amd64` (Linux Intel/AMD)
- `tmq-linux-arm64` (Linux ARM)
- `tmq-darwin-amd64` (macOS Intel)
- `tmq-darwin-arm64` (macOS Apple Silicon)
- `tmq-windows-amd64.exe` (Windows)

## Limitations

### Current Limitations
- No JSON/YAML → TOML conversion (planned)
- No comment preservation in modifications (planned)
- No plugin system (planned)
- No library API (planned)

### File Size Limits
- Maximum file size: 100MB
- Recommended: < 10MB for optimal performance

### Path Length Limits
- Maximum path length: 1024 characters
- Maximum nesting depth: 100 levels

## Troubleshooting

### Common Issues

#### "command not found"
```bash
# Check if tmq is in PATH
which tmq

# Or use full path
/path/to/tmq --version
```

#### "permission denied"
```bash
# Make executable
chmod +x tmq

# Or run with full permissions
sudo ./tmq --version
```

#### "directory not found"
```bash
# Check file exists
ls -la config.toml

# Check current directory
pwd

# Use absolute path
tmq '.' /full/path/to/config.toml
```

#### Invalid query
```bash
# Check syntax
tmq --help

# Test with simple query
tmq '.' config.toml
```

### Debug Mode
```bash
# Use verbose output (when available)
tmq --help

# Check file contents
cat config.toml

# Validate file
tmq --validate config.toml
```

### Getting Help
```bash
# Show help
tmq --help

# Show version
tmq --version

# Report issues: https://github.com/azolfagharj/tmq/issues
```

## Version History

### v1.0.1
- Initial release
- Query, set, delete operations
- Validation and comparison
- Bulk operations
- JSON/YAML output
- Cross-platform binaries