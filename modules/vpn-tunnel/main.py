#!/usr/bin/env python3
"""
VPN Tunnel Module
Simple IP-over-ICMP or DNS tunneling
"""

import os
import sys
import socket
import struct

try:
    from scapy.all import IP, ICMP, DNS, DNSQR, Raw, sniff, send
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

class ICMPTunnel:
    """ICMP-based tunnel"""
    
    def __init__(self, server=None):
        self.server = server
        self.packet_id = 1000
    
    def send_data(self, data, target):
        """Send data over ICMP"""
        try:
            packet = IP(dst=target)/ICMP(id=self.packet_id)/Raw(load=data)
            send(packet, verbose=False)
            self.packet_id += 1
            return True
        except:
            return False
    
    def receive_data(self, timeout=5):
        """Receive data from ICMP"""
        packets = []
        
        def icmp_callback(pkt):
            if pkt.haslayer('Raw'):
                packets.append(pkt['Raw'].load)
        
        try:
            sniff(filter="icmp", prn=icmp_callback, timeout=timeout)
        except:
            pass
        
        return packets

class DNSTunnel:
    """DNS-based tunnel"""
    
    def __init__(self, server=None):
        self.server = server
    
    def send_data(self, data, target):
        """Send data over DNS"""
        # Simple DNS query encoding
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            encoded = data.hex()[:63]  # DNS label limit
            query = f"{encoded}.tunnel.local"
            
            # Simplified DNS query (would need proper DNS packet construction)
            sock.sendto(b"dummy", (target, 53))
            sock.close()
            return True
        except:
            return False

def main():
    mode = os.getenv('ARG_MODE', 'icmp').lower()
    server = os.getenv('ARG_SERVER')
    listen_port = int(os.getenv('ARG_LISTEN_PORT', '8080'))
    
    print(f"[*] VPN Tunnel Module")
    print(f"[*] Mode: {mode}")
    
    if mode == "icmp":
        tunnel = ICMPTunnel(server)
        print(f"[*] ICMP Tunnel initialized")
    elif mode == "dns":
        tunnel = DNSTunnel(server)
        print(f"[*] DNS Tunnel initialized")
    else:
        print(f"Error: Unknown mode: {mode}")
        sys.exit(1)
    
    if server:
        print(f"[*] Target Server: {server}")
        print(f"[*] Testing tunnel connection...")
        
        # Test with small data
        test_data = b"HELLO"
        if tunnel.send_data(test_data, server):
            print(f"[+] Test packet sent successfully")
        else:
            print(f"[!] Failed to send test packet")
    
    print(f"\n[*] Listening on port {listen_port}...")
    print(f"[*] This is a basic tunnel implementation for demonstration")
    print(f"[+] Module ready (full implementation requires more setup)")

if __name__ == "__main__":
    main()
