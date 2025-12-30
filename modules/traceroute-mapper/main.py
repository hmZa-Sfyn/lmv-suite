#!/usr/bin/env python3
"""
Traceroute Mapper Module
Custom traceroute implementation with hop visualization
"""

import os
import sys
import socket

try:
    from scapy.all import IP, ICMP, UDP, send, sr1
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

def traceroute(host, max_hops=30, timeout=2):
    """Perform traceroute"""
    try:
        ip = socket.gethostbyname(host)
    except socket.gaierror:
        return None
    
    hops = []
    for ttl in range(1, max_hops + 1):
        packet = IP(dst=ip, ttl=ttl)/ICMP()
        
        try:
            reply = sr1(packet, timeout=timeout, verbose=False)
            
            if reply is None:
                hops.append({
                    'ttl': ttl,
                    'ip': '*',
                    'host': '*',
                    'rtt': None
                })
            else:
                try:
                    host_name = socket.gethostbyaddr(reply.src)[0]
                except:
                    host_name = reply.src
                
                hops.append({
                    'ttl': ttl,
                    'ip': reply.src,
                    'host': host_name,
                    'rtt': reply.time * 1000
                })
                
                if reply.src == ip:
                    break
        except Exception as e:
            hops.append({
                'ttl': ttl,
                'ip': '*',
                'host': '*',
                'rtt': None
            })
    
    return hops

def main():
    host = os.getenv('ARG_HOST')
    max_hops = int(os.getenv('ARG_MAX_HOPS', '30'))
    timeout = int(os.getenv('ARG_TIMEOUT', '2'))
    
    if not host:
        print("Error: host required")
        sys.exit(1)
    
    print(f"[*] Tracing route to {host}")
    
    hops = traceroute(host, max_hops, timeout)
    
    if not hops:
        print("[!] Failed to resolve host")
        sys.exit(1)
    
    print(f"\n[+] Traceroute results:")
    print(f"{'TTL':<5} {'IP Address':<20} {'Host':<30} {'RTT':<10}")
    print("-" * 65)
    
    for hop in hops:
        ttl = hop['ttl']
        ip = hop['ip']
        hostname = hop['host']
        rtt = f"{hop['rtt']:.2f}ms" if hop['rtt'] else "*"
        
        print(f"{ttl:<5} {ip:<20} {hostname:<30} {rtt:<10}")
    
    print("\n[+] Complete")

if __name__ == "__main__":
    main()
