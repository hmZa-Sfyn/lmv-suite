# port-knocker

## Description
Perform port knocking sequences to open firewall ports stealthily

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
network, firewall, stealth


## Links

## Options

### host
- **Type:** string
- **Description:** Target host
- **Required:** Yes

### sequence
- **Type:** string
- **Description:** Comma-separated port sequence (e.g., 1234,5678,9012)
- **Required:** Yes

### delay
- **Type:** float
- **Description:** Delay between knocks in seconds
- **Required:** No
- **Default:** 0.1

### protocol
- **Type:** string
- **Description:** Protocol to use (tcp, udp)
- **Required:** No
- **Default:** tcp
