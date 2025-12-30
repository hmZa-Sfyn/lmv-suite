#!/usr/bin/env python3
"""
Device Fingerprinter Module
Fingerprint devices by analyzing services and behavior
"""

import os
import sys
import socket
import re

def scan_ports(host, port_range, timeout):
    """Scan open ports"""
    ports = []
    
    # Parse port range
    if '-' in port_range:
        start, end = port_range.split('-')
        start, end = int(start), int(end)
    else:
        start = end = int(port_range)
    
    common_ports = [21, 22, 23, 25, 53, 80, 110, 143, 443, 465, 587, 993, 995, 
                   3306, 3389, 5432, 5900, 8080, 8443, 9200, 27017]
    
    ports_to_check = [p for p in common_ports if start <= p <= end]
    
    for port in ports_to_check:
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(timeout)
            result = sock.connect_ex((host, port))
            sock.close()
            
            if result == 0:
                ports.append(port)
        except:
            pass
    
    return ports

def get_banner(host, port, timeout):
    """Get service banner"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        sock.connect((host, port))
        
        banner = sock.recv(1024).decode('utf-8', errors='ignore').strip()
        sock.close()
        
        return banner
    except:
        return None

def identify_service(port, banner=None):
    """Identify service from port and banner"""
    common_services = {
        21: 'FTP',
        22: 'SSH',
        23: 'Telnet',
        25: 'SMTP',
        53: 'DNS',
        80: 'HTTP',
        110: 'POP3',
        143: 'IMAP',
        443: 'HTTPS',
        465: 'SMTPS',
        587: 'SMTP',
        993: 'IMAPS',
        995: 'POP3S',
        3306: 'MySQL',
        3389: 'RDP',
        5432: 'PostgreSQL',
        5900: 'VNC',
        8080: 'HTTP Alt',
        8443: 'HTTPS Alt',
        9200: 'Elasticsearch',
        27017: 'MongoDB'
    }
    
    service = common_services.get(port, 'Unknown')
    
    if banner:
        # Try to identify from banner
        if 'Apache' in banner:
            service = 'Apache HTTP'
        elif 'nginx' in banner:
            service = 'Nginx HTTP'
        elif 'Microsoft' in banner or 'IIS' in banner:
            service = 'Microsoft IIS'
        elif 'OpenSSH' in banner:
            service = 'OpenSSH'
        elif 'MySQL' in banner:
            service = 'MySQL'
        elif 'PostgreSQL' in banner:
            service = 'PostgreSQL'
    
    return service

def fingerprint_os(ports):
    """Guess OS from open ports"""
    if 3389 in ports and 445 in ports:
        return "Windows"
    elif 22 in ports and 80 in ports:
        return "Linux/Unix"
    elif 22 in ports:
        return "Unix-like"
    else:
        return "Unknown"

def main():
    target = os.getenv('ARG_TARGET')
    port_range = os.getenv('ARG_PORT_RANGE', '1-1000')
    timeout = int(os.getenv('ARG_TIMEOUT', '5'))
    
    if not target:
        print("Error: target required")
        sys.exit(1)
    
    print(f"[*] Device Fingerprinter")
    print(f"[*] Target: {target}")
    print(f"[*] Port Range: {port_range}")
    
    # Resolve hostname
    try:
        ip = socket.gethostbyname(target)
        print(f"[*] Resolved to: {ip}")
    except socket.gaierror:
        print(f"[!] Could not resolve: {target}")
        sys.exit(1)
    
    print(f"[*] Scanning ports...\n")
    
    open_ports = scan_ports(ip, port_range, timeout)
    
    if not open_ports:
        print(f"[!] No open ports found")
        sys.exit(0)
    
    print(f"[+] Found {len(open_ports)} open port(s):")
    print(f"\n{'Port':<8} {'Service':<25} {'Banner':<40}")
    print("-" * 75)
    
    services = []
    for port in sorted(open_ports):
        banner = get_banner(ip, port, timeout)
        service = identify_service(port, banner)
        services.append(service)
        
        banner_display = banner[:37] if banner else "N/A"
        print(f"{port:<8} {service:<25} {banner_display:<40}")
    
    # Fingerprint OS
    guessed_os = fingerprint_os(open_ports)
    
    print(f"\n[+] Fingerprinting Analysis:")
    print(f"    Open Ports: {open_ports}")
    print(f"    Services: {', '.join(set(services))}")
    print(f"    Guessed OS: {guessed_os}")
    
    # Device classification
    if 'MySQL' in services or 'PostgreSQL' in services:
        device_type = "Database Server"
    elif 'HTTP' in services or 'HTTPS' in services:
        device_type = "Web Server"
    elif 'SSH' in services:
        device_type = "Linux/Unix Server"
    elif 'RDP' in services:
        device_type = "Windows Server"
    else:
        device_type = "Network Device"
    
    print(f"    Device Type: {device_type}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
