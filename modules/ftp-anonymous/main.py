#!/usr/bin/env python3
"""
FTP Anonymous Module
Check for anonymous FTP login and list directories
"""

import os
import sys

try:
    from ftplib import FTP, all_errors
except ImportError:
    print("Error: ftplib not available")
    sys.exit(1)

def check_anonymous_ftp(host, port, timeout):
    """Check if anonymous FTP is enabled"""
    try:
        ftp = FTP(timeout=timeout)
        ftp.connect(host, port)
        ftp.login('anonymous', 'anonymous@example.com')
        return True, ftp
    except all_errors:
        return False, None

def list_directory(ftp, path="/", depth=0, max_depth=2, results=None):
    """Recursively list FTP directories"""
    if results is None:
        results = []
    
    if depth > max_depth:
        return results
    
    try:
        files = ftp.nlst(path)
        
        for file in files:
            try:
                # Try to change to directory
                ftp.cwd(file)
                indent = "  " * depth
                results.append({
                    'type': 'dir',
                    'path': file,
                    'depth': depth
                })
                ftp.cwd('..')
                
                # Recurse
                list_directory(ftp, file, depth + 1, max_depth, results)
            except:
                # It's a file
                indent = "  " * depth
                results.append({
                    'type': 'file',
                    'path': file,
                    'depth': depth
                })
    except:
        pass
    
    return results

def main():
    host = os.getenv('ARG_HOST')
    port = int(os.getenv('ARG_PORT', '21'))
    timeout = int(os.getenv('ARG_TIMEOUT', '5'))
    depth = int(os.getenv('ARG_DEPTH', '2'))
    
    if not host:
        print("Error: host required")
        sys.exit(1)
    
    print(f"[*] FTP Anonymous Checker")
    print(f"[*] Target: {host}:{port}")
    print(f"[*] Timeout: {timeout}s")
    
    # Check for anonymous access
    print(f"\n[*] Checking for anonymous FTP access...")
    
    success, ftp = check_anonymous_ftp(host, port, timeout)
    
    if not success:
        print(f"[!] Anonymous FTP not available")
        sys.exit(1)
    
    print(f"[+] Anonymous FTP is ENABLED!")
    
    # Get banner
    try:
        banner = ftp.getwelcome()
        print(f"\n[+] Server Banner:")
        for line in banner.split('\n'):
            if line.strip():
                print(f"    {line}")
    except:
        pass
    
    # List directories
    print(f"\n[*] Listing directories (depth: {depth})...")
    
    try:
        contents = list_directory(ftp, "/", 0, depth)
        
        print(f"\n[+] Found {len(contents)} items:")
        for item in contents:
            indent = "  " * item['depth']
            icon = "[DIR]" if item['type'] == 'dir' else "[FILE]"
            print(f"{indent}{icon} {item['path']}")
    
    except Exception as e:
        print(f"[!] Error listing: {e}")
    
    finally:
        try:
            ftp.quit()
        except:
            pass
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
