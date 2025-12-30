#!/usr/bin/env python3
"""
WiFi Network Analyzer Module
Analyze detailed WiFi network information
"""

import os
import sys
import time

try:
    from scapy.all import sniff, Dot11, Dot11Beacon, Dot11Elt
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

networks = {}

def packet_callback(packet):
    """Process WiFi packets"""
    global networks
    
    if packet.haslayer(Dot11Beacon):
        bssid = packet[Dot11].addr3
        ssid = packet[Dot11Beacon].network_stats.get("SSID", b"").decode('utf-8', errors='ignore')
        signal = packet.dBm_AntSignal if hasattr(packet, 'dBm_AntSignal') else -100
        
        if bssid not in networks:
            networks[bssid] = {
                'ssid': ssid,
                'bssid': bssid,
                'signal': signal,
                'encryption': 'Unknown',
                'cipher': 'Unknown',
                'auth': 'Unknown',
                'channel': 0,
                'capabilities': []
            }
            
            # Parse capabilities
            if packet.haslayer(Dot11Beacon):
                cap = packet[Dot11Beacon].cap
                capabilities = []
                
                if cap & 0x0001:
                    capabilities.append('ESS')
                if cap & 0x0002:
                    capabilities.append('IBSS')
                if cap & 0x0004:
                    capabilities.append('CF_POLLABLE')
                if cap & 0x0008:
                    capabilities.append('CF_POLL_REQUEST')
                if cap & 0x0010:
                    capabilities.append('PRIVACY')
                if cap & 0x0020:
                    capabilities.append('SHORT_PREAMBLE')
                if cap & 0x0100:
                    capabilities.append('SHORT_SLOT_TIME')
                if cap & 0x0200:
                    capabilities.append('RADIO_MEASUREMENT')
                
                networks[bssid]['capabilities'] = capabilities
            
            # Parse RSN/WPA info
            if packet.haslayer(Dot11Elt):
                elem = packet[Dot11Elt]
                while elem:
                    if elem.ID == 48:  # RSN
                        networks[bssid]['encryption'] = 'WPA2'
                    elif elem.ID == 221:  # Vendor Specific (WPA)
                        networks[bssid]['encryption'] = 'WPA'
                    elem = elem.payload if isinstance(elem.payload, Dot11Elt) else None

def main():
    interface = os.getenv('ARG_INTERFACE')
    target_bssid = os.getenv('ARG_TARGET_BSSID')
    duration = int(os.getenv('ARG_DURATION', '15'))
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[*] WiFi Network Analyzer")
    print(f"[*] Interface: {interface}")
    print(f"[*] Duration: {duration}s")
    print(f"[*] Analyzing networks...\n")
    
    try:
        sniff(iface=interface, prn=packet_callback, timeout=duration, verbose=False)
        
        if not networks:
            print("[!] No networks found")
            sys.exit(0)
        
        # Filter if target specified
        if target_bssid:
            networks = {k: v for k, v in networks.items() if k == target_bssid}
        
        print(f"\n[+] Found {len(networks)} network(s):")
        
        for i, (bssid, info) in enumerate(sorted(networks.items(), key=lambda x: x[1]['signal'], reverse=True), 1):
            print(f"\n[+] Network #{i}:")
            print(f"    SSID: {info['ssid'] if info['ssid'] else '<HIDDEN>'}")
            print(f"    BSSID: {bssid}")
            print(f"    Signal: {info['signal']}dBm")
            print(f"    Encryption: {info['encryption']}")
            print(f"    Capabilities: {', '.join(info['capabilities'])}")
            
            # Security assessment
            if 'PRIVACY' in info['capabilities']:
                print(f"    [!] Network is encrypted")
            else:
                print(f"    [!] WARNING: Open network - no encryption!")
        
        print(f"\n[+] Complete")
    except PermissionError:
        print("[!] This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
