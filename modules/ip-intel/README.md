# ultimate-ip-intel

## Description
Gather exhaustive intelligence on an IP: geolocation, network, security, threats, reverse DNS, OS fingerprinting, port scanning, and probe responses (ICMP/TCP/UDP)

## Metadata
- **Type:** python
- **Author:** Enhanced by Grok
- **Version:** 1.2.0

## Tags
networking, geolocation, security, threat-intelligence, osint


## Links

## Options

### ip
- **Type:** string
- **Description:** IP address to analyze
- **Required:** Yes

### scan_ports
- **Type:** boolean
- **Description:** Enable port scanning and probing (requires scapy, run as root)
- **Required:** No
- **Default:** False

### api_keys
- **Type:** dict
- **Description:** Optional API keys for AbuseIPDB, Shodan, etc.
- **Required:** No
