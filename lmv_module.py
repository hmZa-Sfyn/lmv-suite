import argparse
import os
import subprocess
import shutil
import glob
import sys
import yaml

HOME = os.path.expanduser("~")
LANMANVAN_DIR = os.path.join(HOME, "lanmanvan")
MODULES_DIR = os.path.join(LANMANVAN_DIR, "modules")
REPO_FILE = os.path.join(LANMANVAN_DIR, "repo_url.yaml")

def load_repos():
    if os.path.exists(REPO_FILE):
        with open(REPO_FILE, "r") as f:
            repos = yaml.safe_load(f) or {}
        return repos
    return {}

def clone_repo(url, tmp_dir):
    try:
        subprocess.run(["git", "clone", "--quiet", url, tmp_dir], check=True)
        return True
    except subprocess.CalledProcessError:
        return False

def main():
    parser = argparse.ArgumentParser(description="LanManVan Module Manager")
    subparsers = parser.add_subparsers(dest="command")

    install_parser = subparsers.add_parser("install", help="Install modules from repo")
    install_parser.add_argument("repo", help="Repo name or URL (repo=<name|url>)")

    remove_parser = subparsers.add_parser("remove", help="Remove modules by pattern")
    remove_parser.add_argument("name", help="Module name pattern (name=<pattern>)")

    subparsers.add_parser("list", help="List installed modules")

    subparsers.add_parser("help", help="Show help")

    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(0)

    args = parser.parse_args()

    if args.command == "help":
        parser.print_help()
        sys.exit(0)

    if args.command == "list":
        modules = [d for d in os.listdir(MODULES_DIR) if os.path.isdir(os.path.join(MODULES_DIR, d))]
        if not modules:
            print("[!] No modules installed")
        else:
            print("Installed modules:")
            for i, mod in enumerate(sorted(modules), 1):
                print(f"  {i}. {mod}")
        sys.exit(0)

    if args.command == "install":
        repo_input = args.repo.replace("repo=", "") if args.repo.startswith("repo=") else args.repo
        repos = load_repos()
        url = repos.get(repo_input) if repo_input in repos else repo_input
        if not url.startswith("http"):
            print("[✗] Invalid repo or URL")
            sys.exit(1)

        name = url.split("/")[-1].replace(".git", "")
        tmp_dir = os.path.join("/tmp", f"lmv_module_{os.getpid()}_{name}")
        os.makedirs(tmp_dir, exist_ok=True)

        print(f"[+] Cloning {url}...")
        if not clone_repo(url, tmp_dir):
            print("[✗] Failed to clone repo")
            shutil.rmtree(tmp_dir)
            sys.exit(1)

        count = 0
        for mod_dir in glob.glob(os.path.join(tmp_dir, "*/")):
            mod_name = os.path.basename(os.path.normpath(mod_dir))
            if mod_name == ".git":
                continue
            dest = os.path.join(MODULES_DIR, mod_name)
            shutil.copytree(mod_dir, dest, dirs_exist_ok=True)
            print(f"[+] Installed/Updated: {mod_name}")
            count += 1

        shutil.rmtree(tmp_dir)
        print(f"[+] Installed {count} modules from {name}")
        sys.exit(0)

    if args.command == "remove":
        pattern = args.name.replace("name=", "") if args.name.startswith("name=") else args.name
        matched = glob.glob(os.path.join(MODULES_DIR, pattern))
        if not matched:
            print(f"[!] No modules matched: {pattern}")
            sys.exit(0)

        mods = [os.path.basename(d) for d in matched if os.path.isdir(d)]
        print("Will remove:")
        for mod in mods:
            print(f"  - {mod}")
        confirm = input("Confirm? [y/N]: ").strip().lower()
        if confirm != "y":
            print("[+] Cancelled")
            sys.exit(0)

        for mod in mods:
            shutil.rmtree(os.path.join(MODULES_DIR, mod))
            print(f"[+] Removed: {mod}")
        sys.exit(0)

    parser.print_help()

if __name__ == "__main__":
    main()