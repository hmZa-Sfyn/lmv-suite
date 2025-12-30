#!/usr/bin/env python3
"""Generate README.md for all modules"""

import os
import yaml
from pathlib import Path

MODULES_DIR = Path(__file__).parent / "modules"

def generate_readme(module_path):
    """Generate README.md for a module"""
    yaml_file = module_path / "module.yaml"
    readme_file = module_path / "README.md"
    
    if not yaml_file.exists():
        return False
    
    # Parse YAML
    try:
        with open(yaml_file, 'r') as f:
            config = yaml.safe_load(f)
    except yaml.YAMLError:
        # Fallback for malformed YAML
        config = {'name': module_path.name, 'description': 'Module', 'type': 'unknown', 'author': 'unknown', 'version': '1.0.0', 'tags': []}
    
    if not config:
        return False
    
    # Generate README content
    content = f"""# {config.get('name', module_path.name)}

## Description
{config.get('description', 'No description available')}

## Metadata
- **Type:** {config.get('type', 'unknown')}
- **Author:** {config.get('author', 'unknown')}
- **Version:** {config.get('version', '1.0.0')}

## Tags
{', '.join(config.get('tags', [])) if config.get('tags') else 'No tags'}

## Options
"""
    
    options = config.get('options', {})
    if options:
        for opt_name, opt_config in options.items():
            content += f"\n### {opt_name}\n"
            content += f"- **Type:** {opt_config.get('type', 'string')}\n"
            content += f"- **Description:** {opt_config.get('description', 'No description')}\n"
            content += f"- **Required:** {'Yes' if opt_config.get('required', False) else 'No'}\n"
            if 'default' in opt_config:
                content += f"- **Default:** {opt_config['default']}\n"
    else:
        content += "\nNo options available.\n"
    
    # Write README
    with open(readme_file, 'w') as f:
        f.write(content)
    
    return True

def main():
    """Generate READMEs for all modules"""
    if not MODULES_DIR.exists():
        print(f"Modules directory not found: {MODULES_DIR}")
        return
    
    count = 0
    for module_path in sorted(MODULES_DIR.iterdir()):
        if module_path.is_dir():
            if generate_readme(module_path):
                count += 1
                print(f"✓ {module_path.name}")
    
    print(f"\nGenerated {count} README files")

if __name__ == "__main__":
    main()
