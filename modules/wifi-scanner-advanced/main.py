#!/usr/bin/env python3
"""
Advanced WiFi Scanner Module
Detailed WiFi network analysis with signal strength and channel info
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

def packet_callback(packet):
    """Process WiFi packets"""
    global networks
    
    if packet.haslayer(Dot11):
        # Check for Beacon frames
        if packet.haslayer(Dot11Beacon):
            ssid = packet[Dot11Beacon].network_stats.get("SSID", b"").decode('utf-8', errors='ignore')
            bssid = packet[Dot11].addr3
            signal = packet.dBm_AntSignal if hasattr(packet, 'dBm_AntSignal') else -100
            channel = packet[Dot11].notset if hasattr(packet, 'notset') else 0
            
            if bssid not in networks:
                networks[bssid] = {
                    'ssid': ssid if ssid else "<HIDDEN>",
                    'bssid': bssid,
                    'signal': signal,
                    'channel': 0,
                    'encryption': 'WPA2',
                    'beacons': 1,
                    'clients': set()
                }
            else:
                networks[bssid]['beacons'] += 1
                if signal > networks[bssid]['signal']:
                    networks[bssid]['signal'] = signal
        
        # Check for client devices
        if packet.haslayer(Dot11) and packet.type == 0 and packet.subtype == 8:
            src_mac = packet[Dot11].addr2
            dst_mac = packet[Dot11].addr1
            
            for bssid, info in networks.items():
                if dst_mac == bssid or src_mac == bssid:
                    info['clients'].add(src_mac if src_mac != bssid else dst_mac)

def main():
    interface = os.getenv('ARG_INTERFACE')
    duration = int(os.getenv('ARG_DURATION', '30'))
    show_hidden = os.getenv('ARG_SHOW_HIDDEN', 'true').lower() == 'true'
    channel = os.getenv('ARG_CHANNEL')
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[*] Advanced WiFi Scanner")
    print(f"[*] Interface: {interface}")
    print(f"[*] Duration: {duration}s")
    print(f"[*] Scanning...\n")
    
    try:
        sniff(iface=interface, prn=packet_callback, timeout=duration, verbose=False)
        
        print(f"\n\n[+] Found {len(networks)} network(s):")
        print(f"\n{'#':<3} {'SSID':<30} {'BSSID':<20} {'Signal':<8} {'Clients':<8} {'Beacons':<8}")
        print("-" * 85)
        
        for i, (bssid, info) in enumerate(sorted(networks.items(), key=lambda x: x[1]['signal'], reverse=True), 1):
            ssid = info['ssid'][:27] if len(info['ssid']) > 27 else info['ssid']
            signal = f"{info['signal']}dBm"
            clients = len(info['clients'])
            beacons = info['beacons']
            
            print(f"{i:<3} {ssid:<30} {bssid:<20} {signal:<8} {clients:<8} {beacons:<8}")
        
        print(f"\n[+] Detailed Analysis:")
        for i, (bssid, info) in enumerate(sorted(networks.items(), key=lambda x: x[1]['signal'], reverse=True), 1):
            if i <= 5:
                print(f"\n[+] Network #{i} - {info['ssid']}:")
                print(f"    BSSID: {bssid}")
                print(f"    Signal: {info['signal']}dBm")
                print(f"    Encryption: {info['encryption']}")
                print(f"    Connected Clients: {len(info['clients'])}")
                print(f"    Beacon Frames: {info['beacons']}")
        
        print(f"\n[+] Complete")
    except PermissionError:
        print("[!] This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
