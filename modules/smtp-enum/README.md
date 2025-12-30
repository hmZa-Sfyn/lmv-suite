# smtp-enum

## Description
Enumerate valid email users on SMTP servers using VRFY/EXPN/RCPT

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
smtp, enumeration, email


## Links

## Options

### host
- **Type:** string
- **Description:** SMTP server address
- **Required:** Yes

### port
- **Type:** integer
- **Description:** SMTP port
- **Required:** No
- **Default:** 25

### wordlist
- **Type:** string
- **Description:** User wordlist file
- **Required:** Yes

### method
- **Type:** string
- **Description:** Enumeration method (vrfy, expn, rcpt)
- **Required:** No
- **Default:** rcpt
