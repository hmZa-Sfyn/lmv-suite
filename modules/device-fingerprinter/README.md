# device-fingerprinter

## Description
Fingerprint connected devices by analyzing network behavior, ports, and services

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat
- **Version:** 1.0.0

## Tags
fingerprinting, reconnaissance, network


## Links

## Options

### target
- **Type:** string
- **Description:** Target IP address to fingerprint
- **Required:** Yes

### port_range
- **Type:** string
- **Description:** Port range to scan (e.g., 1-65535)
- **Required:** No
- **Default:** 1-1000

### timeout
- **Type:** integer
- **Description:** Connection timeout
- **Required:** No
- **Default:** 5
