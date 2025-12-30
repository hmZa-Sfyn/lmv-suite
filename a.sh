#!/bin/bash

# Script: rename_modules_swap_underscore_to_dash.sh
# Purpose: Recursively rename all directories (and only directories) inside ./modules
#          by replacing underscores (_) with hyphens (-) in their names.
#          Files are left untouched.

TARGET_DIR="./modules"

echo "Starting directory rename process in $TARGET_DIR"
echo "Replacing '_' with '-' in directory names only"
echo "----------------------------------------"

# Check if the modules directory exists
if [[ ! -d "$TARGET_DIR" ]]; then
    echo "Error: Directory '$TARGET_DIR' not found!"
    exit 1
fi

# Find all directories inside ./modules (any depth) and process them bottom-up
# Processing bottom-up prevents issues when renaming parent dirs before children
find "$TARGET_DIR" -depth -type d | while read -r dir; do
    # Skip the base ./modules directory itself to avoid unnecessary rename
    [[ "$dir" == "$TARGET_DIR" ]] && continue

    # Get the parent directory and the current name
    parent="$(dirname "$dir")"
    current_name="$(basename "$dir")"

    # Create new name: replace all underscores with hyphens
    new_name="${current_name//_/-}"

    # If the name already has hyphens only (or no change needed), skip
    if [[ "$current_name" == "$new_name" ]]; then
        continue
    fi

    # Full new path
    new_dir="$parent/$new_name"

    # Check if the new directory name already exists (avoid conflicts)
    if [[ -e "$new_dir" ]]; then
        echo "SKIP (conflict): '$dir' → '$new_dir' (already exists)"
        continue
    fi

    # Perform the rename
    echo "RENAMING: '$dir' → '$new_dir'"
    mv "$dir" "$new_dir"
done

echo "----------------------------------------"
echo "All directory renames completed!"
echo "Note: Only directories were renamed. Files inside remain unchanged."