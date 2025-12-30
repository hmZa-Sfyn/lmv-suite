#!/usr/bin/env python3
"""
Rogue AP Detector Module
Detect rogue WiFi access points and unauthorized devices
"""

import os
import sys

try:
    from scapy.all import sniff, Dot11, Dot11Beacon
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

networks = {}
suspicious = []

def packet_callback(packet):
    """Analyze WiFi packets for rogue APs"""
    global networks, suspicious
    
    if packet.haslayer(Dot11Beacon):
        bssid = packet[Dot11].addr3
        ssid = packet[Dot11Beacon].network_stats.get("SSID", b"").decode('utf-8', errors='ignore')
        signal = packet.dBm_AntSignal if hasattr(packet, 'dBm_AntSignal') else -100
        
        if bssid not in networks:
            networks[bssid] = {
                'ssid': ssid,
                'signal': signal,
                'count': 1
            }
        else:
            networks[bssid]['count'] += 1

def detect_rogue_patterns(networks, known_networks=None):
    """Detect suspicious patterns"""
    suspicious_aps = []
    
    for bssid, info in networks.items():
        ssid = info['ssid']
        
        # Check for common rogue AP names
        if any(x in ssid.lower() for x in ['free wifi', 'open', 'guest', 'public']):
            if info['signal'] > -70:
                suspicious_aps.append({
                    'bssid': bssid,
                    'ssid': ssid,
                    'reason': 'Suspicious SSID with strong signal',
                    'risk': 'HIGH'
                })
        
        # Check for hidden networks
        if not ssid:
            suspicious_aps.append({
                'bssid': bssid,
                'ssid': '<HIDDEN>',
                'reason': 'Hidden network',
                'risk': 'MEDIUM'
            })
        
        # Check for signal cloning (same SSID with multiple BSSIDs)
        for other_bssid, other_info in networks.items():
            if other_bssid != bssid and other_info['ssid'] == ssid:
                suspicious_aps.append({
                    'bssid': bssid,
                    'ssid': ssid,
                    'reason': f'Duplicate SSID detected ({other_bssid})',
                    'risk': 'CRITICAL'
                })
    
    return suspicious_aps

def main():
    interface = os.getenv('ARG_INTERFACE')
    known_networks = os.getenv('ARG_KNOWN_NETWORKS')
    duration = int(os.getenv('ARG_DURATION', '30'))
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[*] Rogue AP Detector")
    print(f"[*] Interface: {interface}")
    print(f"[*] Duration: {duration}s")
    
    if known_networks and os.path.exists(known_networks):
        with open(known_networks, 'r') as f:
            known = [line.strip() for line in f if line.strip()]
        print(f"[*] Known networks: {len(known)}")
    else:
        known = None
    
    print(f"[*] Scanning for rogue APs...\n")
    
    try:
        sniff(iface=interface, prn=packet_callback, timeout=duration, verbose=False)
        
        if not networks:
            print("[!] No networks found")
            sys.exit(0)
        
        # Detect rogue patterns
        suspicious_aps = detect_rogue_patterns(networks, known)
        
        print(f"[+] Found {len(networks)} network(s)")
        
        if suspicious_aps:
            print(f"\n[!] SUSPICIOUS NETWORKS DETECTED ({len(suspicious_aps)}):")
            print(f"\n{'#':<3} {'BSSID':<20} {'SSID':<30} {'Risk':<10} {'Reason':<35}")
            print("-" * 100)
            
            for i, ap in enumerate(suspicious_aps, 1):
                ssid = ap['ssid'][:27] if len(ap['ssid']) > 27 else ap['ssid']
                reason = ap['reason'][:32] if len(ap['reason']) > 32 else ap['reason']
                print(f"{i:<3} {ap['bssid']:<20} {ssid:<30} {ap['risk']:<10} {reason:<35}")
        else:
            print(f"\n[+] No obvious rogue APs detected")
        
        print(f"\n[+] All Networks:")
        print(f"\n{'#':<3} {'BSSID':<20} {'SSID':<30} {'Beacons':<10}")
        print("-" * 65)
        
        for i, (bssid, info) in enumerate(sorted(networks.items()), 1):
            ssid = info['ssid'][:27] if len(info['ssid']) > 27 else info['ssid']
            print(f"{i:<3} {bssid:<20} {ssid:<30} {info['count']:<10}")
        
        print(f"\n[+] Complete")
    except PermissionError:
        print("[!] This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
