# LanManVan Framework

LanManVan is a **Metasploit-like framework** built in Go, designed to make it easy to create, manage, and execute modules. The framework supports modules written in **Python3** and **Bash**, allowing you to create powerful security tools with minimal effort.

### Partly VibeCoded: 70% Human written code, 30% AI!

## Features

✨ **Easy Module Creation** - Create modules in Python3 or Bash  
🚀 **Command-Line Interface** - Interactive shell-like interface  
📦 **Modular Design** - Load and execute modules on demand  
🔧 **Flexible Arguments** - Pass arguments like normal bash commands  
📝 **YAML Metadata** - Module configuration and documentation  
🎯 **Real-time Execution** - Execute modules with instant feedback  
🌐 **Network Tools** - Built-in examples for scanning, hashing, shells

## Installation

```bash
cd LanManVan
go mod tidy
go build -o lanmanvan main.go
```

## Quick Start

### Running the Framework

```bash
./lanmanvan
```

Or specify a custom modules directory:

```bash
./lanmanvan -modules ./my_modules
```

### Available Commands

```
hmza@0root ❯ help

╔════════════════════════════════════════════════════════════════╗
║                    📚 AVAILABLE COMMANDS                        ║
╚════════════════════════════════════════════════════════════════╝

  ❯help, h, ?              Show this help message
  ❯list, ls                List all modules
  ❯search <keyword>        Search modules by name/tag
  ❯info <module>           Show detailed module information
  ❯<module>!               Quick show module options and usage
  ❯run <module> [args]     Execute a module with arguments
  ❯<module> [args]         Shorthand: <module> arg_key=value
  ❯<module> arg_key = value Format with spaces (alternative)
  ❯env, envs               Show all global environment variables
  ❯key=value               Set global environment variable (persistent)
  ❯key=?                   View global environment variable value
  ❯create <name> [type]    Create a new module (python/bash)
  ❯edit <module>           Edit module files
  ❯delete <module>         Delete a module
  ❯history                 Show command history
  ❯clear                   Clear screen
  ❯exit, quit, q           Exit framework

hmza@0root ❯  
```

## Using Modules

### List Available Modules

```
[user@host]$ list
```

### Get Module Information

```
[user@host]$ info portscan
```

### Run a Module

```
[user@host]$ run portscan host=192.168.1.1 ports=80,443,22
```

Or use shorthand:

```
[user@host]$ portscan host=192.168.1.1 ports=80,443,22
```

## Creating Modules

### Python3 Module Structure

Create a directory under `modules/`:

```
modules/mymodule/
├── main.py          # Your Python script
└── module.yaml      # Module metadata
```

#### Python3 Module Example

**modules/mymodule/main.py:**

```python
#!/usr/bin/env python3
import os

def main():
    # Get arguments from environment variables
    target = os.getenv('ARG_TARGET') or 'localhost'
    port = os.getenv('ARG_PORT') or '80'
    
    print(f"[*] Scanning {target}:{port}")
    # Your code here
    print("[+] Scan complete!")

if __name__ == '__main__':
    main()
```

**modules/mymodule/module.yaml:**

```yaml
name: mymodule
description: "My custom module"
type: python
author: Your Name
version: 1.0.0
tags:
  - custom
  - scanning
options:
  target:
    type: string
    description: Target host
    required: true
  port:
    type: string
    description: Target port
    default: "80"
    required: false
required:
  - target
```

### Bash Module Structure

```
modules/mybashmodule/
├── main.sh          # Your Bash script
└── module.yaml      # Module metadata
```

#### Bash Module Example

**modules/mybashmodule/main.sh:**

```bash
#!/bin/bash

TARGET="${ARG_TARGET:-localhost}"
PORT="${ARG_PORT:-80}"

echo "[*] Scanning $TARGET:$PORT"
# Your code here
echo "[+] Scan complete!"
```

**modules/mybashmodule/module.yaml:**

```yaml
name: mybashmodule
description: "My bash module"
type: bash
author: Your Name
version: 1.0.0
tags:
  - custom
options:
  target:
    type: string
    description: Target host
    required: true
  port:
    type: string
    description: Target port
    default: "80"
required:
  - target
```

## Built-in Modules

### portscan
Port scanner for hosts
```
portscan host=192.168.1.1 ports=80,443,22
```

### hashgen
Generate MD5, SHA1, SHA256, SHA512 hashes
```
hashgen data="hello world"
```

### httpreq
Make HTTP requests to targets
```
httpreq host=example.com method=GET path=/
```

### revshell
Generate reverse shell payloads
```
revshell lhost=10.0.0.5 lport=4444 type=bash
```

## Module Argument Syntax

Arguments can be passed in multiple ways:

### Key=Value Format
```
run mymodule key1=value1 key2=value2
```

### Positional Arguments
```
run mymodule arg1 arg2 arg3
```

## Environment Variables

When a module executes, arguments are available as environment variables:

- `ARG_KEY` (uppercase) - For key=value arguments
- `ARG_ARG0`, `ARG_ARG1` - For positional arguments

Example:
```
portscan host=192.168.1.1
```

In Python:
```python
import os
host = os.getenv('ARG_HOST')
```

In Bash:
```bash
HOST="${ARG_HOST}"
```

## Project Structure

```
LanManVan/
├── main.go              # Entry point
├── go.mod              # Go module file
├── cli/
│   └── cli.go          # CLI implementation
├── core/
│   ├── types.go        # Type definitions
│   ├── manager.go      # Module manager
│   └── loader.go       # Module loader
├── modules/            # Modules directory
│   ├── portscan/
│   ├── hashgen/
│   ├── httpreq/
│   └── revshell/
└── README.md           # This file
```

## Tips & Tricks

### Creating Advanced Modules

1. **Use metadata extensively** - Document all options in module.yaml
2. **Error handling** - Return appropriate exit codes (0 for success, 1+ for errors)
3. **User feedback** - Use `[*]`, `[+]`, `[!]` prefixes in output for clarity
4. **Test thoroughly** - Test with various argument combinations

### Module Development Best Practices

```python
#!/usr/bin/env python3
"""
Module Description
"""

import os
import sys

def validate_args():
    """Validate required arguments"""
    required = ['ARG_TARGET', 'ARG_PORT']
    for arg in required:
        if not os.getenv(arg):
            print(f"[!] Missing required argument: {arg}")
            return False
    return True

def main():
    if not validate_args():
        sys.exit(1)
    
    target = os.getenv('ARG_TARGET')
    port = os.getenv('ARG_PORT')
    
    print(f"[*] Executing on {target}:{port}")
    try:
        # Your code
        print("[+] Success!")
    except Exception as e:
        print(f"[!] Error: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()
```

## Troubleshooting

### Module Not Found
- Ensure the module directory exists in `modules/`
- Check that module.yaml is properly formatted
- Verify the module type matches the script (python/bash)

### Module Fails to Execute
- Check Python3/Bash is installed
- Ensure scripts have execute permissions
- Verify environment variables are set correctly

### Permission Denied
```bash
chmod +x modules/*/main.py
chmod +x modules/*/main.sh
```

## Contributing

To contribute new modules:

1. Create a new directory under `modules/`
2. Add your main script (main.py or main.sh)
3. Create a module.yaml with proper metadata
4. Test thoroughly
5. Submit!

## License

MIT License - Feel free to use and modify!

## Support

For issues, questions, or contributions, feel free to reach out!

---

**Happy Hacking! 🚀**
