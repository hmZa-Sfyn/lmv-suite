#!/usr/bin/env python3
"""
WiFi Scanner Module
Passive Wi-Fi network scanner to detect nearby access points and clients
"""

import os
import sys
import time

try:
    from scapy.all import sniff, Dot11, Dot11Beacon, Dot11ProbeResp
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

networks = {}
clients = set()

def packet_callback(packet):
    """Process captured packets"""
    global networks, clients
    
    # Check for Beacon frames
    if packet.haslayer(Dot11Beacon):
        ssid = packet[Dot11Beacon].network_stats.get("SSID", b"").decode('utf-8', errors='ignore')
        bssid = packet[Dot11].addr3
        signal = packet.dBm_AntSignal if hasattr(packet, 'dBm_AntSignal') else -100
        
        if bssid not in networks:
            networks[bssid] = {
                'ssid': ssid if ssid else "<hidden>",
                'signal': signal,
                'count': 1
            }
        else:
            networks[bssid]['count'] += 1
            if signal > networks[bssid]['signal']:
                networks[bssid]['signal'] = signal
    
    # Check for Probe Response frames
    elif packet.haslayer(Dot11ProbeResp):
        ssid = packet[Dot11ProbeResp].network_stats.get("SSID", b"").decode('utf-8', errors='ignore')
        bssid = packet[Dot11].addr3
        
        if bssid not in networks:
            networks[bssid] = {
                'ssid': ssid if ssid else "<hidden>",
                'signal': -100,
                'count': 1
            }

def main():
    interface = os.getenv('ARG_INTERFACE')
    duration = int(os.getenv('ARG_DURATION', '10'))
    packet_count = int(os.getenv('ARG_PACKET_COUNT', '0'))
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[*] WiFi Scanner on {interface}")
    print(f"[*] Duration: {duration}s, Packets: {packet_count if packet_count > 0 else 'unlimited'}")
    print(f"[*] Scanning for networks...\n")
    
    try:
        # Set interface to monitor mode would go here (requires root)
        # For now, just do passive scanning on the interface
        
        if packet_count > 0:
            sniff(iface=interface, prn=packet_callback, count=packet_count, timeout=duration)
        else:
            sniff(iface=interface, prn=packet_callback, timeout=duration)
        
        print(f"\n\n[+] Found {len(networks)} networks:")
        print(f"\n{'BSSID':<20} {'SSID':<30} {'Signal':<10} {'Frames':<8}")
        print("-" * 68)
        
        for bssid, info in sorted(networks.items(), key=lambda x: x[1]['signal'], reverse=True):
            ssid = info['ssid'][:27] if len(info['ssid']) > 27 else info['ssid']
            signal = f"{info['signal']}dBm"
            print(f"{bssid:<20} {ssid:<30} {signal:<10} {info['count']:<8}")
        
        print(f"\n[+] Complete")
    except PermissionError:
        print("[!] This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
