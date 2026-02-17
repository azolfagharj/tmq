# Getting Started

This guide covers the basics of using tmq to work with TOML files.

## Basic Concepts

### TOML Structure
TOML files use a simple, human-readable format:

```toml
# config.toml
title = "My Application"

[database]
host = "localhost"
port = 5432
enabled = true

[server]
host = "0.0.0.0"
port = 8080

[[users]]
name = "Alice"
role = "admin"

[[users]]
name = "Bob"
role = "user"
```

### Query Syntax
tmq uses dot notation to access TOML values:
- `title` → root-level keys
- `database.host` → nested table values
- `users[0].name` → array element access

## First Steps

### Check Installation
```bash
tmq --version
tmq --help
```

### Basic Query
Create a test TOML file:
```bash
cat > config.toml << 'EOF'
title = "My App"
version = "1.0.0"

[database]
host = "localhost"
port = 5432
EOF
```

Query values:
```bash
# Get title
tmq '.title' config.toml
# Output: "My App"

# Get version
tmq '.version' config.toml
# Output: "1.0.0"

# Get database host
tmq '.database.host' config.toml
# Output: "localhost"

# Get database port
tmq '.database.port' config.toml
# Output: 5432
```

### Display All Data
```bash
# Show entire file
tmq '.' config.toml

# Format as JSON
tmq '.' config.toml -o json

# Format as YAML
tmq '.' config.toml -o yaml
```

## Common Patterns

### Stdin Input
```bash
# Read from stdin
cat config.toml | tmq '.version'

# Use with other tools
echo 'version = "2.0.0"' | tmq '.version'
```

### Error Handling in Scripts
```bash
#!/bin/bash
VERSION=$(tmq '.version' config.toml)
if [ $? -ne 0 ]; then
    echo "Error reading version from config.toml"
    exit 1
fi
echo "Version: $VERSION"
```

### File Operations
```bash
# Check if file exists and is valid TOML
if tmq '.' config.toml > /dev/null 2>&1; then
    echo "config.toml is valid"
else
    echo "config.toml is invalid or missing"
fi
```

## Next Steps

Now that you know the basics:

1. **Query Operations**: Learn advanced querying in [Query Operations](Query-Operations.md)
2. **Modifications**: Learn to modify TOML files in [Modification Operations](Modification-Operations.md)
3. **Validation**: Check file validity in [Validation & Comparison](Validation-and-Comparison.md)
4. **Examples**: See comprehensive examples in [Examples](Examples.md)

## Quick Reference

| Command | Description |
|---------|-------------|
| `tmq '.key' file.toml` | Query a value |
| `tmq '.' file.toml` | Show entire file |
| `tmq '.' file.toml -o json` | Output as JSON |
| `cat file.toml \| tmq '.key'` | Read from stdin |
| `tmq --validate file.toml` | Validate TOML syntax |
| `tmq --help` | Show help |

## Exit Codes

tmq uses standard exit codes:
- `0`: Success
- `1`: Parse/runtime error
- `2`: Usage error
- `3`: Security error
- `4`: File error