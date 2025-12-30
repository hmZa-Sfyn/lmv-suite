# ipv6-ra-spoof

## Description
Spoof IPv6 Router Advertisements to perform MITM on IPv6 networks

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
ipv6, ra, mitm


## Links

## Options

### interface
- **Type:** string
- **Description:** Network interface
- **Required:** Yes

### target_mac
- **Type:** string
- **Description:** Target MAC address (optional)
- **Required:** No

### prefix
- **Type:** string
- **Description:** IPv6 prefix to advertise
- **Required:** No
- **Default:** fd00::/64

### lifetime
- **Type:** integer
- **Description:** Router lifetime in seconds
- **Required:** No
- **Default:** 3600
