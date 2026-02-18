# Bulk Operations

tmq supports processing multiple TOML files simultaneously, making it ideal for batch processing, configuration management, and automation tasks.

## Processing Multiple Files

### Query Multiple Files
```bash
# Get version from all config files
tmq '.version' config/*.toml

# Output format:
# config/app.toml: "1.0.0"
# config/database.toml: "2.1.0"
# config/server.toml: "1.5.0"
```

### Bulk Validation
```bash
# Validate all TOML files in directory
tmq --validate config/*.toml

# Validate recursively
find . -name "*.toml" -exec tmq --validate {} \;
```

### Bulk Modification
```bash
# Update version in all config files
tmq '.version = "2.0.0"' -i config/*.toml

# This updates all matching files in-place
```

## Output Formatting

### Query Results
```bash
# Default output (filename: value)
tmq '.database.host' config/*.toml
# config/app.toml: "localhost"
# config/prod.toml: "prod-db.example.com"

# JSON output
tmq '.database.host' config/*.toml -o json
# {"config/app.toml": "localhost", "config/prod.toml": "prod-db.example.com"}

# Custom formatting with scripts
tmq '.version' config/*.toml | while IFS=: read -r file version; do
    echo "File: $file, Version: $version"
done
```

## Error Handling

### Continue on Error
```bash
# Process all files even if some fail
tmq '.version' config/*.toml || true

# Check individual exit codes
for file in config/*.toml; do
    if tmq '.version' "$file" >/dev/null 2>&1; then
        echo "$file: OK"
    else
        echo "$file: FAILED"
    fi
done
```

### Error Aggregation
```bash
#!/bin/bash
# Script to handle bulk operations with error reporting

process_files() {
    local operation="$1"
    local files="$2"
    local errors=()

    for file in $files; do
        if ! tmq "$operation" "$file" >/dev/null 2>&1; then
            errors+=("$file")
        fi
    done

    if [ ${#errors[@]} -gt 0 ]; then
        echo "Failed files:"
        printf '  %s\n' "${errors[@]}"
        return 1
    fi
}

# Validate all configs
if ! process_files "--validate" "config/*.toml"; then
    echo "Some files failed validation"
    exit 1
fi

echo "All files validated successfully"
```

## Advanced Bulk Operations

### Conditional Processing
```bash
# Only process files with specific criteria
for file in config/*.toml; do
    # Check if file has a specific key
    if tmq '.production' "$file" >/dev/null 2>&1; then
        echo "Processing production config: $file"
        tmq '.version = "stable"' -i "$file"
    fi
done
```

### Batch Updates with Backup
```bash
#!/bin/bash
# Safe bulk update with backup

BACKUP_DIR="config-backup-$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

# Create backups
for file in config/*.toml; do
    cp "$file" "$BACKUP_DIR/"
done

echo "Backups created in $BACKUP_DIR"

# Perform bulk update
tmq '.updated_at = "$(date)"' -i config/*.toml

# Validate all files
if tmq --validate config/*.toml >/dev/null 2>&1; then
    echo "✓ All files updated and validated successfully"
else
    echo "✗ Some files failed validation, restoring backups..."
    for file in "$BACKUP_DIR"/*.toml; do
        basename_file=$(basename "$file")
        cp "$file" "config/$basename_file"
    done
    echo "Backups restored"
    exit 1
fi
```

### Configuration Migration
```bash
#!/bin/bash
# Migrate configuration structure across multiple files

echo "Starting configuration migration..."

# Rename keys
for file in config/*.toml; do
    echo "Migrating $file..."

    # Rename database.host to database.hostname
    host=$(tmq '.database.host' "$file" 2>/dev/null)
    if [ -n "$host" ]; then
        tmq ".database.hostname = $host" -i "$file"
        tmq 'del(.database.host)' -i "$file"
    fi

    # Add new default values
    tmq '.database.port = 5432' -i "$file"
    tmq '.features.logging = true' -i "$file"
done

# Validate all migrated files
if tmq --validate config/*.toml >/dev/null 2>&1; then
    echo "✓ Migration completed successfully"
else
    echo "✗ Migration failed - check files manually"
    exit 1
fi
```

## Directory Operations

### Recursive Processing
```bash
# Process all TOML files in directory tree
find . -name "*.toml" -type f -exec tmq '.version' {} \;

# Process files in specific directories
tmq '.version' config/**/*.toml
```

### Organized Output
```bash
# Group output by directory
for dir in config/*/ ; do
    echo "=== $dir ==="
    tmq '.version' "$dir"*.toml
done
```

## Performance Optimization

### Parallel Processing
```bash
# Process files in parallel (basic)
for file in config/*.toml; do
    tmq '.version' "$file" &
done
wait

# Advanced parallel processing
#!/bin/bash
MAX_JOBS=4
process_file() {
    local file="$1"
    echo "Processing $file..."
    tmq '.version' "$file"
}

export -f process_file

# Use xargs for parallel execution
find config/ -name "*.toml" -print0 |
    xargs -0 -n1 -P"$MAX_JOBS" bash -c 'process_file "$@"' --
```

