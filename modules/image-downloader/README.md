# image_downloader

## Description
Bulk download images from URLs or web pages, with options for resizing and format conversion

## Metadata
- **Type:** python
- **Author:** l0n3ly_nat & hmza
- **Version:** 1.0.0

## Tags
web, images, download


## Links

## Options

### url
- **Type:** string
- **Description:** URL to scrape images from or direct image URL
- **Required:** Yes

### output_dir
- **Type:** string
- **Description:** Output directory for downloaded images
- **Required:** No
- **Default:** ./images

### max_images
- **Type:** integer
- **Description:** Maximum images to download
- **Required:** No
- **Default:** 10

### timeout
- **Type:** integer
- **Description:** Request timeout
- **Required:** No
- **Default:** 10
