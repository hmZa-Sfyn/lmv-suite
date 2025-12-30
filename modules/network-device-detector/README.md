# network-device-detector

## Description
Detect and identify devices connected to the network with OS fingerprinting capabilities

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat
- **Version:** 1.0.0

## Tags
network, discovery, fingerprinting


## Links

## Options

### network
- **Type:** string
- **Description:** Network range (e.g., 192.168.1.0/24)
- **Required:** Yes

### threads
- **Type:** integer
- **Description:** Number of threads for scanning
- **Required:** No
- **Default:** 20

### timeout
- **Type:** integer
- **Description:** Timeout per host
- **Required:** No
- **Default:** 2

### fingerprint
- **Type:** boolean
- **Description:** Fingerprint OS of discovered devices
- **Required:** No
- **Default:** True
