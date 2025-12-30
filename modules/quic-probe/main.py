#!/usr/bin/env python3
"""
QUIC Probe Module
Probe and analyze QUIC/HTTP3 traffic
"""

import os
import sys
import socket
import struct
import time

def send_quic_initial(host, port, timeout):
    """Send QUIC Initial Packet"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.settimeout(timeout)
        
        # Minimal QUIC Initial packet header
        # QUIC packet format:
        # - Packet Type + Flags (1 byte)
        # - Version (4 bytes)
        # - DCID Len (1 byte)
        # - DCID (varies)
        # - SCID Len (1 byte)
        # - SCID (varies)
        # - Token Length (var)
        # - Length (var)
        # - Packet Number (1-4 bytes)
        # - Payload
        
        import secrets
        dcid = secrets.token_bytes(8)
        scid = secrets.token_bytes(8)
        
        # Build minimal QUIC packet
        packet = bytearray()
        packet.append(0xc0)  # Initial packet, long header
        packet.extend(struct.pack('>I', 0x00000001))  # Version 1
        packet.append(0x08)  # DCID Length
        packet.extend(dcid)
        packet.append(0x08)  # SCID Length
        packet.extend(scid)
        
        sock.sendto(bytes(packet), (host, port))
        
        try:
            response, addr = sock.recvfrom(1024)
            return True, response
        except socket.timeout:
            return False, None
    except:
        return False, None
    finally:
        sock.close()

def check_quic_support(host, port, timeout):
    """Check if host supports QUIC"""
    success, response = send_quic_initial(host, port, timeout)
    
    if success and response:
        # Check for Version Negotiation (0xc0 version negotiation)
        if len(response) > 5:
            packet_type = response[0] & 0xf0
            if packet_type == 0xc0:
                # Likely QUIC version negotiation
                return True, response
        return True, response
    
    return False, None

def main():
    host = os.getenv('ARG_HOST')
    port = int(os.getenv('ARG_PORT', '443'))
    timeout = int(os.getenv('ARG_TIMEOUT', '5'))
    probe_type = os.getenv('ARG_PROBE_TYPE', 'handshake').lower()
    
    if not host:
        print("Error: host required")
        sys.exit(1)
    
    print(f"[*] QUIC Probe Module")
    print(f"[*] Target: {host}:{port}")
    print(f"[*] Probe Type: {probe_type}")
    print(f"[*] Timeout: {timeout}s")
    
    # Try to resolve hostname
    try:
        ip = socket.gethostbyname(host)
        print(f"[*] Resolved to: {ip}")
    except socket.gaierror:
        print(f"[!] Could not resolve host: {host}")
        sys.exit(1)
    
    print(f"\n[*] Probing QUIC support...")
    
    # Check QUIC support
    success, response = check_quic_support(ip, port, timeout)
    
    if success:
        print(f"[+] QUIC endpoint detected!")
        
        if response:
            print(f"\n[+] Response Details:")
            print(f"    Length: {len(response)} bytes")
            print(f"    First bytes (hex): {response[:20].hex()}")
            
            # Try to parse response
            if len(response) > 0:
                packet_type = response[0]
                print(f"    Packet Type: 0x{packet_type:02x}")
                
                if len(response) >= 5:
                    version = struct.unpack('>I', response[1:5])[0]
                    print(f"    QUIC Version: 0x{version:08x}")
                    
                    if version == 0x00000001:
                        print(f"    -> Version 1 (RFC 9000)")
                    elif version == 0xff000000:
                        print(f"    -> Version Negotiation")
    else:
        print(f"[!] No QUIC response detected")
        print(f"[*] Host may not support QUIC on port {port}")
        print(f"[*] Common QUIC ports: 80, 443, 8080")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
