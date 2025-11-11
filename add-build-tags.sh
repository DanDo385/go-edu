#!/bin/bash

cd /home/user/go-edu

for file in minis/*/exercise/exercise.go; do
    if [ -f "$file" ]; then
        # Check if already has build tag
        if ! head -1 "$file" | grep -q "go:build"; then
            # Create temp file
            tmpfile=$(mktemp)
            # Add build tag header
            echo "//go:build !solution" > "$tmpfile"
            echo "// +build !solution" >> "$tmpfile"
            echo "" >> "$tmpfile"
            # Append original content
            cat "$file" >> "$tmpfile"
            # Replace original
            mv "$tmpfile" "$file"
            echo "✓ Added build tag to: $file"
        else
            echo "  Already has build tag: $file"
        fi
    fi
done

echo ""
echo "Done! All exercise.go files now have build tags."
