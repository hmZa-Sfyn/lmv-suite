# tag_matcher

## Description
Find and collect all HTML tags matching specific patterns (e.g., <a>, <p>, or custom selectors) from a webpage

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
html, scraping, parsing


## Links

## Options

### url
- **Type:** string
- **Description:** URL to scrape
- **Required:** Yes

### tag
- **Type:** string
- **Description:** HTML tag to match (e.g., a, p, div, or CSS selector)
- **Required:** Yes

### timeout
- **Type:** integer
- **Description:** Request timeout
- **Required:** No
- **Default:** 10

### attribute
- **Type:** string
- **Description:** Attribute to extract (optional)
- **Required:** No
