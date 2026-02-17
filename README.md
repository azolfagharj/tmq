# tmq — TOML Query

**Complete standalone TOML CLI tool.** Like jq for JSON, yq for YAML — but for TOML.

```
tmq = ToMl + Query
```

## Installation

*Coming soon.* Single binary for Linux, macOS, Windows.

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

# Error handling in scripts
VERSION=$(tmq '.project.version' pyproject.toml) || exit 1

# Set value in-place
tmq '.project.version = "2.0"' -i pyproject.toml

# Delete a key
tmq 'del(.optional_dependency)' -i file.toml

# Validate TOML syntax
tmq --validate config.toml

# Compare two TOML files
tmq --compare config1.toml config2.toml

# Process multiple files
tmq '.version' config/*.toml

# Bulk validation
tmq --validate config/*.toml

# Bulk update
tmq '.version = "3.0.0"' -i config/*.toml
```

## Features

| Feature | Status |
|---------|--------|
| Query (read) | ✅ Implemented |
| Set / delete | ✅ Implemented |
| In-place edit | ✅ Implemented |
| TOML → JSON | ✅ Implemented |
| TOML → YAML | ✅ Implemented |
| JSON/YAML → TOML | ⏳ Planned |
| Pipe, stdin, stdout | ✅ Implemented |
| Validation mode | ✅ Implemented |
| Comparison mode | ✅ Implemented |
| Bulk operations | ✅ Implemented |
| Comment preservation | ⏳ Planned |
| Library API | ⏳ Planned |
| Plugin system | ⏳ Planned |
| Performance (< 100ms) | ⏳ Planned |
| Memory (< 10MB) | ⏳ Planned |
| Cross-platform | ⏳ Planned |

## Requirements

- Go 1.21+

## Development

This project follows strict development standards with 100% test coverage (quality-focused) and comprehensive documentation.

See [.cursorrules](.cursorrules) for detailed development guidelines.

## License

MIT (or to be decided)

## Contributing

See `protectedocs/docs/` for project vision and roadmap. Contributions welcome once the foundation is in place.
