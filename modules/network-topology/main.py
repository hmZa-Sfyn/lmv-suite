#!/usr/bin/env python3
"""
Network Topology Module
Discover and map local network devices
"""

import os
import sys
import socket
import threading
from threading import Lock
from queue import Queue
from ipaddress import ip_network

try:
    from scapy.all import ICMP, IP, sr1
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

discovered_hosts = []
hosts_lock = Lock()

def ping_host(ip, timeout):
    """Ping a single host"""
    try:
        packet = IP(dst=str(ip))/ICMP()
        reply = sr1(packet, timeout=timeout, verbose=False)
        
        if reply:
            try:
                hostname = socket.gethostbyaddr(str(ip))[0]
            except:
                hostname = "unknown"
            
            return True, hostname
    except:
        pass
    
    return False, None

def worker(queue, timeout):
    """Worker thread"""
    while True:
        ip = queue.get()
        if ip is None:
            break
        
        found, hostname = ping_host(ip, timeout)
        if found:
            with hosts_lock:
                discovered_hosts.append({
                    'ip': str(ip),
                    'hostname': hostname
                })
        
        queue.task_done()

def main():
    network = os.getenv('ARG_NETWORK')
    threads = int(os.getenv('ARG_THREADS', '10'))
    timeout = int(os.getenv('ARG_TIMEOUT', '2'))
    
    if not network:
        print("Error: network required")
        sys.exit(1)
    
    print(f"[*] Network Topology Discovery")
    print(f"[*] Network: {network}")
    print(f"[*] Threads: {threads}, Timeout: {timeout}s")
    
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
        t = threading.Thread(target=worker, args=(queue, timeout))
        t.start()
        thread_list.append(t)
    
    # Queue all hosts
    for host in hosts:
        queue.put(host)
    
    # Wait for completion
    queue.join()
    
    # Stop workers
    for i in range(threads):
        queue.put(None)
    for t in thread_list:
        t.join()
    
    print(f"\n[+] Found {len(discovered_hosts)} active hosts:")
    print(f"\n{'IP Address':<20} {'Hostname':<30}")
    print("-" * 50)
    
    for host in sorted(discovered_hosts, key=lambda x: x['ip']):
        print(f"{host['ip']:<20} {host['hostname']:<30}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
