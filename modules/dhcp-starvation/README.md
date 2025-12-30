# dhcp-starvation

## Description
Exhaust DHCP pool by requesting multiple IP addresses rapidly

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
dos, dhcp, network


## Links

## Options

### interface
- **Type:** string
- **Description:** Network interface to use
- **Required:** Yes

### count
- **Type:** integer
- **Description:** Number of DHCP requests
- **Required:** No
- **Default:** 100

### timeout
- **Type:** integer
- **Description:** Timeout for each request
- **Required:** No
- **Default:** 5