### Batch Size Management
```bash
# Process files in batches to avoid memory issues
BATCH_SIZE=10
files=(config/*.toml)
total_files=${#files[@]}

for ((i=0; i<total_files; i+=BATCH_SIZE)); do
    batch_files=("${files[@]:i:BATCH_SIZE}")
    echo "Processing batch ${i}/${total_files}..."

    # Process batch
    tmq '.version' "${batch_files[@]}" > "batch_$((i/BATCH_SIZE)).log"
done
```

## Integration Examples

### CI/CD Pipeline
```yaml
# .github/workflows/bulk-update.yml
name: Bulk Configuration Update

on:
  workflow_dispatch:
    inputs:
      new_version:
        description: 'New version to set'
        required: true

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: checkout@v4

      - name: Setup tmq
        run: |
          wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
          chmod +x tmq-linux-amd64
          sudo mv tmq-linux-amd64 /usr/local/bin/tmq

      - name: Update versions
        run: |
          tmq ".version = \"${{ github.event.inputs.new_version }}\"" -i config/*.toml

      - name: Validate changes
        run: tmq --validate config/*.toml

      - name: Create pull request
        uses: peter-evans/create-pull-request@v5
        with:
          title: "Update configuration versions to ${{ github.event.inputs.new_version }}"
          body: "Bulk update of version fields across all configuration files"
```

### Configuration Management
```bash
#!/bin/bash
# Configuration management script

CONFIG_DIR="config"
ENVIRONMENTS=("development" "staging" "production")

# Update environment-specific settings
for env in "${ENVIRONMENTS[@]}"; do
    echo "Updating $env configuration..."

    tmq ".environment = \"$env\"" -i "$CONFIG_DIR/$env.toml"

    case $env in
        "development")
            tmq '.debug = true' -i "$CONFIG_DIR/$env.toml"
            tmq '.database.host = "localhost"' -i "$CONFIG_DIR/$env.toml"
            ;;
        "staging")
            tmq '.debug = false' -i "$CONFIG_DIR/$env.toml"
            tmq '.database.host = "staging-db.example.com"' -i "$CONFIG_DIR/$env.toml"
            ;;
        "production")
            tmq '.debug = false' -i "$CONFIG_DIR/$env.toml"
            tmq '.database.host = "prod-db.example.com"' -i "$CONFIG_DIR/$env.toml"
            ;;
    esac
done

# Validate all configurations
echo "Validating all configurations..."
tmq --validate "$CONFIG_DIR"/*.toml

echo "Configuration update complete"
```

## Monitoring and Reporting

### Progress Tracking
```bash
#!/bin/bash
# Bulk operation with progress reporting

files=(config/*.toml)
total=${#files[@]}
processed=0

echo "Processing $total files..."

for file in "${files[@]}"; do
    ((processed++))
    echo -ne "Progress: $processed/$total\r"

    if ! tmq '.version' "$file" >/dev/null; then
        echo "Failed: $file"
    fi
done

echo -e "\nProcessing complete"
```

### Summary Reports
```bash
#!/bin/bash
# Generate summary report of configuration status

echo "=== Configuration Status Report ==="
echo

# Count total files
total_files=$(ls config/*.toml 2>/dev/null | wc -l)
echo "Total configuration files: $total_files"

# Count valid files
valid_files=$(tmq --validate config/*.toml >/dev/null 2>&1 && echo "$total_files" || echo "0")
echo "Valid TOML files: $valid_files"

# List versions
echo
echo "Version summary:"
tmq '.version // "unknown"' config/*.toml | sort | uniq -c | sort -nr

# Check for common issues
echo
echo "Potential issues:"
for file in config/*.toml; do
    # Check for deprecated settings
    if tmq '.deprecated_setting' "$file" >/dev/null 2>&1; then
        echo "  ⚠️  $file contains deprecated setting"
    fi
done
```

## Best Practices

### Backup Strategy
```bash
# Always backup before bulk operations
backup_dir="backup-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$backup_dir"
cp config/*.toml "$backup_dir/"
echo "Backup created in $backup_dir"
```

### Dry Run Testing
```bash
# Test bulk operations with dry-run first
echo "Testing bulk update (dry-run)..."
if tmq '.version = "test"' --dry-run config/*.toml >/dev/null; then
    echo "✓ Dry-run successful, proceeding with actual update..."
    tmq '.version = "test"' -i config/*.toml
else
    echo "✗ Dry-run failed, aborting"
    exit 1
fi
```

### Rollback Capability
```bash
#!/bin/bash
# Bulk update with rollback capability

ROLLBACK_FILE="rollback-$(date +%s).sh"

# Generate rollback script
cat > "$ROLLBACK_FILE" << 'EOF'
#!/bin/bash
echo "Rolling back configuration changes..."
EOF

# Store current values for rollback
for file in config/*.toml; do
    version=$(tmq '.version' "$file" 2>/dev/null || echo '""')
    echo "tmq '.version = $version' -i '$file'" >> "$ROLLBACK_FILE"
done

chmod +x "$ROLLBACK_FILE"

# Perform bulk update
tmq '.version = "2.0.0"' -i config/*.toml

# Test if everything works
if ./test-configuration.sh; then
    echo "✓ Update successful"
    rm "$ROLLBACK_FILE"  # Clean up rollback script
else
    echo "✗ Update failed, rolling back..."
    ./"$ROLLBACK_FILE"
    echo "Rollback complete"
fi
```
