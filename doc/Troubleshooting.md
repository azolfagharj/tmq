# Troubleshooting

Common issues and solutions when using tmq.

## Installation Issues

### Command Not Found
```bash
tmq --version
# bash: tmq: command not found
```

**Solutions:**
1. **Check PATH:**
   ```bash
   echo $PATH
   which tmq
   ```

2. **Use full path:**
   ```bash
   /usr/local/bin/tmq --version
   ./tmq --version
   ```

3. **Reinstall:**
   ```bash
   # Download again
   wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
   chmod +x tmq-linux-amd64
   sudo mv tmq-linux-amd64 /usr/local/bin/tmq
   ```

### Permission Denied
```bash
./tmq --version
# bash: ./tmq: Permission denied
```

**Solution:**
```bash
chmod +x tmq
./tmq --version
```

### Unsupported Platform
```bash
# Trying to run ARM binary on x86 system
./tmq-linux-arm64 --version
# exec format error
```

**Solution:** Download the correct binary for your platform:
- Linux x86_64: `tmq-linux-amd64`
- Linux ARM64: `tmq-linux-arm64`
- macOS Intel: `tmq-darwin-amd64`
- macOS Apple Silicon: `tmq-darwin-arm64`
- Windows: `tmq-windows-amd64.exe`

## Query Issues

### Invalid Query Path
```bash
tmq '.invalid..path' config.toml
# ERROR: Invalid query path
```

**Common mistakes:**
1. **Empty path segments:**
   ```bash
   # Wrong
   tmq '.database..' config.toml

   # Right
   tmq '.database' config.toml
   ```

2. **Invalid characters:**
   ```bash
   # Wrong (spaces in path)
   tmq '.my key' config.toml

   # Right (use underscores or different structure)
   tmq '.my_key' config.toml
   ```

3. **Unmatched brackets:**
   ```bash
   # Wrong
   tmq '.array[0' config.toml

   # Right
   tmq '.array[0]' config.toml
   ```

### Key Not Found
```bash
tmq '.nonexistent' config.toml
# ERROR: key 'nonexistent' not found
```

**Solutions:**
1. **Check key exists:**
   ```bash
   # List all keys
   tmq '.' config.toml

   # Check specific section
   tmq '.database' config.toml
   ```

2. **Case sensitivity:**
   ```toml
   [Database]  # Capital D
   host = "localhost"
   ```
   ```bash
   # Wrong
   tmq '.database.host' config.toml

   # Right
   tmq '.Database.host' config.toml
   ```

3. **Array vs object:**
   ```toml
   # This is an array of tables
   [[servers]]
   name = "server1"

   # This is a table
   [database]
   host = "localhost"
   ```

### Array Index Out of Bounds
```bash
tmq '.servers[10]' config.toml
# ERROR: array index 10 out of bounds
```

**Solution:** Check array length first:
```bash
# See the array
tmq '.servers' config.toml

# Get array length (in scripts)
length=$(tmq '.servers | length' config.toml)
```

## Modification Issues

### Cannot Modify File
```bash
tmq '.version = "2.0"' -i config.toml
# ERROR: File operation error
# DETAILS: Cannot write to file 'config.toml'
```

**Solutions:**
1. **Check permissions:**
   ```bash
   ls -la config.toml
   chmod 644 config.toml
   ```

2. **Check disk space:**
   ```bash
   df -h
   ```

3. **File locked:**
   ```bash
   # Check if file is open by another process
   lsof config.toml
   ```

### Invalid Set Expression
```bash
tmq '.version = ' -i config.toml
# ERROR: Invalid set expression
```

**Common syntax errors:**
1. **Missing value:**
   ```bash
   # Wrong
   tmq '.version =' -i config.toml

   # Right
   tmq '.version = "2.0"' -i config.toml
   ```

2. **Unquoted strings:**
   ```bash
   # Wrong (if value contains spaces or special chars)
   tmq '.path = /usr/local/bin' -i config.toml

   # Right
   tmq '.path = "/usr/local/bin"' -i config.toml
   ```

3. **Invalid TOML values:**
   ```bash
   # Wrong
   tmq '.value = [unclosed array' -i config.toml

   # Right
   tmq '.value = ["item1", "item2"]' -i config.toml
   ```

## Validation Issues

### TOML Parse Error
```bash
tmq --validate config.toml
# ERROR: TOML parsing failed
# DETAILS: Expected newline at line 5, column 10
```

**Common TOML syntax errors:**
1. **Missing quotes:**
   ```toml
   # Wrong
   title = My App

   # Right
   title = "My App"
   ```

2. **Invalid array syntax:**
   ```toml
   # Wrong
   items = [1, 2, 3

   # Right
   items = [1, 2, 3]
   ```

3. **Invalid table definition:**
   ```toml
   # Wrong
   [section
   key = "value"

   # Right
   [section]
   key = "value"
   ```

### File Encoding Issues
```bash
tmq --validate config.toml
# ERROR: Invalid UTF-8 encoding
```

