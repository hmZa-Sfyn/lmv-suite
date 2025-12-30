# data_aggregator

## Description
Aggregate harvested data from multiple sources (e.g., web, files) into structured formats like databases or reports

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
data, aggregation, reporting


## Links

## Options

### sources
- **Type:** string
- **Description:** Comma-separated list of source files or URLs
- **Required:** Yes

### output_format
- **Type:** string
- **Description:** Output format (json, csv, text, html)
- **Required:** No
- **Default:** json

### output_file
- **Type:** string
- **Description:** Output file path
- **Required:** No

### merge_strategy
- **Type:** string
- **Description:** How to merge data (append, merge, unique)
- **Required:** No
- **Default:** append
