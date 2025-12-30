# wifi-security-tester

## Description
Test WiFi network security by checking for weak passwords, common vulnerabilities, and WPS

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat
- **Version:** 1.0.0

## Tags
wifi, security, testing


## Links

## Options

### interface
- **Type:** string
- **Description:** Wireless interface
- **Required:** Yes

### bssid
- **Type:** string
- **Description:** Target AP MAC address
- **Required:** Yes

### wordlist
- **Type:** string
- **Description:** Password wordlist file
- **Required:** No

### test_wps
- **Type:** boolean
- **Description:** Test for WPS vulnerability
- **Required:** No
- **Default:** True
