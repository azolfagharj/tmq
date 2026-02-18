# Examples

Comprehensive collection of tmq usage examples for common scenarios.

## Basic Configuration Management

### Application Configuration
```toml
# app.toml
[app]
name = "myapp"
version = "1.0.0"
debug = true

[database]
host = "localhost"
port = 5432
name = "myapp_db"

[logging]
level = "info"
file = "/var/log/myapp.log"
```

```bash
# Get application info
tmq '.app.name' app.toml          # "myapp"
tmq '.app.version' app.toml       # "1.0.0"

# Database configuration
tmq '.database' app.toml
# Output: host = "localhost"
#         port = 5432
#         name = "myapp_db"

# Update for production
tmq '.app.debug = false' -i app.toml
tmq '.database.host = "prod-db.example.com"' -i app.toml
```

### Environment-Specific Configs
```toml
# config.toml
environment = "development"

[database]
host = "localhost"
port = 5432

[features]
debug = true
metrics = false
```

```bash
# Switch to production
tmq '.environment = "production"' -i config.toml
tmq '.database.host = "prod-db.cluster.com"' -i config.toml
tmq '.database.port = 3306' -i config.toml
tmq '.features.debug = false' -i config.toml
tmq '.features.metrics = true' -i config.toml

# Verify changes
tmq '.' config.toml
```

## Web Server Configuration

### Nginx-style Config
```toml
# nginx.toml
[server]
listen = 80
server_name = "example.com"
root = "/var/www/html"

[ssl]
enabled = false
certificate = "/etc/ssl/certs/example.crt"
key = "/etc/ssl/private/example.key"

[[locations]]
path = "/"
try_files = ["$uri", "$uri/", "/index.html"]

[[locations]]
path = "/api"
proxy_pass = "http://localhost:3000"
```

```bash
# Server configuration
tmq '.server' nginx.toml

# SSL setup
tmq '.ssl.enabled = true' -i nginx.toml
tmq '.ssl.certificate = "/etc/letsencrypt/live/example.com/fullchain.pem"' -i nginx.toml
tmq '.ssl.key = "/etc/letsencrypt/live/example.com/privkey.pem"' -i nginx.toml

# Add API location
tmq '.locations = [{ path = "/api", proxy_pass = "http://api:8080" }]' -i nginx.toml
```

### Database Configuration

```toml
# database.toml
[connection]
host = "localhost"
port = 5432
database = "myapp"
username = "app_user"
ssl_mode = "require"

[pool]
min_connections = 2
max_connections = 20
idle_timeout = 300

[migrations]
enabled = true
path = "./migrations"
```

```bash
# Connection string
HOST=$(tmq '.connection.host' database.toml)
PORT=$(tmq '.connection.port' database.toml)
DB=$(tmq '.connection.database' database.toml)
USER=$(tmq '.connection.username' database.toml)

echo "postgresql://$USER@$HOST:$PORT/$DB"

# Pool configuration
tmq '.pool' database.toml -o json

# Enable migrations
tmq '.migrations.enabled = true' -i database.toml
```

## CI/CD Pipeline Examples

### GitHub Actions Config
```toml
# .github/config.toml
[build]
go_version = "1.23"
platforms = ["linux/amd64", "linux/arm64", "darwin/amd64", "darwin/arm64"]

[test]
race_detector = true
coverage = true
threshold = 80

[release]
draft = false
prerelease = false
generate_release_notes = true
```

```bash
# Update Go version across all workflows
find .github/workflows -name "*.yml" -exec \
    tmq '.build.go_version = "1.24"' -i {} \;

# Update test coverage threshold
tmq '.test.threshold = 85' -i .github/config.toml

# Enable race detector
tmq '.test.race_detector = true' -i .github/config.toml
```

### Docker Compose Override
```toml
# docker-compose.override.toml
[services.app]
image = "myapp:latest"
ports = ["8080:8080"]
environment = ["NODE_ENV=development"]

[services.db]
image = "postgres:13"
environment = ["POSTGRES_DB=myapp_dev"]
volumes = ["./data:/var/lib/postgresql/data"]
```

```bash
# Switch to production images
tmq '.services.app.image = "myapp:v1.2.0"' -i docker-compose.override.toml
tmq '.services.db.image = "postgres:15"' -i docker-compose.override.toml

# Update environment
tmq '.services.app.environment = ["NODE_ENV=production"]' -i docker-compose.override.toml
```

## Package Management

### Package Configuration
```toml
# package.toml
[package]
name = "tmq"
version = "1.0.0"
description = "TOML Query Tool"
authors = ["Your Name <your@email.com>"]

[dependencies]
toml = "1.0.0"
clap = "4.0"

[build]
target = "x86_64-unknown-linux-gnu"
release = true
```

```bash
# Version management
CURRENT_VERSION=$(tmq '.package.version' package.toml)
echo "Current version: $CURRENT_VERSION"

# Update version
NEW_VERSION="1.0.1"
tmq ".package.version = \"$NEW_VERSION\"" -i package.toml

# Update dependencies
tmq '.dependencies.toml = "1.1.0"' -i package.toml
```

