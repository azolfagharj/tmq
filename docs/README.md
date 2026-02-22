# tmq (TOML Query)

**tmq** is a lightweight and flexible command-line TOML processor. Like jq for JSON, yq for YAML — tmq is for TOML.

tmq lets you query, modify, validate, and compare TOML files with simple dot-notation syntax. Slice nested tables, update values in-place, batch process configs, or convert to JSON and YAML — all from the command line with a single binary.

tmq is written in Go: zero dependencies, one executable. Download for Linux, macOS, or Windows and run it anywhere.

```
tmq = TOML + Query
```

## Quick Links

!!! note "🚀 Getting Started"
    Get up and running with tmq in minutes

    - [Installation](Installation.md) - Download and setup
    - [Quick Start](Getting-Started.md) - Basic usage and examples

!!! info "📖 User Guide"
    Learn how to query, modify, validate, and process TOML files

    - [Query Operations](Query-Operations.md) - Read data with dot notation
    - [Modification Operations](Modification-Operations.md) - Set and delete values
    - [Validation & Comparison](Validation-and-Comparison.md) - Check and compare files
    - [Bulk Operations](Bulk-Operations.md) - Process multiple files at once

!!! tip "📄 Reference"
    Complete command reference and practical examples

    - [Command Reference](Command-Reference.md) - Full CLI documentation
    - [Examples](Examples.md) - Real-world usage examples

!!! warning "❓ Troubleshooting"
    Solutions to common issues and how to get help

    - [Troubleshooting](Troubleshooting.md) - Common issues and solutions

## Features Overview

### ✅ Implemented Features
- **Query & Read**: Access nested TOML values with dot notation
- **Modify & Write**: Set, delete, and update TOML values in-place
- **Validation**: Check TOML file syntax and structure
- **Comparison**: Compare two TOML files for differences
- **Bulk Operations**: Process multiple files at once
- **Output Formats**: JSON, YAML, and TOML output
- **Script-Friendly**: Clear exit codes and error messages
- **Cross-Platform**: Linux, macOS, Windows binaries

### 🚧 Planned Features
- JSON/YAML → TOML conversion
- Comment preservation
- Library API
- Plugin system
- Performance optimizations

## Getting Help

- **Issues**: Report bugs or request features on [GitHub Issues](https://github.com/azolfagharj/tmq/issues)
- **Discussions**: Join community discussions on [GitHub Discussions](https://github.com/azolfagharj/tmq/discussions)

## Contributing

Contributions welcome! See the main repository README for development setup.