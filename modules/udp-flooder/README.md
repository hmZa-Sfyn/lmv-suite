# udp-flooder

## Description
UDP-based flood tool targeting specific ports/services

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
dos, udp, stress-test


## Links

## Options

### target
- **Type:** string
- **Description:** Target IP address
- **Required:** Yes

### port
- **Type:** integer
- **Description:** Target port
- **Required:** Yes

### count
- **Type:** integer
- **Description:** Number of packets
- **Required:** No
- **Default:** 1000

### payload_size
- **Type:** integer
- **Description:** Payload size in bytes
- **Required:** No
- **Default:** 512
