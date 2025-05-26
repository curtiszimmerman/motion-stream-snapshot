#!/bin/bash

# Exit on error
set -e

VERSION="1.0.0"
SERVICE_FILE="motion-snapshot.service"
MOTION_CONF="/etc/motion/motion.conf"
MOTION_LIB="/var/lib/motion"

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if a package is installed
package_installed() {
    dpkg -l "$1" | grep -q "^ii"
}

# Function to check version
check_version() {
    if command_exists motion-snapshot-server; then
        current_version=$(motion-snapshot-server --version 2>/dev/null || echo "unknown")
        echo "[*] Current version: $current_version"
        echo "[*] Installing version: $VERSION"
    fi
}

# Function to uninstall
uninstall() {
    echo "[*] Uninstalling motion-snapshot-server..."
    # Stop and disable the service if it exists
    if systemctl is-active --quiet motion-snapshot; then
        systemctl stop motion-snapshot
        systemctl disable motion-snapshot
    fi
    rm -f /usr/bin/motion-snapshot-server
    rm -f /etc/systemd/system/motion-snapshot.service
    systemctl daemon-reload
    echo "[*] Uninstallation complete! You may wish to uninstall motion as well (sudo apt-get remove motion)."
    exit 0
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --uninstall)
            uninstall
            ;;
        --version)
            echo "[*] motion-snapshot-server installer version: $VERSION"
            exit 0
            ;;
        *)
            echo "[!] Unknown option: $1"
            echo "    Usage: $0 [--uninstall|--version]"
            exit 1
            ;;
    esac
    shift
done

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "[!] Please run as root (use sudo)"
    exit 1
fi

# Check for existing installation
if command_exists motion-snapshot-server; then
    echo "[!] motion-snapshot-server is already installed."
    check_version
    read -p "[!] Do you want to reinstall? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

# Check if motion is installed
if ! package_installed motion; then
    echo "[*] Installing motion package..."
    apt-get update
    apt-get install -y motion

    echo "[*] Creating motion library directory..."
    mkdir -p "$MOTION_LIB"
    echo "[*] Installing default motion configuration..."
    cp motion.conf "$MOTION_CONF"
    chmod 644 "$MOTION_CONF"
else
    echo "[!] Motion package is already installed."
fi

echo "[*] Building motion-snapshot-server..."
make build

echo "[*] Installing motion-snapshot-server to /usr/bin/..."
cp bin/motion-snapshot-server /usr/bin/
chmod +x /usr/bin/motion-snapshot-server

echo "[*] Installing systemd service..."
cp "$SERVICE_FILE" /etc/systemd/system/
systemctl daemon-reload
systemctl enable motion-snapshot
systemctl start motion-snapshot

echo "[*] Installation complete!"
echo "[*] You can now run motion-snapshot-server from anywhere using the command: motion-snapshot-server"
echo "[*] The service has been installed and started. You can manage it using:"
echo "  sudo systemctl start motion-snapshot    # Start the service"
echo "  sudo systemctl stop motion-snapshot     # Stop the service"
echo "  sudo systemctl status motion-snapshot   # Check service status"
echo "[*] To uninstall, run: sudo $0 --uninstall" 