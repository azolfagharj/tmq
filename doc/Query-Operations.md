# Query Operations

tmq provides powerful querying capabilities to extract data from TOML files using a simple dot notation syntax.

## Basic Queries

### Root Level Keys
```toml
# config.toml
title = "My App"
version = "1.0.0"
enabled = true
```

```bash
tmq '.title' config.toml     # "My App"
tmq '.version' config.toml   # "1.0.0"
tmq '.enabled' config.toml   # true
```

### Nested Tables
```toml
[database]
host = "localhost"
port = 5432
ssl = false

[server]
host = "0.0.0.0"
port = 8080
```

```bash
tmq '.database.host' config.toml    # "localhost"
tmq '.database.port' config.toml    # 5432
tmq '.server.host' config.toml      # "0.0.0.0"
```

### Deep Nesting
```toml
[app]
[app.config]
[app.config.database]
host = "db.example.com"
port = 3306

[app.config.cache]
type = "redis"
ttl = 3600
```

```bash
tmq '.app.config.database.host' config.toml   # "db.example.com"
tmq '.app.config.cache.type' config.toml      # "redis"
tmq '.app.config.cache.ttl' config.toml       # 3600
```

## Array Operations

### Access Array Elements
```toml
[[servers]]
name = "web1"
ip = "192.168.1.1"

[[servers]]
name = "web2"
ip = "192.168.1.2"

[[servers]]
name = "db1"
ip = "192.168.1.10"
```

```bash
# First server
tmq '.servers[0].name' config.toml    # "web1"
tmq '.servers[0].ip' config.toml      # "192.168.1.1"

# Second server
tmq '.servers[1].name' config.toml    # "web2"

# Database server
tmq '.servers[2].name' config.toml    # "db1"
```

### Array of Values
```toml
ports = [8080, 8443, 9000]
tags = ["web", "api", "admin"]
```

```bash
# Access entire arrays
tmq '.ports' config.toml      # [8080, 8443, 9000]
tmq '.tags' config.toml       # ["web", "api", "admin"]

# Access array elements
tmq '.ports[0]' config.toml   # 8080
tmq '.ports[1]' config.toml   # 8443
tmq '.tags[2]' config.toml    # "admin"
```

## Output Formats

### Default TOML Output
```bash
tmq '.database' config.toml
# Output: host = "localhost"
#         port = 5432
```

### JSON Output
```bash
tmq '.database' config.toml -o json
# Output: {"host":"localhost","port":5432}
```

### YAML Output
```bash
tmq '.database' config.toml -o yaml
# Output: host: localhost
#         port: 5432
```

## Advanced Queries

### Complex Structures
```toml
[app]
name = "myapp"

[app.database]
host = "db.example.com"
credentials = { username = "admin", password = "secret" }

[app.features]
logging = true
metrics = { enabled = true, port = 9090 }
```

```bash
# Access nested objects
tmq '.app.database.credentials' config.toml
# Output: username = "admin"
#         password = "secret"

tmq '.app.database.credentials.username' config.toml    # "admin"
tmq '.app.features.metrics' config.toml
# Output: enabled = true
#         port = 9090
```

### Mixed Data Types
```toml
# config.toml
version = "1.2.3"
debug = true
timeout = 30
pi = 3.14159

[metadata]
created = 2024-01-15T10:30:00Z
tags = ["production", "stable"]
```

```bash
tmq '.version' config.toml     # "1.2.3"
tmq '.debug' config.toml       # true
tmq '.timeout' config.toml     # 30
tmq '.pi' config.toml          # 3.14159
tmq '.metadata.tags' config.toml    # ["production", "stable"]
```

## Query Root

### Access Everything
```bash
# Show entire file
tmq '.' config.toml

# Same as above (explicit root)
tmq '. .' config.toml
```

### Root with Different Formats
```bash
# JSON output of entire file
tmq '.' config.toml -o json

# YAML output of entire file
tmq '.' config.toml -o yaml
```

## Error Handling

### Non-existent Keys
```bash
tmq '.nonexistent' config.toml
# Error: key 'nonexistent' not found
# Exit code: 1
```

### Invalid Paths
```bash
tmq '.invalid..path' config.toml
# Error: invalid query path
# Exit code: 1
```

### Type Mismatches
```bash
# Trying to access array index on a non-array
tmq '.title[0]' config.toml
# Error: cannot index into string
# Exit code: 1
```

## Performance Notes

- Queries are evaluated in constant time O(1)
- Memory usage scales with the size of the queried data
- Large files (>100MB) may require increased memory limits

## Best Practices

### Scripting
```bash
#!/bin/bash
# Safe querying with error handling
DB_HOST=$(tmq '.database.host' config.toml 2>/dev/null) || {
    echo "Error: Could not read database host from config"
    exit 1
}
echo "Database host: $DB_HOST"
```

### Validation Before Query
```bash
# Check if file is valid before querying
if tmq '.' config.toml >/dev/null 2>&1; then
    VERSION=$(tmq '.version' config.toml)
    echo "Version: $VERSION"
else
    echo "Invalid TOML file"
    exit 1
fi
```

### Use Appropriate Output Format
```bash
# For scripts, use raw output (default)
HOST=$(tmq '.database.host' config.toml)

# For data processing, use JSON
tmq '.servers' config.toml -o json | jq '.[0].name'
```