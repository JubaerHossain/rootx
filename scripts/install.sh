#!/bin/bash

set -e

if [[ "$EUID" -ne 0 ]]; then
  echo "Please run as root"
  exit 1
fi

BIN_DIR="/usr/local/bin"
BIN_NAME="rootx"

cp $BIN_NAME $BIN_DIR
chmod +x $BIN_DIR/$BIN_NAME

echo "rootx installed successfully to $BIN_DIR"
