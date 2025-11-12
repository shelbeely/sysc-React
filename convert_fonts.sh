#!/bin/bash
# Batch convert FIGlet fonts (.flf) to .bit format

if [ -z "$1" ]; then
    echo "Usage: $0 <directory-with-flf-files> [output-directory]"
    echo ""
    echo "Example:"
    echo "  $0 ~/figlet-fonts assets/fonts"
    echo "  $0 /usr/share/figlet"
    exit 1
fi

INPUT_DIR="$1"
OUTPUT_DIR="${2:-assets/fonts}"

# Check if flf2bit exists
if [ ! -f "./flf2bit" ]; then
    echo "Building flf2bit..."
    go build ./cmd/flf2bit
    if [ $? -ne 0 ]; then
        echo "Error: Failed to build flf2bit"
        exit 1
    fi
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Count total .flf files
total=$(find "$INPUT_DIR" -name "*.flf" | wc -l)
if [ $total -eq 0 ]; then
    echo "No .flf files found in $INPUT_DIR"
    exit 1
fi

echo "Found $total FIGlet fonts in $INPUT_DIR"
echo "Converting to $OUTPUT_DIR..."
echo ""

# Convert each .flf file
count=0
success=0
failed=0

find "$INPUT_DIR" -name "*.flf" | while read -r flf_file; do
    count=$((count + 1))

    # Get base name without extension
    base=$(basename "$flf_file" .flf)

    # Clean up name (replace spaces with underscores, lowercase)
    clean_name=$(echo "$base" | tr '[:upper:]' '[:lower:]' | tr ' ' '_')

    output_file="$OUTPUT_DIR/${clean_name}.bit"

    echo "[$count/$total] Converting: $base"

    # Convert the font
    if ./flf2bit "$flf_file" "$output_file" > /dev/null 2>&1; then
        success=$((success + 1))
        echo "  ✓ Saved to: $output_file"
    else
        failed=$((failed + 1))
        echo "  ✗ FAILED: $flf_file"
    fi
done

echo ""
echo "Conversion complete!"
echo "  Success: $success fonts"
echo "  Failed:  $failed fonts"
echo ""
echo "Fonts saved to: $OUTPUT_DIR"