### Cargo.toml Example
```toml
# Cargo.toml
[package]
name = "my-crate"
version = "0.1.0"
edition = "2021"

[dependencies]
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1.0", features = ["full"] }

[profile.release]
opt-level = 3
lto = true
```

```bash
# Dependency updates
tmq '.dependencies.serde.version = "1.1"' -i Cargo.toml
tmq '.dependencies.tokio.version = "1.1"' -i Cargo.toml

# Profile optimization
tmq '.profile.release.opt-level = 3' -i Cargo.toml
tmq '.profile.release.lto = true' -i Cargo.toml
```

## Configuration Validation

### Schema Validation
```bash
#!/bin/bash
# Validate configuration against schema

validate_config() {
    local file="$1"

    # Check required fields
    required_fields=(".app.name" ".database.host" ".server.port")
    for field in "${required_fields[@]}"; do
        if ! tmq "$field" "$file" >/dev/null 2>&1; then
            echo "Missing required field: $field"
            return 1
        fi
    done

    # Check port range
    port=$(tmq '.server.port' "$file")
    if [ "$port" -lt 1024 ] || [ "$port" -gt 65535 ]; then
        echo "Invalid port: $port"
        return 1
    fi

    return 0
}

# Validate all configs
for config in config/*.toml; do
    if validate_config "$config"; then
        echo "✓ $config is valid"
    else
        echo "✗ $config is invalid"
        exit 1
    fi
done
```

### Configuration Migration
```bash
#!/bin/bash
# Migrate old configuration format to new format

migrate_config() {
    local file="$1"

    # Migrate old field names
    if tmq '.old_database_host' "$file" >/dev/null 2>&1; then
        host=$(tmq '.old_database_host' "$file")
        tmq ".database.host = $host" -i "$file"
        tmq 'del(.old_database_host)' -i "$file"
    fi

    # Add new required fields with defaults
    tmq '.database.port = 5432' -i "$file"
    tmq '.features.logging = true' -i "$file"
    tmq '.version = "1.0.0"' -i "$file"
}

# Migrate all configurations
for config in config/*.toml; do
    echo "Migrating $config..."
    migrate_config "$config"
done

# Validate migrated configs
tmq --validate config/*.toml
```

## Monitoring and Alerting

### Health Check Configuration
```toml
# health.toml
[checks]
[checks.database]
enabled = true
query = "SELECT 1"
timeout = 5

[checks.redis]
enabled = true
key = "health"
timeout = 2

[alerts]
[alerts.email]
enabled = true
recipients = ["admin@example.com"]

[alerts.slack]
enabled = false
webhook = "https://hooks.slack.com/..."
```

```bash
# Enable health checks
tmq '.checks.database.enabled = true' -i health.toml
tmq '.checks.redis.enabled = true' -i health.toml

# Configure alerts
tmq '.alerts.email.enabled = true' -i health.toml
tmq '.alerts.slack.enabled = false' -i health.toml

# Update timeouts
tmq '.checks.database.timeout = 10' -i health.toml
tmq '.checks.redis.timeout = 3' -i health.toml
```

### Log Configuration
```toml
# logging.toml
[logger]
level = "info"
format = "json"

[outputs]
[outputs.console]
enabled = true
level = "debug"

[outputs.file]
enabled = true
path = "/var/log/app.log"
max_size = "10MB"
max_age = "30d"

[outputs.syslog]
enabled = false
facility = "user"
```

```bash
# Change log level
tmq '.logger.level = "debug"' -i logging.toml

# Enable file logging
tmq '.outputs.file.enabled = true' -i logging.toml
tmq '.outputs.file.path = "/var/log/myapp.log"' -i logging.toml

# Configure log rotation
tmq '.outputs.file.max_size = "50MB"' -i logging.toml
tmq '.outputs.file.max_age = "7d"' -i logging.toml
```

## Multi-Environment Management

### Environment Overrides
```bash
#!/bin/bash
# Manage multiple environments

ENVIRONMENTS=("development" "staging" "production")
BASE_CONFIG="config/base.toml"

for env in "${ENVIRONMENTS[@]}"; do
    config_file="config/$env.toml"

    # Copy base config
    cp "$BASE_CONFIG" "$config_file"

    # Apply environment-specific overrides
    case $env in
        "development")
            tmq '.debug = true' -i "$config_file"
            tmq '.database.host = "localhost"' -i "$config_file"
            tmq '.logging.level = "debug"' -i "$config_file"
            ;;
        "staging")
            tmq '.debug = false' -i "$config_file"
            tmq '.database.host = "staging-db.example.com"' -i "$config_file"
            tmq '.logging.level = "info"' -i "$config_file"
            ;;
        "production")
            tmq '.debug = false' -i "$config_file"
            tmq '.database.host = "prod-db-cluster.example.com"' -i "$config_file"
            tmq '.logging.level = "warn"' -i "$config_file"
            ;;
    esac

    echo "Generated $config_file"
done

# Validate all generated configs
tmq --validate config/*.toml
```

