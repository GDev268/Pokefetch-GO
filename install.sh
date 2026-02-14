#!/bin/sh
set -e

CUR_DIR=$(pwd)
BIN_NAME="out/pokefetch"        # your GO binary path
INSTALL_DIR="$HOME/.local/bin"  # target install directory

# Extract just the file name
BIN_FILE=$(basename "$BIN_NAME")

echo "==> Installing $CUR_DIR/$BIN_NAME to $INSTALL_DIR/$BIN_FILE"

# Ensure install directory exists
mkdir -p "$INSTALL_DIR"

# Copy binary
cp "$CUR_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_FILE"
chmod +x "$INSTALL_DIR/$BIN_FILE"

echo "Installation complete."

# Optional PATH warning
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo
    echo "WARNING: $INSTALL_DIR is not in your PATH."
    echo "Add this line to your shell config (~/.zshrc or ~/.bashrc):"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi
