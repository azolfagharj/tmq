#!/bin/bash

# Build script for MkDocs documentation
# This script builds the documentation and copies it to docs/statics/
# It automatically sets up a Python virtual environment and installs dependencies
#
# Environment variables:
#   CI_MODE: If set to "true", keeps venv for caching and disables colors
#   KEEP_VENV: If set to "true", keeps venv after build (useful for CI caching)

set -e

# Detect CI environment
if [ -n "$CI" ] || [ -n "$GITHUB_ACTIONS" ] || [ "${CI_MODE:-false}" = "true" ]; then
    CI_MODE=true
else
    CI_MODE=false
fi

# Colors for output (disabled in CI)
if [ "$CI_MODE" = "true" ]; then
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
else
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    NC='\033[0m'
fi

# Get the script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
MKDOCS_CONFIG="$SCRIPT_DIR/mkdocs.yaml"
REQUIREMENTS_FILE="$SCRIPT_DIR/requirements.txt"
VENV_DIR="$SCRIPT_DIR/venv"
DOCS_DIR="$PROJECT_ROOT/docs"
STATICS_DIR="$DOCS_DIR/statics"
SITE_DIR="$PROJECT_ROOT/docs-site-temp"

echo -e "${GREEN}=== GajIn Documentation Build ===${NC}"
if [ "$CI_MODE" = "true" ]; then
    echo "Running in CI mode"
fi
echo ""

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}Error: python3 is not installed${NC}"
    echo "Please install Python 3.7 or higher"
    exit 1
fi

# Check Python version (minimum 3.7)
PYTHON_VERSION=$(python3 -c 'import sys; print(".".join(map(str, sys.version_info[:2])))')
PYTHON_MAJOR=$(echo "$PYTHON_VERSION" | cut -d. -f1)
PYTHON_MINOR=$(echo "$PYTHON_VERSION" | cut -d. -f2)

if [ "$PYTHON_MAJOR" -lt 3 ] || { [ "$PYTHON_MAJOR" -eq 3 ] && [ "$PYTHON_MINOR" -lt 7 ]; }; then
    echo -e "${RED}Error: Python 3.7 or higher is required${NC}"
    echo "Current version: $PYTHON_VERSION"
    exit 1
fi

echo -e "${BLUE}Python version: $PYTHON_VERSION${NC}"
echo ""

# Check if config file exists
if [ ! -f "$MKDOCS_CONFIG" ]; then
    echo -e "${RED}Error: mkdocs.yaml not found at $MKDOCS_CONFIG${NC}"
    exit 1
fi

# Check if requirements file exists
if [ ! -f "$REQUIREMENTS_FILE" ]; then
    echo -e "${RED}Error: requirements.txt not found at $REQUIREMENTS_FILE${NC}"
    exit 1
fi

# Step 1: Setup virtual environment
echo -e "${YELLOW}Step 1: Setting up Python virtual environment...${NC}"

if [ ! -d "$VENV_DIR" ]; then
    echo "  Creating virtual environment..."
    python3 -m venv "$VENV_DIR"
    echo "  ✓ Virtual environment created"
else
    echo "  ✓ Virtual environment already exists (reusing for faster build)"
fi

# Activate virtual environment
echo "  Activating virtual environment..."
# shellcheck disable=SC1091
source "$VENV_DIR/bin/activate"

# Upgrade pip
echo "  Upgrading pip..."
if [ "$CI_MODE" = "true" ]; then
    pip install --upgrade pip > /dev/null 2>&1
else
    pip install --quiet --upgrade pip
fi

# Install/upgrade dependencies
echo "  Installing dependencies from requirements.txt..."
if [ "$CI_MODE" = "true" ]; then
    if ! pip install -r "$REQUIREMENTS_FILE" > /dev/null 2>&1; then
        echo -e "${RED}Error: Failed to install dependencies${NC}"
        echo "Retrying with verbose output..."
        pip install -r "$REQUIREMENTS_FILE"
        exit 1
    fi
else
    if ! pip install --quiet -r "$REQUIREMENTS_FILE"; then
        echo -e "${RED}Error: Failed to install dependencies${NC}"
        exit 1
    fi
fi

echo "  ✓ Dependencies installed successfully"
echo ""

# Change to project root directory
cd "$PROJECT_ROOT"

# Step 2: Clean old build artifacts
echo -e "${YELLOW}Step 2: Cleaning old build artifacts...${NC}"
# Remove old site directory if exists
if [ -d "$SITE_DIR" ]; then
    rm -rf "$SITE_DIR"
    echo "  ✓ Removed $SITE_DIR"
fi

# Remove old statics directory if exists
if [ -d "$STATICS_DIR" ]; then
    rm -rf "$STATICS_DIR"
    echo "  ✓ Removed $STATICS_DIR"
fi

echo ""
echo -e "${YELLOW}Step 3: Building documentation with MkDocs...${NC}"
# Build documentation using mkdocs from venv
# Change to script directory so relative paths in mkdocs.yaml work correctly
cd "$SCRIPT_DIR"
# Use --strict only in CI mode, allow warnings in local builds
if [ "$CI_MODE" = "true" ]; then
    if ! "$VENV_DIR/bin/mkdocs" build --config-file "$(basename "$MKDOCS_CONFIG")" --strict; then
        echo -e "${RED}Error: MkDocs build failed${NC}"
        exit 1
    fi
else
    if ! "$VENV_DIR/bin/mkdocs" build --config-file "$(basename "$MKDOCS_CONFIG")"; then
        echo -e "${RED}Error: MkDocs build failed${NC}"
        exit 1
    fi
fi
# Return to project root
cd "$PROJECT_ROOT"

echo "  ✓ Documentation built successfully"

# Check if site directory was created
if [ ! -d "$SITE_DIR" ]; then
    echo -e "${RED}Error: Site directory was not created at $SITE_DIR${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Step 4: Copying files to docs/statics/...${NC}"
# Copy site directory to statics
if ! cp -r "$SITE_DIR" "$STATICS_DIR"; then
    echo -e "${RED}Error: Failed to copy files to statics directory${NC}"
    exit 1
fi

echo "  ✓ Files copied successfully"

echo ""
echo -e "${YELLOW}Step 5: Cleaning up temporary files...${NC}"
# Remove temporary site directory
rm -rf "$SITE_DIR"
echo "  ✓ Removed temporary site directory"

# Deactivate virtual environment before deletion
deactivate 2>/dev/null || true

# Remove virtual environment (unless in CI mode or KEEP_VENV is set)
if [ "$CI_MODE" = "true" ] || [ "${KEEP_VENV:-false}" = "true" ]; then
    echo "  ℹ Keeping virtual environment (CI mode or KEEP_VENV=true)"
    if [ "$CI_MODE" = "true" ]; then
        echo "  ℹ Venv location: $VENV_DIR (can be cached in CI)"
    fi
else
    # Remove virtual environment
    if [ -d "$VENV_DIR" ]; then
        rm -rf "$VENV_DIR"
        echo "  ✓ Removed virtual environment"
    fi
fi

echo ""
echo -e "${GREEN}=== Build Complete ===${NC}"
echo ""
echo "Documentation has been built and is available in:"
echo "  $STATICS_DIR"
echo ""

if [ "$CI_MODE" = "true" ]; then
    echo "Build completed successfully in CI environment."
    echo "Output directory: $STATICS_DIR"
else
    echo "You can view it by opening index.html in a web browser:"
    echo "  file://$STATICS_DIR/index.html"
fi
echo ""

