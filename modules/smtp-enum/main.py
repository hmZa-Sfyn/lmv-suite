#!/usr/bin/env python3
"""
SMTP Enum Module
Enumerate valid email users on SMTP servers
"""

import os
import sys
import socket
import time

def enum_vrfy(host, port, username, timeout=5):
    """Enumerate user using VRFY"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        sock.connect((host, port))
        
        sock.recv(1024)  # Receive greeting
        sock.send(f"VRFY {username}\r\n".encode())
        
        response = sock.recv(1024).decode('utf-8', errors='ignore')
        sock.close()
        
        # 250 response means user exists
        return response.startswith('250')
    except:
        return False

def enum_expn(host, port, username, timeout=5):
    """Enumerate user using EXPN"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        sock.connect((host, port))
        
        sock.recv(1024)  # Receive greeting
        sock.send(f"EXPN {username}\r\n".encode())
        
        response = sock.recv(1024).decode('utf-8', errors='ignore')
        sock.close()
        
        return response.startswith('250')
    except:
        return False

def enum_rcpt(host, port, username, domain, timeout=5):
    """Enumerate user using RCPT"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        sock.connect((host, port))
        
        sock.recv(1024)  # Greeting
        sock.send(b"MAIL FROM:<test@test.com>\r\n")
        sock.recv(1024)
        
        sock.send(f"RCPT TO:<{username}@{domain}>\r\n".encode())
        response = sock.recv(1024).decode('utf-8', errors='ignore')
        sock.close()
        
        # 250 or 251 means user exists
        return response.startswith('25')
    except:
        return False

def main():
    host = os.getenv('ARG_HOST')
    port = int(os.getenv('ARG_PORT', '25'))
    wordlist = os.getenv('ARG_WORDLIST')
    method = os.getenv('ARG_METHOD', 'rcpt').lower()
    
    if not host or not wordlist:
        print("Error: host and wordlist required")
        sys.exit(1)
    
    # Extract domain from host for RCPT method
    domain = host.split('@')[-1] if '@' in host else host
    
    print(f"[*] SMTP Enumeration")
    print(f"[*] Server: {host}:{port}")
    print(f"[*] Method: {method}")
    print(f"[*] Wordlist: {wordlist}")
    
    if not os.path.exists(wordlist):
        print(f"[!] Wordlist not found: {wordlist}")
        print(f"[*] Creating sample wordlist...")
        users = ["admin", "info", "support", "test", "postmaster"]
    else:
        with open(wordlist, 'r') as f:
            users = [line.strip() for line in f if line.strip()]
    
    print(f"[*] Testing {len(users)} users...\n")
    
    found_users = []
    
    # Select enumeration method
    if method == "vrfy":
        enum_func = enum_vrfy
    elif method == "expn":
        enum_func = enum_expn
    else:  # rcpt
        enum_func = lambda h, p, u, t: enum_rcpt(h, p, u, domain, t)
    
    for i, user in enumerate(users, 1):
        print(f"[*] Testing {i}/{len(users)}: {user}", end='\r')
        
        if enum_func(host, port, user):
            found_users.append(user)
            print(f"[+] Found: {user:<30}")
    
    print(f"\n\n[+] Complete!")
    print(f"[+] Found {len(found_users)} valid users:")
    for user in found_users:
        print(f"    - {user}")

if __name__ == "__main__":
    main()
