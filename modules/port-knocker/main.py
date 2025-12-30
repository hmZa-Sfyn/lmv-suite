#!/usr/bin/env python3
"""
Port Knocker Module
Perform port knocking sequences to open firewall ports
"""

import os
import sys
import socket
import time

def knock_tcp(host, port, timeout=1):
    """Send TCP knock to port"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        sock.connect((host, port))
        sock.close()
        return True
    except:
        return False

def knock_udp(host, port, timeout=1):
    """Send UDP knock to port"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.settimeout(timeout)
        sock.sendto(b'', (host, port))
        sock.close()
        return True
    except:
        return False

def main():
    host = os.getenv('ARG_HOST')
    sequence = os.getenv('ARG_SEQUENCE')
    delay = float(os.getenv('ARG_DELAY', '0.1'))
    protocol = os.getenv('ARG_PROTOCOL', 'tcp').lower()
    
    if not host or not sequence:
        print("Error: host and sequence required")
        sys.exit(1)
    
    try:
        ports = [int(p.strip()) for p in sequence.split(',')]
    except ValueError:
        print("Error: Invalid port sequence")
        sys.exit(1)
    
    print(f"[*] Port Knocker for {host}")
    print(f"[*] Sequence: {ports}")
    print(f"[*] Protocol: {protocol}, Delay: {delay}s")
    
    if protocol == "udp":
        knock_func = knock_udp
    else:
        knock_func = knock_tcp
    
    print(f"\n[*] Starting port knock sequence...")
    
    for i, port in enumerate(ports, 1):
        result = knock_func(host, port)
        status = "OK" if result else "TIMEOUT"
        print(f"[+] Knock {i}/{len(ports)}: Port {port} [{status}]")
        
        if i < len(ports):
            time.sleep(delay)
    
    print(f"\n[+] Knock sequence complete")
    print(f"[+] Check if port is now open on the target")

if __name__ == "__main__":
    main()
