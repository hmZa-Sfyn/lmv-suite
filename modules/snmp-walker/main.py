#!/usr/bin/env python3
"""
SNMP Walker Module
Perform SNMP walks to extract system information
"""

import os
import sys

try:
    from pysnmp.hlapi import getCmd, SnmpEngine, CommunityData, UdpTransportTarget, ContextData, ObjectType, ObjectIdentity
    from pysnmp.smi import builder, view
except ImportError:
    print("Error: pysnmp not installed. Install with: pip install pysnmp")
    sys.exit(1)

def snmp_get(host, community, oid, version="2c"):
    """Perform SNMP GET"""
    try:
        snmp_engine = SnmpEngine()
        
        iterator = getCmd(
            snmp_engine,
            CommunityData(community, mpModel=1 if version == "2c" else 0),
            UdpTransportTarget((host, 161)),
            ContextData(),
            ObjectType(ObjectIdentity(oid))
        )
        
        error_indication, error_status, error_index, var_binds = next(iterator)
        
        if error_indication:
            return None
        
        if error_status:
            return None
        
        for var_bind in var_binds:
            return str(var_bind[1])
    
    except Exception as e:
        return None

def get_system_info(host, community, version="2c"):
    """Get system information via SNMP"""
    info = {}
    
    # Common OIDs
    oids = {
        'sysDescr': '1.3.6.1.2.1.1.1.0',
        'sysObjectID': '1.3.6.1.2.1.1.2.0',
        'sysUpTime': '1.3.6.1.2.1.1.3.0',
        'sysContact': '1.3.6.1.2.1.1.4.0',
        'sysName': '1.3.6.1.2.1.1.5.0',
        'sysLocation': '1.3.6.1.2.1.1.6.0'
    }
    
    for name, oid in oids.items():
        value = snmp_get(host, community, oid, version)
        if value:
            info[name] = value
    
    return info

def main():
    host = os.getenv('ARG_HOST')
    community = os.getenv('ARG_COMMUNITY', 'public')
    version = os.getenv('ARG_VERSION', '2c')
    oid = os.getenv('ARG_OID', '1.3.6.1.2.1')
    
    if not host:
        print("Error: host required")
        sys.exit(1)
    
    print(f"[*] SNMP Walker")
    print(f"[*] Host: {host}")
    print(f"[*] Community: {community}")
    print(f"[*] Version: {version}")
    
    # Try to get system info
    print(f"\n[*] Retrieving system information...")
    
    info = get_system_info(host, community, version)
    
    if not info:
        print(f"[!] Failed to retrieve SNMP data")
        print(f"[*] SNMP may not be enabled or community string is incorrect")
        sys.exit(1)
    
    print(f"\n[+] System Information:")
    print(f"{'Key':<20} {'Value':<50}")
    print("-" * 70)
    
    for key, value in info.items():
        val_short = str(value)[:47]
        print(f"{key:<20} {val_short:<50}")
    
    # Additional info retrieval
    print(f"\n[*] Getting interface information...")
    
    # sysUpTime
    uptime_val = snmp_get(host, community, '1.3.6.1.2.1.1.3.0', version)
    if uptime_val:
        print(f"[+] System Uptime: {uptime_val}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
