#!/bin/bash
# LanManVan CLI setup script
set -e

BIN_DIR="$HOME/bin"
LANMANVAN_DIR="$HOME/lanmanvan"
MODULES_SRC="./modules"
MODULES_DEST="$LANMANVAN_DIR/modules"

mkdir -p "$BIN_DIR" "$LANMANVAN_DIR"

# Ensure ~/bin is in PATH
for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zprofile"; do
    if [ -f "$rc" ] && ! grep -q 'export PATH="$HOME/bin:$PATH"' "$rc"; then
        echo 'export PATH="$HOME/bin:$PATH"' >> "$rc"
    fi
done

# Build binary
go mod tidy
go build -o "$BIN_DIR/lanmanvan"

# Copy modules
if [ -d "$MODULES_SRC" ]; then
    mkdir -p "$MODULES_DEST"
    rsync -a --delete "$MODULES_SRC/" "$MODULES_DEST/"
fi

# Copy VERSION file if present
for version_file in VERSION.taml VERSION.yaml VERSION.yml; do
    if [ -f "./$version_file" ]; then
        cp "./$version_file" "$LANMANVAN_DIR/"
        break
    fi
done

# Alias helper
add_or_update_alias() {
    local rc_file="$1"
    local name="$2"
    local cmd="$3"

    sed -i.bak "/alias $name=/d" "$rc_file" 2>/dev/null || true
    rm -f "$rc_file.bak" 2>/dev/null || true

    echo "alias $name='$cmd'" >> "$rc_file"
}

# Add aliases
for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zprofile"; do
    [ -f "$rc" ] || continue
    add_or_update_alias "$rc" "lanmanvan" "lanmanvan -modules $MODULES_DEST"
    add_or_update_alias "$rc" "lmvconsole" "lanmanvan -modules $MODULES_DEST"
    add_or_update_alias "$rc" "lmv_update" \
        "cd /tmp && rm -rf lanmanvan && git clone https://github.com/hmZa-Sfyn/lanmanvan && cd lanmanvan && chmod +x setup.sh && ./setup.sh"
done

echo "✔ LanManVan installed"
echo "✔ Binary: $BIN_DIR/lanmanvan"
echo "✔ Modules: $MODULES_DEST"
echo "✔ Reload shell or run: source ~/.zshrc || source ~/.bashrc"
