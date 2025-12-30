#!/usr/bin/env python3
"""
Network Device Detector Module
Detect and identify devices on network with OS fingerprinting
"""

import os
import sys
import socket
import threading
from threading import Lock
from queue import Queue
from ipaddress import ip_network

try:
    from scapy.all import ICMP, IP, sr1, TCP
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

discovered = []
devices_lock = Lock()

def ping_host(ip, timeout):
    """Ping a host"""
    try:
        packet = IP(dst=str(ip))/ICMP()
        reply = sr1(packet, timeout=timeout, verbose=False)
        
        if reply:
            return True
    except:
        pass
    
    return False

def get_hostname(ip):
    """Get hostname from IP"""
    try:
        return socket.gethostbyaddr(str(ip))[0]
    except:
        return "unknown"

def fingerprint_os(ip):
    """Simple OS fingerprinting based on TTL"""
    try:
        packet = IP(dst=str(ip))/ICMP()
        reply = sr1(packet, timeout=2, verbose=False)
        
        if reply:
            ttl = reply.ttl
            
            # TTL-based fingerprinting
            if ttl == 64:
                return "Linux/Unix"
            elif ttl == 128:
                return "Windows"
            elif ttl == 255:
                return "Router/Network Device"
            elif ttl <= 32:
                return "Unknown"
            else:
                return "Unknown"
    except:
        pass
    
    return "Unknown"

def worker(queue, timeout, fingerprint):
    """Worker thread"""
    while True:
        ip = queue.get()
        if ip is None:
            break
        
        if ping_host(ip, timeout):
            hostname = get_hostname(ip)
            os_type = fingerprint_os(ip) if fingerprint else "N/A"
            
            with devices_lock:
                discovered.append({
                    'ip': str(ip),
                    'hostname': hostname,
                    'os': os_type
                })
        
        queue.task_done()

def main():
    network = os.getenv('ARG_NETWORK')
    threads = int(os.getenv('ARG_THREADS', '20'))
    timeout = int(os.getenv('ARG_TIMEOUT', '2'))
    fingerprint = os.getenv('ARG_FINGERPRINT', 'true').lower() == 'true'
    
    if not network:
        print("Error: network required")
        sys.exit(1)
    
    print(f"[*] Network Device Detector")
    print(f"[*] Network: {network}")
    print(f"[*] Threads: {threads}, Timeout: {timeout}s")
    print(f"[*] OS Fingerprinting: {'Yes' if fingerprint else 'No'}")
    
    try:
        net = ip_network(network, strict=False)
        hosts = list(net.hosts())
    except ValueError:
        print(f"Error: Invalid network: {network}")
        sys.exit(1)
    
    print(f"[*] Scanning {len(hosts)} hosts...\n")
    
    queue = Queue()
    thread_list = []
    
    # Start workers
    for i in range(threads):
        t = threading.Thread(target=worker, args=(queue, timeout, fingerprint))
        t.start()
        thread_list.append(t)
    
    # Queue hosts
    for host in hosts:
        queue.put(host)
    
    # Wait completion
    queue.join()
    
    # Stop workers
    for i in range(threads):
        queue.put(None)
    for t in thread_list:
        t.join()
    
    if not discovered:
        print(f"[!] No devices found")
        sys.exit(0)
    
    print(f"[+] Found {len(discovered)} device(s):")
    print(f"\n{'#':<3} {'IP Address':<20} {'Hostname':<30} {'OS':<20}")
    print("-" * 75)
    
    for i, device in enumerate(sorted(discovered, key=lambda x: x['ip']), 1):
        hostname = device['hostname'][:27] if len(device['hostname']) > 27 else device['hostname']
        print(f"{i:<3} {device['ip']:<20} {hostname:<30} {device['os']:<20}")
    
    print(f"\n[+] OS Summary:")
    os_count = {}
    for device in discovered:
        os_type = device['os']
        os_count[os_type] = os_count.get(os_type, 0) + 1
    
    for os_type, count in sorted(os_count.items(), key=lambda x: x[1], reverse=True):
        print(f"    {os_type}: {count}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
