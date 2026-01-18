#!/bin/bash

# Helper script to record animations
# Usage: ./record-demo.sh [fire|matrix|all]

set -e

EFFECT=${1:-all}

echo "🎬 Recording Animation Demo"
echo "============================"
echo ""

check_vhs() {
    if ! command -v vhs &> /dev/null; then
        echo "❌ VHS is not installed"
        echo ""
        echo "Install VHS:"
        echo "  macOS:  brew install vhs"
        echo "  Linux:  go install github.com/charmbracelet/vhs@latest"
        echo ""
        echo "Or use alternative recording tools (see examples/RECORDING.md)"
        exit 1
    fi
    echo "✅ VHS found"
}

record_animation() {
    local name=$1
    echo ""
    echo "📹 Recording $name animation..."
    npm run record:$name
    
    if [ -f "demos/${name}.gif" ]; then
        echo "✅ Saved to demos/${name}.gif"
        
        # Show file size
        local size=$(du -h "demos/${name}.gif" | cut -f1)
        echo "   File size: $size"
    else
        echo "❌ Failed to create demos/${name}.gif"
    fi
}

cd "$(dirname "$0")"

echo "Checking dependencies..."
check_vhs

if [ "$EFFECT" = "all" ]; then
    record_animation "fire"
    record_animation "matrix"
elif [ "$EFFECT" = "fire" ] || [ "$EFFECT" = "matrix" ]; then
    record_animation "$EFFECT"
else
    echo "❌ Unknown effect: $EFFECT"
    echo "Usage: $0 [fire|matrix|all]"
    exit 1
fi

echo ""
echo "✨ Recording complete!"
echo ""
echo "To optimize file size:"
echo "  gifsicle -O3 --colors 256 demos/fire.gif -o demos/fire-optimized.gif"
