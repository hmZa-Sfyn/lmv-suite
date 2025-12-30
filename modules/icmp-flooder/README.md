# icmp-flooder

## Description
Simple ICMP ping flood tool for stress testing network responsiveness

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
dos, icmp, stress-test


## Links

## Options

### target
- **Type:** string
- **Description:** Target IP address
- **Required:** Yes

### count
- **Type:** integer
- **Description:** Number of packets to send
- **Required:** No
- **Default:** 100

### packet_size
- **Type:** integer
- **Description:** Packet size in bytes
- **Required:** No
- **Default:** 56

### interval
- **Type:** float
- **Description:** Interval between packets in seconds
- **Required:** No
- **Default:** 0.1
