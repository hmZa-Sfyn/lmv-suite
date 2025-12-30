#!/usr/bin/env python3
"""
VoIP Sniffer Module
Detect and capture VoIP (SIP/RTP) packets
"""

import os
import sys

try:
    from scapy.all import sniff, Raw
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

sip_calls = []
rtp_streams = set()
packet_count = 0

def packet_callback(packet):
    """Process VoIP packets"""
    global packet_count, sip_calls, rtp_streams
    
    packet_count += 1
    
    # Check for SIP packets
    if packet.haslayer(Raw):
        payload = bytes(packet[Raw].load)
        
        # Check for SIP (port 5060)
        if b'SIP/2.0' in payload:
            try:
                sip_info = payload.decode('utf-8', errors='ignore')
                lines = sip_info.split('\r\n')
                
                method = lines[0].split()[0] if lines else "UNKNOWN"
                
                sip_calls.append({
                    'method': method,
                    'src': packet[0].src if packet.haslayer('IP') else 'unknown',
                    'dst': packet[0].dst if packet.haslayer('IP') else 'unknown'
                })
                
                print(f"[+] SIP Packet: {method}")
            except:
                pass
        
        # Check for RTP (port range 16000-32000)
        if packet.haslayer('UDP'):
            sport = packet['UDP'].sport
            dport = packet['UDP'].dport
            
            if (16000 <= sport <= 32000) or (16000 <= dport <= 32000):
                stream = f"{packet[0].src}:{sport} -> {packet[0].dst}:{dport}"
                if stream not in rtp_streams:
                    rtp_streams.add(stream)
                    print(f"[+] RTP Stream Detected: {stream}")

def main():
    interface = os.getenv('ARG_INTERFACE')
    duration = int(os.getenv('ARG_DURATION', '60'))
    max_packets = int(os.getenv('ARG_PACKET_COUNT', '0'))
    
    if not interface:
        print("Error: interface required")
        sys.exit(1)
    
    print(f"[*] VoIP Sniffer on {interface}")
    print(f"[*] Duration: {duration}s")
    print(f"[*] Listening for SIP and RTP traffic...\n")
    
    try:
        if max_packets > 0:
            sniff(iface=interface, prn=packet_callback, count=max_packets, timeout=duration)
        else:
            sniff(iface=interface, prn=packet_callback, timeout=duration)
        
        print(f"\n\n[+] Capture Complete")
        print(f"[+] Total packets: {packet_count}")
        print(f"[+] SIP Calls detected: {len(sip_calls)}")
        print(f"[+] RTP Streams detected: {len(rtp_streams)}")
        
        if sip_calls:
            print(f"\n[+] SIP Call Details:")
            for i, call in enumerate(sip_calls, 1):
                print(f"    {i}. {call['method']}: {call['src']} -> {call['dst']}")
        
        if rtp_streams:
            print(f"\n[+] RTP Streams:")
            for stream in rtp_streams:
                print(f"    {stream}")
    
    except PermissionError:
        print("[!] This requires root/administrator privileges")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
