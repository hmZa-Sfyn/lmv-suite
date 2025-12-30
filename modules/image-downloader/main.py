#!/usr/bin/env python3
"""
Image Downloader Module
Bulk download images from URLs or web pages
"""

import os
import sys
from urllib.parse import urljoin, urlparse
import requests
from pathlib import Path

try:
    from bs4 import BeautifulSoup
except ImportError:
    print("Error: beautifulsoup4 and requests required")
    sys.exit(1)

def extract_image_urls(html, base_url):
    """Extract image URLs from HTML"""
    soup = BeautifulSoup(html, 'html.parser')
    images = []
    
    for img in soup.find_all('img'):
        src = img.get('src')
        if src:
            full_url = urljoin(base_url, src)
            images.append({
                'url': full_url,
                'alt': img.get('alt', ''),
                'title': img.get('title', '')
            })
    
    return images

def download_image(url, output_dir, index, timeout=10):
    """Download a single image"""
    try:
        response = requests.get(url, timeout=timeout)
        response.raise_for_status()
        
        # Get filename from URL
        parsed = urlparse(url)
        filename = os.path.basename(parsed.path)
        
        if not filename:
            filename = f"image_{index}.jpg"
        
        filepath = os.path.join(output_dir, filename)
        
        # Save image
        with open(filepath, 'wb') as f:
            f.write(response.content)
        
        return True, filepath
    except Exception as e:
        return False, str(e)

def main():
    url = os.getenv('ARG_URL')
    output_dir = os.getenv('ARG_OUTPUT_DIR', './images')
    max_images = int(os.getenv('ARG_MAX_IMAGES', '10'))
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    
    if not url:
        print("Error: url required")
        sys.exit(1)
    
    print(f"[*] Image Downloader")
    print(f"[*] Target: {url}")
    print(f"[*] Output: {output_dir}, Max: {max_images}")
    
    # Create output directory
    try:
        Path(output_dir).mkdir(parents=True, exist_ok=True)
    except Exception as e:
        print(f"[!] Error creating directory: {e}")
        sys.exit(1)
    
    # Check if URL is direct image or HTML page
    if url.lower().endswith(('.jpg', '.jpeg', '.png', '.gif', '.webp')):
        print(f"[*] Direct image URL detected")
        success, result = download_image(url, output_dir, 1, timeout)
        
        if success:
            print(f"[+] Downloaded: {result}")
        else:
            print(f"[!] Error: {result}")
        
        sys.exit(0)
    
    # Fetch and parse HTML
    try:
        print(f"[*] Fetching page...")
        response = requests.get(url, timeout=timeout)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"[!] Error: {e}")
        sys.exit(1)
    
    print(f"[*] Extracting images...")
    images = extract_image_urls(response.text, url)
    
    if not images:
        print(f"[!] No images found")
        sys.exit(0)
    
    print(f"[+] Found {len(images)} image(s), downloading {min(max_images, len(images))}...")
    
    downloaded = 0
    for i, img in enumerate(images[:max_images], 1):
        success, result = download_image(img['url'], output_dir, i, timeout)
        
        if success:
            print(f"[+] {i}. Downloaded: {os.path.basename(result)}")
            downloaded += 1
        else:
            print(f"[!] {i}. Failed: {result[:50]}")
    
    print(f"\n[+] Complete: {downloaded} images downloaded to {output_dir}")

if __name__ == "__main__":
    main()
