# arp-spoofer

## Description
Perform ARP spoofing to intercept traffic on a local network (MITM setup)

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
network, mitm, arp


## Links

## Options

### target_ip
- **Type:** string
- **Description:** Target IP address to spoof
- **Required:** Yes

### gateway_ip
- **Type:** string
- **Description:** Gateway IP address
- **Required:** Yes

### interface
- **Type:** string
- **Description:** Network interface to use
- **Required:** No
- **Default:** eth0

### duration
- **Type:** integer
- **Description:** Duration in seconds (0 for infinite)
- **Required:** No
- **Default:** 60
