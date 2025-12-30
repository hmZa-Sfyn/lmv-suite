#!/usr/bin/env python3
"""
SSL Analyzer Module
Inspect TLS/SSL handshakes and certificates from target hosts
"""

import os
import sys
import socket
import ssl
from datetime import datetime

def get_certificate_info(host, port, timeout):
    """Get certificate information from host"""
    try:
        context = ssl.create_default_context()
        context.check_hostname = False
        context.verify_mode = ssl.CERT_NONE
        
        with socket.create_connection((host, port), timeout=timeout) as sock:
            with context.wrap_socket(sock, server_hostname=host) as ssock:
                cert = ssock.getpeercert()
                cert_der = ssock.getpeercert(binary_form=True)
                
                return {
                    'subject': dict(x[0] for x in cert.get('subject', [])),
                    'issuer': dict(x[0] for x in cert.get('issuer', [])),
                    'version': cert.get('version'),
                    'serial': cert.get('serialNumber'),
                    'notBefore': cert.get('notBefore'),
                    'notAfter': cert.get('notAfter'),
                    'subjectAltName': cert.get('subjectAltName', []),
                    'cipher': ssock.cipher()[0],
                    'protocol': ssock.version()
                }
    except Exception as e:
        return None

def parse_cert_date(date_str):
    """Parse certificate date"""
    try:
        return datetime.strptime(date_str, "%b %d %H:%M:%S %Y %Z")
    except:
        return date_str

def main():
    host = os.getenv('ARG_HOST')
    port = int(os.getenv('ARG_PORT', '443'))
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    
    if not host:
        print("Error: host required")
        sys.exit(1)
    
    print(f"[*] Analyzing SSL/TLS for {host}:{port}")
    
    cert_info = get_certificate_info(host, port, timeout)
    
    if not cert_info:
        print("[!] Failed to retrieve certificate")
        sys.exit(1)
    
    print("\n[+] Certificate Information:")
    
    subject = cert_info.get('subject', {})
    print(f"\n    Subject:")
    for key, value in subject.items():
        print(f"      {key}: {value}")
    
    issuer = cert_info.get('issuer', {})
    print(f"\n    Issuer:")
    for key, value in issuer.items():
        print(f"      {key}: {value}")
    
    print(f"\n    Certificate Details:")
    print(f"      Serial: {cert_info.get('serial')}")
    print(f"      Valid From: {cert_info.get('notBefore')}")
    print(f"      Valid Until: {cert_info.get('notAfter')}")
    print(f"      Version: {cert_info.get('version')}")
    
    print(f"\n    TLS/SSL Information:")
    print(f"      Protocol: {cert_info.get('protocol')}")
    print(f"      Cipher: {cert_info.get('cipher')}")
    
    san = cert_info.get('subjectAltName', [])
    if san:
        print(f"\n    Subject Alternative Names:")
        for name_type, name_value in san:
            print(f"      {name_type}: {name_value}")
    
    # Check expiration
    try:
        expiry_date = parse_cert_date(cert_info.get('notAfter', ''))
        if isinstance(expiry_date, datetime):
            days_left = (expiry_date - datetime.utcnow()).days
            if days_left < 0:
                print(f"\n[!] Certificate EXPIRED {abs(days_left)} days ago!")
            elif days_left < 30:
                print(f"\n[!] Certificate expires in {days_left} days")
            else:
                print(f"\n[+] Certificate valid for {days_left} more days")
    except:
        pass
    
    print("\n[+] Complete")

if __name__ == "__main__":
    main()
