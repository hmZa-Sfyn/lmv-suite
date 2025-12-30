# quic-probe

## Description
Probe and analyze QUIC/HTTP3 traffic characteristics and server support

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
quic, http3, network


## Links

## Options

### host
- **Type:** string
- **Description:** Target host
- **Required:** Yes

### port
- **Type:** integer
- **Description:** QUIC port (usually 443)
- **Required:** No
- **Default:** 443

### timeout
- **Type:** integer
- **Description:** Connection timeout
- **Required:** No
- **Default:** 5

### probe_type
- **Type:** string
- **Description:** Probe type (handshake, version, initial)
- **Required:** No
- **Default:** handshake
