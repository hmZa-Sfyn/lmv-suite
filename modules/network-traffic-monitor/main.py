#!/usr/bin/env python3
"""
Network Traffic Monitor Module
Monitor and analyze network traffic in real-time
"""

import os
import sys

try:
    from scapy.all import sniff, IP, TCP, UDP, ICMP
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

stats = {
    'total_packets': 0,
    'tcp': 0,
    'udp': 0,
    'icmp': 0,
    'other': 0,
    'bytes': 0,
    'src_ips': set(),
    'dst_ips': set()
}

def packet_callback(packet):
    """Analyze packets"""
    global stats
    
    stats['total_packets'] += 1
    stats['bytes'] += len(packet)
    
    if packet.haslayer('IP'):
        ip = packet['IP']
        stats['src_ips'].add(ip.src)
        stats['dst_ips'].add(ip.dst)
        
        if packet.haslayer('TCP'):
            stats['tcp'] += 1
        elif packet.haslayer('UDP'):
            stats['udp'] += 1
        elif packet.haslayer('ICMP'):
            stats['icmp'] += 1
        else:
            stats['other'] += 1
    else:
        stats['other'] += 1

def main():
    interface = os.getenv('ARG_INTERFACE')
    duration = int(os.getenv('ARG_DURATION', '60'))
    bpf_filter = os.getenv('ARG_FILTER')
    packet_count = int(os.getenv('ARG_PACKET_COUNT', '0'))
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[*] Network Traffic Monitor")
    print(f"[*] Interface: {interface}")
    print(f"[*] Duration: {duration}s")
    
    if bpf_filter:
        print(f"[*] Filter: {bpf_filter}")
    
    print(f"[*] Capturing traffic...\n")
    
    try:
        if packet_count > 0:
            sniff(iface=interface, prn=packet_callback, count=packet_count, 
                  filter=bpf_filter if bpf_filter else None, verbose=False)
        else:
            sniff(iface=interface, prn=packet_callback, timeout=duration, 
                  filter=bpf_filter if bpf_filter else None, verbose=False)
        
        print(f"\n[+] Capture Complete:")
        print(f"\n    Total Packets: {stats['total_packets']}")
        print(f"    Total Bytes: {stats['bytes']}")
        print(f"    TCP: {stats['tcp']}")
        print(f"    UDP: {stats['udp']}")
        print(f"    ICMP: {stats['icmp']}")
        print(f"    Other: {stats['other']}")
        
        print(f"\n    Unique Source IPs: {len(stats['src_ips'])}")
        print(f"    Unique Dest IPs: {len(stats['dst_ips'])}")
        
        if stats['total_packets'] > 0:
            avg_size = stats['bytes'] / stats['total_packets']
            print(f"    Avg Packet Size: {avg_size:.1f} bytes")
        
        print(f"\n[+] Complete")
    except PermissionError:
        print("[!] This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
