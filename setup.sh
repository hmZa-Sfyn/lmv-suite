#!/bin/bash
# LanManVan CLI setup script (updated: dynamic module cloning)
set -e

BIN_DIR="$HOME/bin"
LANMANVAN_DIR="$HOME/lanmanvan"
MODULES_DEST="$LANMANVAN_DIR/modules"
REPO_FILE="./modules/repo_url.yaml"
LMV_MODULE_PY="./lmv_module.py"  # Assume this file exists in current dir

# Check if lmv_module.py exists in current dir
if [ ! -f "$LMV_MODULE_PY" ]; then
    echo "Error: $LMV_MODULE_PY not found in current directory! Please create it first." >&2
    exit 1
fi

mkdir -p "$BIN_DIR" "$LANMANVAN_DIR" "$MODULES_DEST"

# Ensure ~/bin is in PATH
for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zprofile"; do
    if [ -f "$rc" ] && ! grep -q 'export PATH="$HOME/bin:$PATH"' "$rc"; then
        echo 'export PATH="$HOME/bin:$PATH"' >> "$rc"
    fi
done

# Build binary
go mod tidy
go build -o "$BIN_DIR/lanmanvan"

# Check if repo_url.yaml exists
if [ ! -f "$REPO_FILE" ]; then
    echo "Error: $REPO_FILE not found in current directory!" >&2
    exit 1
fi

# Copy repo_url.yaml to LANMANVAN_DIR for later use by lmv_module.py
cp "$REPO_FILE" "$LANMANVAN_DIR/repo_url.yaml"
echo " Copied repo_url.yaml to $LANMANVAN_DIR"

# Load repo URLs from repo_url.yaml using a simple parser (supports key: "url" format)
declare -A REPOS
while IFS=":" read -r key url; do
    # Skip comments and empty lines
    [[ "$key" =~ ^[[:space:]]*# ]] && continue
    [[ -z "$key" ]] && continue

    # Trim whitespace
    key=$(echo "$key" | xargs)
    url=$(echo "$url" | sed -e 's/^ *"'// -e 's/"$//' | xargs)

    if [[ -n "$key" && -n "$url" ]]; then
        REPOS["$key"]="$url"
    fi
done < <(grep -v '^$' "$REPO_FILE" | grep ':' )

if [ ${#REPOS[@]} -eq 0 ]; then
    echo "Warning: No repositories found in $REPO_FILE"
fi

# Temporary directory for cloning
TMP_DIR="/tmp/lanmanvan_modules_$$"
mkdir -p "$TMP_DIR"

# Cleanup temp dir on exit
trap 'rm -rf "$TMP_DIR"' EXIT

# Function to print in red
red() {
    echo -e "\033[31m$*\033[0m"
}

# Ask user for each repository
for name in "${!REPOS[@]}"; do
    url="${REPOS[$name]}"
    while true; do
        read -p "Download module repo '$name' ($url)? [Y/n]: " answer
        answer=${answer:-Y}
        case "$answer" in
            [Yy]* )
                echo "Cloning $name..."
                repo_tmp="$TMP_DIR/$name"
                if git clone "$url" "$repo_tmp" 2>/dev/null; then
                    echo " Cloned $name successfully"
                else
                    red "✗ Failed to clone $url"
                    continue
                fi

                # Find all subdirectories that contain module.yaml
                find "$repo_tmp" -type f -name "module.yaml" | while read -r yaml_file; do
                    module_dir=$(dirname "$yaml_file")
                    rel_dir=$(realpath --relative-to="$repo_tmp" "$module_dir")
                    dest_dir="$MODULES_DEST/$rel_dir"

                    mkdir -p "$dest_dir"
                    rsync -a "$module_dir/" "$dest_dir/"
                    echo " Copied module: $rel_dir"
                done

                # Check if any module.yaml was found in this repo
                if ! find "$repo_tmp" -name "module.yaml" | grep -q .; then
                    red "Warning: No module.yaml found in any subdirectory of $name – nothing copied!"
                fi

                break
                ;;
            [Nn]* )
                echo "Skipping $name"
                break
                ;;
            * )
                echo "Please answer Y or n."
                ;;
        esac
    done
done

# Copy VERSION file if present
for version_file in VERSION.taml VERSION.yaml VERSION.yml; do
    if [ -f "./$version_file" ]; then
        cp "./$version_file" "$LANMANVAN_DIR/"
        break
    fi
done

# Copy lmv_module.py to LANMANVAN_DIR
cp "$LMV_MODULE_PY" "$LANMANVAN_DIR/lmv_module.py"
echo " Copied lmv_module.py to $LANMANVAN_DIR"

# Alias helper
add_or_update_alias() {
    local rc_file="$1"
    local name="$2"
    local cmd="$3"

    sed -i.bak "/alias $name=/d" "$rc_file" 2>/dev/null || true
    rm -f "$rc_file.bak" 2>/dev/null || true

    echo "alias $name='$cmd'" >> "$rc_file"
}

# Add aliases, including lmv_module to run the Python script
for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zprofile"; do
    [ -f "$rc" ] || continue
    add_or_update_alias "$rc" "lanmanvan" "lanmanvan -modules $MODULES_DEST"
    add_or_update_alias "$rc" "lmv" "lanmanvan -modules $MODULES_DEST"
    add_or_update_alias "$rc" "lmvconsole" "lanmanvan -modules $MODULES_DEST"
    add_or_update_alias "$rc" "lmv_update" \
        "cd /tmp && rm -rf lanmanvan && git clone https://github.com/hmZa-Sfyn/lanmanvan && cd lanmanvan && chmod +x setup.sh && ./setup.sh"
    add_or_update_alias "$rc" "lmv_module" "python3 $LANMANVAN_DIR/lmv_module.py \"\$@\""
done

echo " LanManVan installed successfully!"
echo " Binary: $BIN_DIR/lanmanvan"
echo " Modules directory: $MODULES_DEST"
echo " lmv_module.py copied and alias added"
echo " Reload your shell or run: source ~/.zshrc || source ~/.bashrc"