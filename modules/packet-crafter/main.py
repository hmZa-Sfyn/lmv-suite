#!/usr/bin/env python3
"""
Packet Crafter Module
Advanced custom packet builder for any protocol
"""

import os
import sys

try:
    from scapy.all import IP, TCP, UDP, ICMP, Raw, send
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

def build_packet():
    """Build custom packet based on environment variables"""
    protocol = os.getenv('ARG_PROTOCOL', 'IP').upper()
    dst = os.getenv('ARG_DST', '127.0.0.1')
    src = os.getenv('ARG_SRC', '0.0.0.0')
    dport = os.getenv('ARG_DPORT')
    sport = os.getenv('ARG_SPORT')
    payload = os.getenv('ARG_PAYLOAD', '')
    
    # Build IP layer
    packet = IP(dst=dst, src=src)
    
    # Add protocol layer
    if protocol == "TCP":
        tcp_kwargs = {}
        if dport:
            tcp_kwargs['dport'] = int(dport)
        if sport:
            tcp_kwargs['sport'] = int(sport)
        packet = packet / TCP(**tcp_kwargs)
    
    elif protocol == "UDP":
        udp_kwargs = {}
        if dport:
            udp_kwargs['dport'] = int(dport)
        if sport:
            udp_kwargs['sport'] = int(sport)
        packet = packet / UDP(**udp_kwargs)
    
    elif protocol == "ICMP":
        packet = packet / ICMP()
    
    elif protocol == "DNS":
        if dport is None:
            dport = 53
        packet = packet / UDP(sport=random.randint(49152, 65535), dport=int(dport))
    
    # Add payload
    if payload:
        packet = packet / Raw(load=payload.encode())
    
    return packet

def display_packet(packet):
    """Display packet information"""
    print("\n[+] Packet Details:")
    print("-" * 50)
    
    if packet.haslayer('IP'):
        ip = packet['IP']
        print(f"[*] IP Layer:")
        print(f"    Source: {ip.src}")
        print(f"    Destination: {ip.dst}")
        print(f"    TTL: {ip.ttl}")
        print(f"    Length: {ip.len}")
    
    if packet.haslayer('TCP'):
        tcp = packet['TCP']
        print(f"\n[*] TCP Layer:")
        print(f"    Source Port: {tcp.sport}")
        print(f"    Dest Port: {tcp.dport}")
        print(f"    Flags: {tcp.flags}")
        print(f"    Seq: {tcp.seq}")
        print(f"    Ack: {tcp.ack}")
    
    elif packet.haslayer('UDP'):
        udp = packet['UDP']
        print(f"\n[*] UDP Layer:")
        print(f"    Source Port: {udp.sport}")
        print(f"    Dest Port: {udp.dport}")
        print(f"    Length: {udp.len}")
    
    elif packet.haslayer('ICMP'):
        icmp = packet['ICMP']
        print(f"\n[*] ICMP Layer:")
        print(f"    Type: {icmp.type}")
        print(f"    Code: {icmp.code}")
    
    if packet.haslayer('Raw'):
        raw = packet['Raw']
        payload = str(raw.load)[:100]
        print(f"\n[*] Payload:")
        print(f"    {payload}")
    
    print(f"\n[+] Packet Size: {len(packet)} bytes")
    print(f"[+] Packet Hex:\n{packet.hexdump()}")

def main():
    protocol = os.getenv('ARG_PROTOCOL')
    send_packet = os.getenv('ARG_SEND', 'false').lower() == 'true'
    
    if not protocol:
        print("Error: protocol required")
        sys.exit(1)
    
    print(f"[*] Packet Crafter - Building {protocol} packet")
    
    try:
        packet = build_packet()
        display_packet(packet)
        
        if send_packet:
            print(f"\n[*] Sending packet...")
            send(packet, verbose=False)
            print("[+] Packet sent!")
        else:
            print(f"\n[*] Packet built but not sent (use ARG_SEND=true to send)")
    
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
