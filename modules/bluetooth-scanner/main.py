#!/usr/bin/env python3
"""
Bluetooth Scanner Module
Scan for nearby Bluetooth devices
"""

import os
import sys
import subprocess

try:
    import pybluez
except ImportError:
    print("Note: pybluez not installed, using hcitools")
    pybluez = None

def scan_bluetooth_hci(duration=10):
    """Scan using hcitools"""
    try:
        result = subprocess.run(['hcitool', 'scan', '--flush'], 
                              capture_output=True, text=True, timeout=duration)
        
        devices = []
        for line in result.stdout.split('\n')[1:]:
            if line.strip():
                parts = line.split()
                if len(parts) >= 2:
                    mac = parts[0]
                    name = ' '.join(parts[1:])
                    devices.append({
                        'mac': mac,
                        'name': name,
                        'rssi': 'N/A'
                    })
        
        return devices
    except Exception as e:
        print(f"Error: {e}")
        return []

def scan_bluetooth_pybluez(duration=10):
    """Scan using pybluez"""
    try:
        from bluetooth import discover_devices, lookup_name
        
        devices = []
        nearby_devices = discover_devices(duration=duration, lookup_names=True)
        
        for addr, name in nearby_devices:
            devices.append({
                'mac': addr,
                'name': name,
                'rssi': 'N/A'
            })
        
        return devices
    except Exception as e:
        print(f"Error: {e}")
        return []

def get_device_type(name):
    """Guess device type from name"""
    if not name:
        return "Unknown"
    
    name_lower = name.lower()
    
    if 'iphone' in name_lower or 'ipad' in name_lower:
        return "Apple Device"
    elif 'samsung' in name_lower or 'pixel' in name_lower:
        return "Android Device"
    elif 'airpod' in name_lower or 'beats' in name_lower or 'bose' in name_lower:
        return "Headset"
    elif 'watch' in name_lower:
        return "Smartwatch"
    elif 'car' in name_lower:
        return "Car"
    else:
        return "Device"

def main():
    duration = int(os.getenv('ARG_DURATION', '10'))
    device_type = os.getenv('ARG_DEVICE_TYPE', 'all').lower()
    show_rssi = os.getenv('ARG_SHOW_RSSI', 'true').lower() == 'true'
    
    print(f"[*] Bluetooth Scanner")
    print(f"[*] Duration: {duration}s")
    print(f"[*] Scanning for Bluetooth devices...\n")
    
    # Try pybluez first, fall back to hcitools
    if pybluez:
        devices = scan_bluetooth_pybluez(duration)
    else:
        devices = scan_bluetooth_hci(duration)
    
    if not devices:
        print("[!] No Bluetooth devices found")
        sys.exit(0)
    
    print(f"[+] Found {len(devices)} device(s):")
    print(f"\n{'#':<3} {'MAC Address':<20} {'Device Name':<35} {'Type':<15}")
    print("-" * 75)
    
    for i, device in enumerate(devices, 1):
        dev_type = get_device_type(device['name'])
        name = device['name'][:32] if len(device['name']) > 32 else device['name']
        
        print(f"{i:<3} {device['mac']:<20} {name:<35} {dev_type:<15}")
    
    print(f"\n[+] Device Type Breakdown:")
    
    type_count = {}
    for device in devices:
        dtype = get_device_type(device['name'])
        type_count[dtype] = type_count.get(dtype, 0) + 1
    
    for dtype, count in sorted(type_count.items(), key=lambda x: x[1], reverse=True):
        print(f"    {dtype}: {count}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
