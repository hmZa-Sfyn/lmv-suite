# icmp-redirect

## Description
Send fake ICMP redirect packets to manipulate routing tables

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
network, mitm, routing


## Links

## Options

### gateway
- **Type:** string
- **Description:** Current gateway IP
- **Required:** Yes

### target
- **Type:** string
- **Description:** Target IP to redirect
- **Required:** Yes

### redirect_to
- **Type:** string
- **Description:** IP to redirect traffic to
- **Required:** Yes

### interface
- **Type:** string
- **Description:** Network interface
- **Required:** No
- **Default:** eth0
