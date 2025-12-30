# http-interceptor

## Description
Proxy-like module to intercept and modify HTTP traffic in real-time

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
web, proxy, mitm


## Links

## Options

### listen_port
- **Type:** integer
- **Description:** Port to listen on
- **Required:** No
- **Default:** 8080

### target_host
- **Type:** string
- **Description:** Target host to proxy to
- **Required:** No

### mode
- **Type:** string
- **Description:** Mode (sniffer, modifier, logger)
- **Required:** No
- **Default:** logger
