# MkDocs Documentation Build Script

This directory contains scripts and configuration for building the GajIn documentation using MkDocs.

## Prerequisites

- **Python 3.7 or higher** (required)
- pip (comes with Python, no separate installation needed)

**Note:** The build script automatically sets up a Python virtual environment and installs all dependencies. You don't need to install anything manually!

## Usage

### Build Documentation (Recommended)

Simply run the build script from the project root:

```bash
./.github/scripts/mkdocs/build.sh
```

Or from the mkdocs directory:

```bash
cd .github/scripts/mkdocs
./build.sh
```

The script will automatically:
1. **Check Python installation** and version
2. **Create a virtual environment** in `.github/scripts/mkdocs/venv/`
3. **Install all dependencies** from `requirements.txt` into the virtual environment
4. **Clean old build artifacts** (`docs/statics/` and `docs/site/`)
5. **Build documentation** using MkDocs
6. **Copy built files** from `docs/site/` to `docs/statics/`
7. **Clean up temporary files**

**Local builds:** The virtual environment is created fresh for each build and automatically deleted at the end.

**CI builds:** In GitHub Actions or CI environments, the virtual environment is kept for caching (set `CI_MODE=true` or `KEEP_VENV=true`).

### Manual Build (Advanced)

If you prefer to build manually without the script:

```bash
# Activate virtual environment (created by build script)
source .github/scripts/mkdocs/venv/bin/activate

# Build documentation
mkdocs build --config-file .github/scripts/mkdocs/mkdocs.yaml

# Copy to statics
rm -rf docs/statics
cp -r docs/site docs/statics
rm -rf docs/site

# Deactivate virtual environment
deactivate
```

Or if you want to use a system-wide installation:

```bash
# Install dependencies globally (not recommended)
pip install -r .github/scripts/mkdocs/requirements.txt

# Build
mkdocs build --config-file .github/scripts/mkdocs/mkdocs.yaml

# Copy to statics
rm -rf docs/statics
cp -r docs/site docs/statics
rm -rf docs/site
```

## Configuration

The MkDocs configuration is in `mkdocs.yaml`. It includes:

- **Theme**: Material for MkDocs
- **Plugins**:
  - Search
  - Git revision date (localized)
  - Git committers
  - Minify (HTML, JS, CSS)
- **Navigation**: Manual ordering (README → USAGE → ARCHITECTURE → MIGRATION)

## GitHub Actions Integration

The script is optimized for use in GitHub Actions. It automatically detects CI environments and:
- Disables colored output for better log readability
- Keeps the virtual environment for caching (speeds up subsequent builds)
- Provides better error messages for CI debugging

### Example Workflow

See `docs-build-example.yml` for a complete example workflow. Here's a minimal example:

```yaml
name: Build Documentation

on:
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - '.github/scripts/mkdocs/**'

jobs:
  build-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Required for git-revision-date-localized plugin

      - uses: actions/setup-python@v5
        with:
          python-version: '3.11'
          cache: 'pip'

      - name: Cache MkDocs dependencies
        uses: actions/cache@v4
        with:
          path: .github/scripts/mkdocs/venv
          key: mkdocs-${{ runner.os }}-${{ hashFiles('.github/scripts/mkdocs/requirements.txt') }}

      - run: chmod +x .github/scripts/mkdocs/build.sh
      - run: ./.github/scripts/mkdocs/build.sh
        env:
          CI_MODE: true
          KEEP_VENV: true

      - uses: actions/upload-artifact@v4
        with:
          name: documentation
          path: docs/statics/**
```

### Environment Variables

- `CI_MODE`: Set to `true` to enable CI optimizations (auto-detected in GitHub Actions)
- `KEEP_VENV`: Set to `true` to keep venv after build (useful for caching)

## Output

Built documentation will be available in `docs/statics/`. You can view it by opening `docs/statics/index.html` in a web browser.

## Troubleshooting

### Python not found

If you get "python3: command not found", make sure Python 3.7+ is installed:

```bash
# Check Python version
python3 --version

# On Ubuntu/Debian
sudo apt-get install python3 python3-venv python3-pip

# On macOS (with Homebrew)
brew install python3
```

### Permission denied

If the build script is not executable:

```bash
chmod +x .github/scripts/mkdocs/build.sh
```

### Virtual environment issues

**Local builds:** The virtual environment is automatically created and deleted with each build. If you encounter issues, make sure Python 3.7+ is installed and try running the script again.

**CI builds:** In CI environments, the venv is kept for caching. If you need to force a fresh build, delete the cache in your CI platform.

### Build errors

Check the MkDocs configuration file for syntax errors:

```bash
# Activate venv first
source .github/scripts/mkdocs/venv/bin/activate

# Build with verbose output
mkdocs build --config-file .github/scripts/mkdocs/mkdocs.yaml --verbose

# Deactivate
deactivate
```

### Dependencies installation fails

If pip installation fails, try upgrading pip first:

```bash
# Activate venv
source .github/scripts/mkdocs/venv/bin/activate

# Upgrade pip
pip install --upgrade pip

# Install dependencies
pip install -r .github/scripts/mkdocs/requirements.txt

# Deactivate
deactivate
```

