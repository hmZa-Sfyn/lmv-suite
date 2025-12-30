#!/usr/bin/env python3
"""
Banner Grabber Module
Connect to services and grab banners to identify software versions
"""

import os
import sys
import socket
import ssl

def grab_raw_banner(host, port, timeout):
    """Grab raw banner from socket"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        sock.connect((host, port))
        
        banner = sock.recv(1024).decode('utf-8', errors='ignore').strip()
        sock.close()
        return banner
    except socket.timeout:
        return "TIMEOUT"
    except socket.error as e:
        return f"ERROR: {e}"

def grab_http_banner(host, port, timeout):
    """Grab HTTP banner"""
    try:
        import requests
        url = f"http://{host}:{port}/"
        response = requests.get(url, timeout=timeout, allow_redirects=False)
        
        server = response.headers.get('Server', 'Unknown')
        return f"Server: {server}"
    except Exception as e:
        return f"ERROR: {e}"

def grab_ftp_banner(host, port, timeout):
    """Grab FTP banner"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        sock.connect((host, port))
        
        banner = sock.recv(1024).decode('utf-8', errors='ignore').strip()
        sock.close()
        return banner
    except Exception as e:
        return f"ERROR: {e}"

def grab_ssh_banner(host, port, timeout):
    """Grab SSH banner"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        sock.connect((host, port))
        
        banner = sock.recv(1024).decode('utf-8', errors='ignore').strip()
        sock.close()
        return banner
    except Exception as e:
        return f"ERROR: {e}"

def main():
    host = os.getenv('ARG_HOST')
    port = int(os.getenv('ARG_PORT', '80'))
    timeout = int(os.getenv('ARG_TIMEOUT', '5'))
    protocol = os.getenv('ARG_PROTOCOL', 'raw').lower()
    
    if not host:
        print("Error: host required")
        sys.exit(1)
    
    print(f"[*] Grabbing banner from {host}:{port}")
    print(f"[*] Protocol: {protocol}, Timeout: {timeout}s")
    
    if protocol == "http":
        banner = grab_http_banner(host, port, timeout)
    elif protocol == "ftp":
        banner = grab_ftp_banner(host, port, timeout)
    elif protocol == "ssh":
        banner = grab_ssh_banner(host, port, timeout)
    else:
        banner = grab_raw_banner(host, port, timeout)
    
    if banner.startswith("ERROR"):
        print(f"[!] {banner}")
        sys.exit(1)
    
    print(f"\n[+] Banner:")
    print(f"    {banner}")
    
    # Try to identify software
    banner_lower = banner.lower()
    if "apache" in banner_lower:
        print("[+] Identified: Apache Web Server")
    elif "nginx" in banner_lower:
        print("[+] Identified: NGINX Web Server")
    elif "microsoft" in banner_lower or "iis" in banner_lower:
        print("[+] Identified: Microsoft IIS")
    elif "openssh" in banner_lower:
        print("[+] Identified: OpenSSH")
    elif "vsftpd" in banner_lower:
        print("[+] Identified: vsftpd")
    
    print("\n[+] Complete")

if __name__ == "__main__":
    main()
