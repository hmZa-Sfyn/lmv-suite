# http-fuzzer

## Description
Basic HTTP request fuzzer to test web servers for vulnerabilities

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
web, fuzzing, vulnerability


## Links

## Options

### url
- **Type:** string
- **Description:** Target URL to fuzz
- **Required:** Yes

### wordlist
- **Type:** string
- **Description:** Path to wordlist file
- **Required:** No
- **Default:** /usr/share/wordlists/common.txt

### threads
- **Type:** integer
- **Description:** Number of threads
- **Required:** No
- **Default:** 5

### timeout
- **Type:** integer
- **Description:** Request timeout in seconds
- **Required:** No
- **Default:** 5
