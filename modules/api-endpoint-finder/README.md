# api_endpoint_finder

## Description
Crawl web apps to discover and list potential API endpoints from JavaScript or network requests

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
api, discovery, reconnaissance


## Links

## Options

### url
- **Type:** string
- **Description:** Target URL to analyze
- **Required:** Yes

### depth
- **Type:** integer
- **Description:** Crawl depth
- **Required:** No
- **Default:** 1

### timeout
- **Type:** integer
- **Description:** Request timeout
- **Required:** No
- **Default:** 10

### api_patterns
- **Type:** boolean
- **Description:** Look for API endpoint patterns
- **Required:** No
- **Default:** True
