# wifi-deauth-attack

## Description
Perform WiFi deauthentication attacks to test network security and disconnect devices

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat
- **Version:** 1.0.0

## Tags
wifi, deauth, security-test


## Links

## Options

### interface
- **Type:** string
- **Description:** Wireless interface name
- **Required:** Yes

### bssid
- **Type:** string
- **Description:** Target AP MAC address
- **Required:** Yes

### target_mac
- **Type:** string
- **Description:** Target client MAC (optional, all if not specified)
- **Required:** No

### count
- **Type:** integer
- **Description:** Number of deauth packets to send
- **Required:** No
- **Default:** 100

### interval
- **Type:** float
- **Description:** Interval between packets in seconds
- **Required:** No
- **Default:** 0.1