### Feature Flags
```toml
# features.toml
[features]
new_ui = false
api_v2 = false
experimental = false
beta_features = false

[rollout]
percentage = 0
user_groups = []
```

```bash
# Enable features for beta users
tmq '.features.new_ui = true' -i features.toml
tmq '.features.api_v2 = true' -i features.toml

# Configure rollout
tmq '.rollout.percentage = 25' -i features.toml
tmq '.rollout.user_groups = ["beta", "premium"]' -i features.toml

# Gradual rollout
tmq '.rollout.percentage = 50' -i features.toml
tmq '.rollout.percentage = 75' -i features.toml
tmq '.rollout.percentage = 100' -i features.toml
```

## Backup and Recovery

### Configuration Backup
```bash
#!/bin/bash
# Backup configuration with versioning

backup_config() {
    local source_dir="$1"
    local backup_dir="$2"

    timestamp=$(date +%Y%m%d_%H%M%S)
    backup_path="$backup_dir/backup_$timestamp"

    mkdir -p "$backup_path"

    # Copy all configs
    cp "$source_dir"/*.toml "$backup_path/"

    # Create manifest
    cat > "$backup_path/manifest.txt" << EOF
Backup created: $(date)
Source: $source_dir
Files: $(ls -1 "$backup_path"/*.toml | wc -l)

Configuration versions:
$(for file in "$backup_path"/*.toml; do
    version=$(tmq '.version // "unknown"' "$file" 2>/dev/null || echo "unknown")
    echo "$(basename "$file"): $version"
done)
EOF

    echo "Backup created: $backup_path"
}

# Usage
backup_config "config" "backups"
```

### Configuration Restore
```bash
#!/bin/bash
# Restore configuration from backup

restore_config() {
    local backup_path="$1"
    local target_dir="$2"

    if [ ! -d "$backup_path" ]; then
        echo "Backup not found: $backup_path"
        return 1
    fi

    echo "Restoring from $backup_path..."

    # Validate backup files
    if ! tmq --validate "$backup_path"/*.toml >/dev/null 2>&1; then
        echo "Backup contains invalid files"
        return 1
    fi

    # Create restore point
    restore_backup="$target_dir/restore_backup_$(date +%s)"
    mkdir -p "$restore_backup"
    cp "$target_dir"/*.toml "$restore_backup/" 2>/dev/null || true

    # Restore files
    cp "$backup_path"/*.toml "$target_dir/"

    # Validate restored files
    if tmq --validate "$target_dir"/*.toml >/dev/null 2>&1; then
        echo "✓ Restore successful"
        echo "Previous state backed up to: $restore_backup"
        return 0
    else
        echo "✗ Restore failed - restoring previous state"
        cp "$restore_backup"/*.toml "$target_dir/" 2>/dev/null || true
        return 1
    fi
}

# Usage
restore_config "backups/backup_20240115_143000" "config"
```

## Advanced Scripting

### Dynamic Configuration Generation
```bash
#!/bin/bash
# Generate configuration based on environment variables

generate_config() {
    local output_file="$1"

    cat > "$output_file" << EOF
[app]
name = "$(tmq '.app.name' base.toml)"
version = "$(tmq '.app.version' base.toml)"
environment = "${ENVIRONMENT:-development}"

[database]
host = "${DB_HOST:-localhost}"
port = ${DB_PORT:-5432}
name = "${DB_NAME:-myapp}"

[features]
debug = ${DEBUG:-false}
metrics = ${METRICS:-true}
EOF

    # Validate generated config
    if tmq --validate "$output_file" >/dev/null 2>&1; then
        echo "Generated valid config: $output_file"
    else
        echo "Generated invalid config"
        return 1
    fi
}

generate_config "config/generated.toml"
```

### Configuration Diff and Patch
```bash
#!/bin/bash
# Show configuration changes between versions

config_diff() {
    local old_file="$1"
    local new_file="$2"

    echo "Configuration changes:"
    echo "======================"

    # Compare versions
    old_version=$(tmq '.version // "unknown"' "$old_file")
    new_version=$(tmq '.version // "unknown"' "$new_file")

    if [ "$old_version" != "$new_version" ]; then
        echo "Version: $old_version → $new_version"
    fi

    # Compare key settings
    keys_to_check=(".database.host" ".database.port" ".debug" ".logging.level")

    for key in "${keys_to_check[@]}"; do
        old_value=$(tmq "$key" "$old_file" 2>/dev/null || echo "not set")
        new_value=$(tmq "$key" "$new_file" 2>/dev/null || echo "not set")

        if [ "$old_value" != "$new_value" ]; then
            echo "$key: $old_value → $new_value"
        fi
    done
}

config_diff "config/v1.0.toml" "config/v1.1.toml"
```