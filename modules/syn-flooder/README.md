# syn-flooder

## Description
Launch SYN flood attacks to test DoS resilience (educational use only)

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
dos, tcp, stress-test


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
- **Description:** Number of SYN packets
- **Required:** No
- **Default:** 1000

### threads
- **Type:** integer
- **Description:** Number of threads
- **Required:** No
- **Default:** 5