**Solution:** Convert file to UTF-8:
```bash
# Check current encoding
file config.toml

# Convert to UTF-8
iconv -f latin1 -t utf8 config.toml > config_utf8.toml
mv config_utf8.toml config.toml
```

## Comparison Issues

### Files Are Different But Should Match
```bash
tmq --compare file1.toml file2.toml
# Files are different
```

**Possible causes:**
1. **Different formatting:**
   ```toml
   # file1.toml
   host="localhost"

   # file2.toml
   host = "localhost"
   ```
   These are semantically identical but formatted differently.

2. **Comments ignored:**
   TOML comments are not part of the data structure.

3. **Order differences:**
   ```toml
   # Different order, same content
   [table]
   b = 2
   a = 1

   [table]
   a = 1
   b = 2
   ```

## Performance Issues

### Slow Operations on Large Files
```bash
time tmq '.' large.toml > /dev/null
# Takes several seconds
```

**Solutions:**
1. **Use specific queries:**
   ```bash
   # Instead of full file
   tmq '.' large.toml

   # Use specific path
   tmq '.database.host' large.toml
   ```

2. **Output to file instead of stdout:**
   ```bash
   tmq '.' large.toml -o json > output.json
   ```

3. **Check file size:**
   ```bash
   ls -lh large.toml
   # If > 50MB, consider splitting
   ```

### High Memory Usage
```bash
# Monitor memory
/usr/bin/time -v tmq '.' large.toml
```

**Solutions:**
1. **Process in chunks:**
   ```bash
   # Instead of whole file, process sections
   tmq '.section1' large.toml
   tmq '.section2' large.toml
   ```

2. **Use streaming for large outputs:**
   ```bash
   tmq '.' large.toml | head -100
   ```

## Bulk Operation Issues

### Too Many Files
```bash
tmq '.version' config/*.toml
# Argument list too long
```

**Solution:** Use `find` or process in batches:
```bash
# Use find
find config/ -name "*.toml" -exec tmq '.version' {} \;

# Process in batches
for file in config/*.toml; do
    tmq '.version' "$file"
done
```

### Inconsistent Results
```bash
tmq '.version' config/*.toml
# Some files show different output format
```

**Cause:** Different files may have different structures.

**Solution:** Handle missing keys:
```bash
# Use default values
tmq '.version // "unknown"' config/*.toml
```

## Scripting Issues

### Exit Code Handling
```bash
#!/bin/bash
tmq '.nonexistent' config.toml
echo "Exit code: $?"
# Always prints exit code, even if command fails
```

**Solution:** Use proper error handling:
```bash
#!/bin/bash
set -e  # Exit on first error

if ! tmq '.nonexistent' config.toml >/dev/null 2>&1; then
    echo "Key not found"
    exit 1
fi
```

### Output Parsing Issues
```bash
# Trying to parse multi-line output
version=$(tmq '.' config.toml | grep version)
# Unreliable
```

**Solution:** Use specific queries:
```bash
version=$(tmq '.version' config.toml)
```

### Variable Interpolation
```bash
# Wrong
tmq '.version = $VERSION' -i config.toml

# Right
tmq ".version = \"$VERSION\"" -i config.toml
```

## Platform-Specific Issues

### Windows Path Issues
```cmd
REM Wrong (backslashes in paths)
tmq ".path = C:\Program Files\app" -i config.toml

REM Right (forward slashes or quoted)
tmq ".path = \"C:/Program Files/app\"" -i config.toml
```

### macOS Permission Issues
```bash
# When installed in /usr/local/bin
sudo tmq --version
# May require sudo if installed system-wide
```

## Network and Download Issues

### Proxy Issues
```bash
# If behind corporate proxy
export https_proxy=http://proxy.company.com:8080
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
```

### Certificate Issues
```bash
# Skip SSL verification (not recommended for production)
wget --no-check-certificate https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
```

### Firewall Blocking Downloads
```bash
# Check connectivity
curl -I https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64

# Use alternative download method
curl -L -o tmq https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
```

## Getting Help

### Debug Information
```bash
# Collect system info
uname -a
tmq --version

# Test with simple file
echo 'test = "value"' > test.toml
tmq '.test' test.toml
```

### Issue Reporting
When reporting issues, include:
1. **tmq version:** `tmq --version`
2. **Platform:** `uname -a`
3. **Command that fails**
4. **Full error output**
5. **Sample TOML file** (if applicable)
6. **Expected vs actual behavior**

### Community Support
- **GitHub Issues:** https://github.com/azolfagharj/tmq/issues
- **Discussions:** https://github.com/azolfagharj/tmq/discussions

## Quick Fixes

| Issue | Quick Fix |
|-------|-----------|
| Command not found | `export PATH=$PATH:/usr/local/bin` |
| Permission denied | `chmod +x tmq` |
| Invalid query | Check syntax: `tmq --help` |
| File not found | `ls -la config.toml` |
| Parse error | Validate: `tmq --validate config.toml` |
| Memory issues | Use specific queries instead of full file |
| Bulk operations slow | Process in smaller batches |