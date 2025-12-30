# email_scraper

## Description
Extract email addresses from text, HTML, or web pages using regex and validation for data harvesting

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
scraping, harvesting, email


## Links

## Options

### url
- **Type:** string
- **Description:** URL to scrape (or local file path)
- **Required:** Yes

### source_type
- **Type:** string
- **Description:** Source type (url, file, text)
- **Required:** No
- **Default:** url

### validate
- **Type:** boolean
- **Description:** Validate email format
- **Required:** No
- **Default:** True

### timeout
- **Type:** integer
- **Description:** Request timeout
- **Required:** No
- **Default:** 10
