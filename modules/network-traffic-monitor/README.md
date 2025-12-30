# network-traffic-monitor

## Description
Monitor and analyze network traffic in real-time, capturing and displaying packet statistics

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat
- **Version:** 1.0.0

## Tags
network, monitoring, traffic


## Links

## Options

### interface
- **Type:** string
- **Description:** Network interface to monitor
- **Required:** Yes

### duration
- **Type:** integer
- **Description:** Monitoring duration in seconds
- **Required:** No
- **Default:** 60

### filter
- **Type:** string
- **Description:** BPF filter (e.g., tcp port 80)
- **Required:** No

### packet_count
- **Type:** integer
- **Description:** Maximum packets to capture
- **Required:** No
- **Default:** 0
