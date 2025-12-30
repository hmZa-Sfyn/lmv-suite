# brute-ssh

## Description
Basic SSH brute-force module with wordlist support (for pentesting)

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
ssh, brute-force, security


## Links

## Options

### host
- **Type:** string
- **Description:** Target SSH host
- **Required:** Yes

### port
- **Type:** integer
- **Description:** SSH port
- **Required:** No
- **Default:** 22

### wordlist
- **Type:** string
- **Description:** Path to wordlist file
- **Required:** Yes

### username
- **Type:** string
- **Description:** Username to brute-force
- **Required:** Yes

### timeout
- **Type:** integer
- **Description:** Connection timeout
- **Required:** No
- **Default:** 5
