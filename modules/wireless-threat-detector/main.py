#!/usr/bin/env python3
"""
Wireless Threat Detector Module
Detect wireless security threats and anomalies
"""

import os
import sys
import time

try:
    from scapy.all import sniff, Dot11, Dot11ProbeReq, Dot11Deauth
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

threats = {
    'deauth': 0,
    'probe_requests': 0,
    'beacon_floods': 0,
    'spoofed_aps': 0
}

def packet_callback(packet):
    """Detect threats in packets"""
    global threats
    
    if packet.haslayer(Dot11):
        # Detect deauth packets
        if packet.haslayer(Dot11Deauth):
            threats['deauth'] += 1
            print(f"[!] Deauth detected: {packet[Dot11].addr2} -> {packet[Dot11].addr1}")
        
        # Detect probe requests
        elif packet.haslayer(Dot11ProbeReq):
            threats['probe_requests'] += 1

def main():
    interface = os.getenv('ARG_INTERFACE')
    duration = int(os.getenv('ARG_DURATION', '60'))
    alert_threshold = int(os.getenv('ARG_ALERT_THRESHOLD', '10'))
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[!] WARNING: This monitors for wireless threats")
    print(f"[!] Only use on networks you own or have permission to monitor")
    print(f"\n[*] Wireless Threat Detector")
    print(f"[*] Interface: {interface}")
    print(f"[*] Duration: {duration}s")
    print(f"[*] Alert Threshold: {alert_threshold}")
    print(f"[*] Monitoring for threats...\n")
    
    try:
        start_time = time.time()
        
        sniff(iface=interface, prn=packet_callback, timeout=duration, verbose=False)
        
        elapsed = time.time() - start_time
        
        print(f"\n\n[+] Threat Summary:")
        print(f"    Deauth Packets: {threats['deauth']}")
        print(f"    Probe Requests: {threats['probe_requests']}")
        print(f"    Beacon Floods: {threats['beacon_floods']}")
        print(f"    Spoofed APs: {threats['spoofed_aps']}")
        
        total_threats = sum(threats.values())
        
        if total_threats > alert_threshold:
            print(f"\n[!] ALERT: High threat activity detected ({total_threats} events)")
        else:
            print(f"\n[+] Network appears secure (threat count: {total_threats})")
        
        print(f"\n[+] Monitoring duration: {elapsed:.1f}s")
        print(f"[+] Complete")
    except PermissionError:
        print("[!] This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
