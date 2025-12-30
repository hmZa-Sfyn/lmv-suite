#!/usr/bin/env python3
"""
DNS Resolver Module
Custom DNS query tool to resolve domains and display all record types
"""

import os
import sys
import socket

try:
    import dns.resolver
    import dns.rdatatype
except ImportError:
    print("Error: dnspython not installed. Install with: pip install dnspython")
    sys.exit(1)

def resolve_domain(domain, record_type="A", nameserver=None):
    """Resolve domain using specified record type"""
    try:
        resolver = dns.resolver.Resolver()
        if nameserver:
            resolver.nameservers = [nameserver]
        
        answers = resolver.resolve(domain, record_type, tcp=False)
        results = []
        for rdata in answers:
            results.append(str(rdata))
        return results
    except dns.resolver.NXDOMAIN:
        return None
    except dns.resolver.NoAnswer:
        return []
    except Exception as e:
        return None

def resolve_all(domain, nameserver=None):
    """Resolve all common record types"""
    record_types = ["A", "AAAA", "MX", "NS", "TXT", "CNAME", "SOA"]
    results = {}
    
    for rtype in record_types:
        try:
            result = resolve_domain(domain, rtype, nameserver)
            if result is not None:
                results[rtype] = result
        except:
            pass
    
    return results

def main():
    domain = os.getenv('ARG_DOMAIN')
    record_type = os.getenv('ARG_RECORD_TYPE', 'A')
    nameserver = os.getenv('ARG_NAMESERVER')
    
    if not domain:
        print("Error: domain required")
        sys.exit(1)
    
    print(f"[*] Resolving: {domain}")
    
    if record_type.upper() == "ALL":
        results = resolve_all(domain, nameserver)
        if not results:
            print("[!] No records found")
            sys.exit(1)
        
        for rtype, addrs in results.items():
            print(f"\n[+] {rtype} Records:")
            for addr in addrs:
                print(f"    {addr}")
    else:
        results = resolve_domain(domain, record_type.upper(), nameserver)
        if results is None:
            print(f"[!] Domain not found")
            sys.exit(1)
        
        if not results:
            print(f"[!] No {record_type} records found")
            sys.exit(0)
        
        print(f"\n[+] {record_type} Records:")
        for result in results:
            print(f"    {result}")
    
    print("\n[+] Complete")

if __name__ == "__main__":
    main()
