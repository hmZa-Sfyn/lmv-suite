#!/bin/bash

set -e

echo "=================================================================="
echo "=============== LanManVan Setup Script STARTED ===================="
echo "=============== (Excessive Logging Edition) ======================="
echo "=============== Current time: $(date) ============================"
echo "=============== User: $USER | Home: $HOME | Shell: $SHELL ========"
echo "=================================================================="
echo ""

BIN_DIR="$HOME/bin"
LANMANVAN_DIR="$HOME/lanmanvan"
MODULES_DEST="$LANMANVAN_DIR/modules"
REPO_FILE="./modules/repo_url.yaml"

echo "[INFO] Creating required directories..."
echo "[INFO] BIN_DIR set to: $BIN_DIR"
echo "[INFO] LANMANVAN_DIR set to: $LANMANVAN_DIR"
echo "[INFO] MODULES_DEST set to: $MODULES_DEST"
mkdir -p "$BIN_DIR" && echo "[SUCCESS] Created/verified $BIN_DIR"
mkdir -p "$LANMANVAN_DIR" && echo "[SUCCESS] Created/verified $LANMANVAN_DIR"
mkdir -p "$MODULES_DEST" && echo "[SUCCESS] Created/verified $MODULES_DEST"
echo "[INFO] All directories ready."
echo ""

echo "[INFO] Ensuring ~/bin is in PATH for common shell rc files..."
for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zprofile"; do
    echo "[CHECK] Examining rc file: $rc"
    if [ -f "$rc" ]; then
        echo "[FOUND] $rc exists"
        if grep -q 'export PATH="$HOME/bin:$PATH"' "$rc"; then
            echo "[SKIP] PATH already includes ~/bin in $rc"
        else
            echo "[ADD] Appending PATH export to $rc"
            echo 'export PATH="$HOME/bin:$PATH"' >> "$rc"
            echo "[SUCCESS] PATH export added to $rc"
        fi
    else
        echo "[MISSING] $rc does not exist – skipping"
    fi
done
echo "[INFO] PATH configuration complete."
echo ""

echo "[INFO] Running Go module tidy..."
go mod tidy && echo "[SUCCESS] go mod tidy completed successfully"
echo "[INFO] Building LanManVan binary..."
go build -o "$BIN_DIR/lanmanvan"
if [ $? -eq 0 ]; then
    echo "[SUCCESS] Binary built and placed at $BIN_DIR/lanmanvan"
    echo "[INFO] You can now run 'lanmanvan' from anywhere (after shell reload)"
else
    echo "[ERROR] Binary build failed!" >&2
    exit 1
fi
echo ""

echo "[CHECK] Looking for repo_url.yaml file at: $REPO_FILE"
if [ ! -f "$REPO_FILE" ]; then
    echo "[ERROR] $REPO_FILE not found in current directory!" >&2
    echo "[FATAL] Cannot proceed without repository configuration."
    exit 1
else
    echo "[SUCCESS] Found $REPO_FILE"
fi
echo ""

