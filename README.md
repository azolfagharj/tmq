# tmq — TOML Query Tool

[![CI](https://github.com/azolfagharj/tmq/actions/workflows/ci.yml/badge.svg)](https://github.com/azolfagharj/tmq/actions) [![Documentation](https://img.shields.io/badge/Documentation-Wiki-blue?logo=github)](https://github.com/azolfagharj/tmq/wiki) [![Donate](https://img.shields.io/badge/Donate-to%20Keep%20This%20Project%20Alive-orange)](https://azolfagharj.github.io/donate/)

**Complete standalone TOML CLI processor .** supporting query, modification, and format conversion

```
tmq = ToMl + Query
```

A fast, script-friendly command-line tool for querying, modifying, and converting TOML files. Works with pipes, supports bulk operations, and provides clear exit codes for automation.

## 📚 Documentation

📖 **[Complete Documentation & Wiki](https://github.com/azolfagharj/tmq/wiki)** - Installation, usage examples, troubleshooting, and command reference

## Installation

### Binary Release (Recommended)

Download the latest release binary for your system from the [GitHub Releases](https://github.com/azolfagharj/tmq/releases) page.

**Available binaries:**

* `tmq-linux-amd64` — Linux (64-bit)
* `tmq-linux-arm64` — Linux (ARM64)
* `tmq-darwin-amd64` — macOS (Intel)
* `tmq-darwin-arm64` — macOS (Apple Silicon)
* `tmq-windows-amd64.exe` — Windows (64-bit)

**Quick setup:**

1. Download the binary for your system:
   ```bash
   # Linux (amd64)
   wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
   chmod +x tmq-linux-amd64
   mv tmq-linux-amd64 tmq

   # Linux (ARM64)
   wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-arm64
   chmod +x tmq-linux-arm64
   mv tmq-linux-arm64 tmq

   # macOS (Apple Silicon)
   wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-darwin-arm64
   chmod +x tmq-darwin-arm64
   mv tmq-darwin-arm64 tmq

   # macOS (Intel)
   wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-darwin-amd64
   chmod +x tmq-darwin-amd64
   mv tmq-darwin-amd64 tmq

   # Windows (amd64)
   # Download: https://github.com/azolfagharj/tmq/releases/latest/download/tmq-windows-amd64.exe
   # Rename it to tmq.exe
   ```
2. Move to PATH (optional):
   ```bash
   sudo mv tmq /usr/local/bin/
   ```

### Build from Source

If you prefer to build from source:

**Prerequisites:**

* Go 1.23 or later

**Build steps:**

```bash
git clone https://github.com/azolfagharj/tmq.git
cd tmq
go build -o bin/tmq ./cmd/tmq
```

## Quick Start

```bash
# Read specific value
tmq '.project.version' pyproject.toml

# Read from stdin
cat pyproject.toml | tmq '.project.name'

# Convert to JSON
tmq '.' config.toml -o json

# Convert to YAML
tmq '.' config.toml -o yaml

# Display all data
tmq config.toml

# Version info
tmq --version

# Help
tmq --help
```

## Features

### Query & Read
- **Query syntax**: Access nested TOML values with dot notation
- **Stdin support**: Pipe TOML data from other commands
- **Multiple files**: Process multiple TOML files in one command
- **Output formats**: JSON, YAML, or TOML output

### Modify & Write
- **In-place editing**: Modify TOML files directly (`-i` flag)
- **Set values**: Update existing keys or create new ones
- **Delete keys**: Remove keys from TOML files
- **Dry-run mode**: Preview changes without modifying files (`--dry-run`)

### Validation & Comparison
- **Syntax validation**: Check TOML file validity (`--validate`)
- **File comparison**: Compare two TOML files for differences (`--compare`)
- **Bulk operations**: Process multiple files at once

### Script-Friendly
- **Clear exit codes**: 0 (success), 1 (parse error), 2 (usage error), 3 (security error), 4 (file error)
- **Structured error output**: Machine-readable error messages
- **Pipe support**: Full stdin/stdout support for scripting
- **No dependencies**: Single binary, no external requirements

## Usage Examples

### Basic Queries
```bash
# Get project name
tmq '.project.name' pyproject.toml

# Get nested configuration
tmq '.database.host' config.toml

# Get array element
tmq '.servers[0].name' config.toml
```

### Modifications
```bash
# Set a value in-place
tmq '.project.version = "2.0.0"' -i pyproject.toml

# Delete a key
tmq 'del(.optional_dependency)' -i file.toml

# Dry-run to preview changes
tmq '.version = "3.0.0"' --dry-run config.toml
```

### Validation & Comparison
```bash
# Validate TOML syntax
tmq --validate config.toml

# Compare two files
tmq --compare config1.toml config2.toml

# Bulk validation
tmq --validate config/*.toml
```

### Bulk Operations
```bash
# Query multiple files
tmq '.version' config/*.toml

# Bulk update
tmq '.version = "3.0.0"' -i config/*.toml
```

### Output Formats
```bash
# JSON output
tmq '.' config.toml -o json

# YAML output
tmq '.database' config.toml -o yaml

# TOML output (default)
tmq '.' config.toml
```

### Scripting Examples
```bash
# Error handling in scripts
VERSION=$(tmq '.project.version' pyproject.toml) || exit 1

# Chain with other tools
tmq '.' config.toml | jq '.database.host'

# Process with find
find . -name "*.toml" -exec tmq '.version' {} \;
```

## Query Syntax

tmq uses a simple dot notation for accessing TOML data:

```toml
[project]
name = "my-app"
version = "1.0.0"

[database]
host = "localhost"
port = 5432

[[servers]]
name = "server1"
ip = "192.168.1.1"
```

```bash
# Access project name
tmq '.project.name' config.toml
# Output: "my-app"

# Access database port
tmq '.database.port' config.toml
# Output: 5432

# Access array element
tmq '.servers[0].name' config.toml
# Output: "server1"
```

## Exit Codes

tmq uses clear exit codes for automation:

- `0` — Success
- `1` — TOML parsing or runtime error
- `2` — Invalid arguments or usage error
- `3` — Security violation (path traversal, etc.)
- `4` — File operation error

## Requirements

- **OS**: Linux, macOS, Windows
- **Architecture**: amd64, arm64
- **No external dependencies** — single binary

## Performance

- **Fast execution**: < 100ms for typical operations
- **Low memory usage**: < 10MB peak
- **Single binary**: No startup overhead


## License

MIT License

---
## Support this Project



 🤝 **Enjoying this free project?** <a href="https://azolfagharj.github.io/donate/">Consider supporting</a> its development

<a href="https://azolfagharj.github.io/donate/"><img src="https://img.shields.io/badge/Donate-Support%20Development-orange?style=for-the-badge" alt="Donate"></a>

---
