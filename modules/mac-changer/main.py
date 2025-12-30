#!/usr/bin/env python3
"""
MAC Changer Module
Change the MAC address of a network interface
"""

import os
import sys
import subprocess
import random
import re

def validate_mac(mac):
    """Validate MAC address format"""
    pattern = r'^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$'
    return re.match(pattern, mac) is not None

def generate_random_mac():
    """Generate random MAC address"""
    mac = [random.randint(0x00, 0xff) for _ in range(6)]
    return ':'.join(map(lambda x: "%02x" % x, mac))

def get_current_mac(interface):
    """Get current MAC address"""
    try:
        result = subprocess.run(['ip', 'link', 'show', interface], 
                              capture_output=True, text=True)
        match = re.search(r'link/ether ([0-9a-f:]+)', result.stdout)
        if match:
            return match.group(1)
    except:
        pass
    return None

def bring_down_interface(interface):
    """Bring down network interface"""
    try:
        subprocess.run(['ip', 'link', 'set', interface, 'down'], 
                      check=True, capture_output=True)
        return True
    except:
        return False

def bring_up_interface(interface):
    """Bring up network interface"""
    try:
        subprocess.run(['ip', 'link', 'set', interface, 'up'], 
                      check=True, capture_output=True)
        return True
    except:
        return False

def change_mac(interface, new_mac):
    """Change MAC address"""
    try:
        subprocess.run(['ip', 'link', 'set', interface, 'address', new_mac], 
                      check=True, capture_output=True)
        return True
    except:
        return False

def main():
    interface = os.getenv('ARG_INTERFACE')
    new_mac = os.getenv('ARG_MAC')
    action = os.getenv('ARG_ACTION', 'change').lower()
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[*] MAC Changer for {interface}")
    
    current_mac = get_current_mac(interface)
    if not current_mac:
        print("[!] Could not get current MAC address")
        sys.exit(1)
    
    print(f"[*] Current MAC: {current_mac}")
    
    if action == "reset":
        print("[*] Action: Reset to original")
        print("[!] Note: Cannot retrieve original MAC (requires system info)")
        sys.exit(1)
    
    if action == "random":
        new_mac = generate_random_mac()
        print(f"[*] Generated random MAC: {new_mac}")
    elif action == "change":
        if not new_mac:
            print("Error: MAC address required for change action")
            sys.exit(1)
        if not validate_mac(new_mac):
            print(f"Error: Invalid MAC format: {new_mac}")
            sys.exit(1)
    else:
        print(f"Error: Unknown action: {action}")
        sys.exit(1)
    
    # Perform MAC change
    print(f"\n[*] Bringing down interface {interface}...")
    if not bring_down_interface(interface):
        print("[!] Failed to bring down interface")
        sys.exit(1)
    
    print(f"[*] Changing MAC to {new_mac}...")
    if not change_mac(interface, new_mac):
        print("[!] Failed to change MAC")
        bring_up_interface(interface)
        sys.exit(1)
    
    print(f"[*] Bringing up interface {interface}...")
    if not bring_up_interface(interface):
        print("[!] Failed to bring up interface")
        sys.exit(1)
    
    # Verify
    updated_mac = get_current_mac(interface)
    if updated_mac == new_mac:
        print(f"\n[+] Success! MAC changed from {current_mac} to {updated_mac}")
    else:
        print(f"[!] MAC change may have failed. Current: {updated_mac}")
    
    print("[+] Complete")

if __name__ == "__main__":
    main()
