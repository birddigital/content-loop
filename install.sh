#!/bin/bash
set -e

INSTALL_DIR="$HOME/.local/bin"
LOG_DIR="$HOME/.local/var/log"
SRC_DIR="$(pwd)"
SERVICE_NAME="content-loop"

echo "🚀 Installing ContentLoop..."

# Create directories
mkdir -p "$INSTALL_DIR"
mkdir -p "$LOG_DIR"

# Copy binary
echo "📦 Installing binary to $INSTALL_DIR"
cp "$SRC_DIR/content-loop" "$INSTALL_DIR/content-loop"
chmod +x "$INSTALL_DIR/content-loop"

# Copy config
if [ ! -f "$SRC_DIR/config.yaml" ]; then
    echo "📝 Creating config.yaml from example"
    cp "$SRC_DIR/config.yaml.example" "$SRC_DIR/config.yaml"
    echo "⚠️  Edit $SRC_DIR/config.yaml with your API credentials"
fi

# Detect OS and install service
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    echo "🍎 Installing macOS launchd service..."
    PLIST="$HOME/Library/LaunchAgents/com.birddigital.content-loop.plist"

    # Update plist with correct paths
    sed "s|/Users/birddigital|$HOME|g" "$SRC_DIR/com.birddigital.content-loop.plist" > "$PLIST"

    # Load the service
    launchctl unload "$PLIST" 2>/dev/null || true
    launchctl load "$PLIST"

    echo "✅ ContentLoop installed and started (macOS)"
    echo "📊 Logs: $LOG_DIR/content-loop.*.log"
    echo "🛑 Stop: launchctl unload $PLIST"
    echo "🔄 Restart: launchctl unload $PLIST && launchctl load $PLIST"

elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    echo "🐧 Installing systemd service..."
    SERVICE_DST="/etc/systemd/system/$SERVICE_NAME.service"

    # Install service file (requires sudo)
    if [ -w "/etc/systemd/system" ]; then
        sed "s|/home/birddigital|$HOME|g" "$SRC_DIR/content-loop.service" > "$SERVICE_DST"
        systemctl daemon-reload
        systemctl enable "$SERVICE_NAME"
        systemctl start "$SERVICE_NAME"
        echo "✅ ContentLoop installed and started (Linux)"
        echo "📊 Logs: journalctl -u $SERVICE_NAME -f"
        echo "🛑 Stop: systemctl stop $SERVICE_NAME"
        echo "🔄 Restart: systemctl restart $SERVICE_NAME"
    else
        echo "⚠️  Requires sudo to install systemd service"
        echo "   Run: sudo cp $SRC_DIR/content-loop.service /etc/systemd/system/"
        echo "   Then: systemctl daemon-reload && systemctl enable $SERVICE_NAME && systemctl start $SERVICE_NAME"
    fi
else
    echo "⚠️  Unsupported OS. Service not installed."
    echo "📝 You can run manually: $INSTALL_DIR/content-loop"
fi

echo ""
echo "✨ Installation complete!"
echo "📖 Edit config: $SRC_DIR/config.yaml"