echo "[INFO] Parsing $REPO_FILE for repository definitions..."
declare -A REPOS
while IFS=":" read -r key url; do
    echo "[PARSE] Raw line: key='$key' url='$url'"
    [[ "$key" =~ ^[[:space:]]*# ]] && { echo "[SKIP] Comment line"; continue; }
    [[ -z "$key" ]] && { echo "[SKIP] Empty line"; continue; }

    key=$(echo "$key" | xargs)
    url=$(echo "$url" | sed -e 's/^ *"'// -e 's/"$//' | xargs)

    echo "[PARSED] Cleaned: key='$key' url='$url'"
    if [[ -n "$key" && -n "$url" ]]; then
        REPOS["$key"]="$url"
        echo "[ADDED] Repository '$key' → '$url'"
    else
        echo "[WARN] Incomplete entry – ignored"
    fi
done < <(grep -v '^$' "$REPO_FILE" | grep ':' )

echo "[INFO] Finished parsing. Found ${#REPOS[@]} repositories."
if [ ${#REPOS[@]} -eq 0 ]; then
    echo "[WARNING] No repositories found in $REPO_FILE – module installation will be empty."
else
    echo "[INFO] Available repositories:"
    for name in "${!REPOS[@]}"; do
        echo "    • $name → ${REPOS[$name]}"
    done
fi
echo ""

TMP_DIR="/tmp/lanmanvan_modules_$$"
echo "[INFO] Creating temporary directory: $TMP_DIR"
mkdir -p "$TMP_DIR" && echo "[SUCCESS] Temp dir created"
trap 'echo "[CLEANUP] Removing temporary directory $TMP_DIR"; rm -rf "$TMP_DIR"' EXIT
echo "[INFO] Cleanup trap installed."
echo ""

# Function for red text
red() { echo -e "\033[31m[ERROR] $*\033[0m"; }

echo "=================================================================="
echo "=============== MODULE DOWNLOAD PHASE BEGIN ======================"
echo "=================================================================="

for name in "${!REPOS[@]}"; do
    url="${REPOS[$name]}"
    echo ""
    echo "[MODULE] Processing repository: $name"
    echo "[URL] $url"
    while true; do
        read -p "[QUESTION] Download module repo '$name' ($url)? [Y/n]: " answer
        answer=${answer:-Y}
        echo "[INPUT] User answered: $answer"
        case "$answer" in
            [Yy]* )
                echo "[ACTION] User approved – starting clone..."
                repo_tmp="$TMP_DIR/$name"
                echo "[CLONE] Cloning into $repo_tmp"
                if git clone "$url" "$repo_tmp" >"$repo_tmp.clone.log" 2>&1; then
                    echo "[SUCCESS] Cloned $name successfully"
                    cat "$repo_tmp.clone.log" | sed 's/^/[GIT] /'
                else
                    echo "[FAIL] Clone failed – details below:"
                    red "✗ Failed to clone $url"
                    cat "$repo_tmp.clone.log" | sed 's/^/[GIT-ERR] /' >&2
                    rm -f "$repo_tmp.clone.log"
                    continue
                fi
                rm -f "$repo_tmp.clone.log"

                echo "[SEARCH] Looking for module.yaml files in cloned repo..."
                module_count=0
                find "$repo_tmp" -type f -name "module.yaml" | while read -r yaml_file; do
                    ((module_count++))
                    module_dir=$(dirname "$yaml_file")
                    rel_dir=$(realpath --relative-to="$repo_tmp" "$module_dir")
                    dest_dir="$MODULES_DEST/$rel_dir"

                    echo "[FOUND] module.yaml at $yaml_file"
                    echo "[COPY] Source: $module_dir/ → Destination: $dest_dir"
                    mkdir -p "$dest_dir" && echo "[MKDIR] Created $dest_dir"
                    rsync -a --log-format="[RSYNC] %o %n" "$module_dir/" "$dest_dir/" | sed "s/^/[RSYNC] /"
                    echo "[SUCCESS] Copied module: $rel_dir"
                done

                if [ "$module_count" -eq 0 ]; then
                    red "Warning: No module.yaml found in any subdirectory of $name – nothing copied!"
                else
                    echo "[SUMMARY] Copied $module_count module(s) from $name"
                fi

                break
                ;;
            [Nn]* )
                echo "[ACTION] User declined – skipping $name"
                break
                ;;
            * )
                echo "[INVALID] Please answer Y or n."
                ;;
        esac
    done
done

echo ""
echo "=================================================================="
echo "=============== MODULE DOWNLOAD PHASE COMPLETE ==================="
echo "=================================================================="

echo "[INFO] Checking for VERSION file to copy..."
for version_file in VERSION.taml VERSION.yaml VERSION.yml; do
    if [ -f "./$version_file" ]; then
        echo "[FOUND] $version_file – copying to $LANMANVAN_DIR/"
        cp "./$version_file" "$LANMANVAN_DIR/" && echo "[SUCCESS] Version file copied"
        break
    else
        echo "[MISSING] $version_file not present"
    fi
done
echo ""

# Alias helper with excessive logging
add_or_update_alias() {
    local rc_file="$1"
    local name="$2"
    local cmd="$3"
    echo "[ALIAS] Processing alias '$name' in $rc_file"

    if sed -i.bak "/alias $name=/d" "$rc_file" 2>/dev/null; then
        echo "[ALIAS] Removed old definition of '$name'"
        rm -f "$rc_file.bak" && echo "[CLEAN] Removed backup file"
    else
        echo "[ALIAS] No previous definition found or sed failed gracefully"
    fi

    echo "[ALIAS] Adding new alias: alias $name='$cmd'"
    echo "alias $name='$cmd'" >> "$rc_file"
    echo "[SUCCESS] Alias '$name' updated in $rc_file"
}

echo "[INFO] Installing shell aliases (lanmanvan, lmvconsole, lmv_update, lmv_module)..."
for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zprofile"; do
    [ -f "$rc" ] || { echo "[SKIP] $rc does not exist"; continue; }
    echo "[PROCESS] Updating aliases in $rc"
    add_or_update_alias "$rc" "lanmanvan" "lanmanvan -modules $MODULES_DEST"
    add_or_update_alias "$rc" "lmvconsole" "lanmanvan -modules $MODULES_DEST"
    add_or_update_alias "$rc" "lmv_update" \
        "cd /tmp && rm -rf lanmanvan && git clone https://github.com/hmZa-Sfyn/lanmanvan && cd lanmanvan && chmod +x setup.sh && ./setup.sh"
    add_or_update_alias "$rc" "lmv_module" "lanmanvan -modules $MODULES_DEST module"
done
echo "[INFO] All aliases installed."
echo ""
echo "Available aliases:"
echo "  • lanmanvan     → lanmanvan -modules $MODULES_DEST"
echo "  • lmvconsole    → same as above (console entrypoint)"
echo "  • lmv_update    → pull and reinstall latest LanManVan"
echo "  • lmv_module    → module management (install/remove/etc)"
echo ""

echo "=================================================================="
echo "===================== INSTALLATION COMPLETE ======================"
echo "=================================================================="
echo "✔ LanManVan installed successfully with EXCESSIVE logging!"
echo "✔ Binary location      : $BIN_DIR/lanmanvan"
echo "✔ Modules directory    : $MODULES_DEST"
echo "✔ To start using immediately, run:"
echo "      source ~/.zshrc 2>/dev/null || source ~/.bashrc"
echo "✔ Module management: lmv_module install repo=https://github.com/Lanmanvan-Org/basic81"
echo "=================================================================="
echo "Thank you for your patience through all that logging! 🚀"
echo ""