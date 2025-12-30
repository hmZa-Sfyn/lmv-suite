# quick-spider

## Description
Simple single-threaded web crawler - lists pages/files with types, saves to TOML progressively

## Metadata
- **Type:** python
- **Author:** hmza & Grok
- **Version:** 3.0.0

## Tags
web, crawler, simple, reconnaissance


## Links

## Options

### url
- **Type:** string
- **Description:** Target URL (e.g. http://testphp.vulnweb.com)
- **Required:** Yes

### depth
- **Type:** integer
- **Description:** Max crawl depth
- **Required:** No
- **Default:** 2

### output
- **Type:** string
- **Description:** Output TOML file
- **Required:** No
- **Default:** spider.toml
