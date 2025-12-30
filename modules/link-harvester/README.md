# link_harvester

## Description
Scrape and harvest hyperlinks from web content, with filtering for internal/external or domain-specific links

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
web, scraping, harvesting


## Links

## Options

### url
- **Type:** string
- **Description:** URL to scrape links from
- **Required:** Yes

### filter
- **Type:** string
- **Description:** Filter (all, internal, external)
- **Required:** No
- **Default:** all

### timeout
- **Type:** integer
- **Description:** Request timeout
- **Required:** No
- **Default:** 10

### unique
- **Type:** boolean
- **Description:** Return unique links only
- **Required:** No
- **Default:** True
