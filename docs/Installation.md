# Installation

tmq is distributed as pre-compiled binaries for multiple platforms. No external dependencies required.

## Download Binary

Download the latest release from [GitHub Releases](https://github.com/azolfagharj/tmq/releases).

### Available Binaries

| Platform | Architecture | Filename |
|----------|-------------|----------|
| Linux | AMD64 | `tmq-linux-amd64` |
| Linux | ARM64 | `tmq-linux-arm64` |
| macOS | Intel | `tmq-darwin-amd64` |
| macOS | Apple Silicon | `tmq-darwin-arm64` |
| Windows | AMD64 | `tmq-windows-amd64.exe` |

## Quick Setup

### Linux (AMD64)
```bash
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
chmod +x tmq-linux-amd64
sudo mv tmq-linux-amd64 /usr/local/bin/tmq
```

### Linux (ARM64)
```bash
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-arm64
chmod +x tmq-linux-arm64
sudo mv tmq-linux-arm64 /usr/local/bin/tmq
```

### macOS (Apple Silicon)
```bash
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-darwin-arm64
chmod +x tmq-darwin-arm64
sudo mv tmq-darwin-arm64 /usr/local/bin/tmq
```

### macOS (Intel)
```bash
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-darwin-amd64
chmod +x tmq-darwin-amd64
sudo mv tmq-darwin-amd64 /usr/local/bin/tmq
```

### Windows (AMD64)
1. Download: https://github.com/azolfagharj/tmq/releases/latest/download/tmq-windows-amd64.exe
2. Rename to `tmq.exe`
3. Add to your PATH

## Manual Installation

1. **Download** the appropriate binary for your system
2. **Make executable** (Linux/macOS):
   ```bash
   chmod +x tmq-*
   ```
3. **Rename** to `tmq` (or `tmq.exe` on Windows):
   ```bash
   mv tmq-linux-amd64 tmq
   ```
4. **Move to PATH** (optional but recommended):
   ```bash
   sudo mv tmq /usr/local/bin/
   ```

## Build from Source

If you prefer to build from source:

### Prerequisites
- Go 1.23 or later

### Build Steps
```bash
git clone https://github.com/azolfagharj/tmq.git
cd tmq
go build -o bin/tmq ./cmd/tmq
```

## Verify Installation

After installation, verify tmq is working:

```bash
tmq --version
# Should show: tmq version 1.0.1

tmq --help
# Should show help text
```

## System Requirements

- **OS**: Linux, macOS, Windows
- **Architecture**: AMD64, ARM64
- **Memory**: Minimal (works with < 10MB RAM)
- **Storage**: ~5MB for the binary
- **No external dependencies** - completely self-contained

## Updating

To update to the latest version:

1. Download the new binary from [releases](https://github.com/azolfagharj/tmq/releases)
2. Replace the old binary
3. Make sure it's executable (`chmod +x tmq`)

## Troubleshooting

### Permission Denied
If you get "permission denied" when running tmq:
```bash
chmod +x /path/to/tmq
```

### Command Not Found
If `tmq` command is not found:
- Make sure it's in your PATH: `echo $PATH`
- Or use the full path: `/usr/local/bin/tmq`

### Download Issues
If wget fails, try using curl:
```bash
curl -L -o tmq https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
```