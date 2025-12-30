#!/usr/bin/env python3
"""
Link Harvester Module
Scrape and harvest hyperlinks from web content
"""

import os
import sys
from urllib.parse import urljoin, urlparse

try:
    from bs4 import BeautifulSoup
    import requests
except ImportError:
    print("Error: beautifulsoup4 and requests required")
    sys.exit(1)

def harvest_links(html, base_url, filter_type="all", unique=True):
    """Extract links from HTML"""
    soup = BeautifulSoup(html, 'html.parser')
    links = []
    seen = set()
    
    base_domain = urlparse(base_url).netloc
    
    for link in soup.find_all('a', href=True):
        href = link.get('href')
        full_url = urljoin(base_url, href)
        
        # Skip duplicates if unique
        if unique and full_url in seen:
            continue
        
        seen.add(full_url)
        
        # Parse URL
        parsed = urlparse(full_url)
        is_internal = parsed.netloc == base_domain
        
        # Apply filter
        if filter_type == "internal" and not is_internal:
            continue
        elif filter_type == "external" and is_internal:
            continue
        
        links.append({
            'url': full_url,
            'text': link.get_text(strip=True)[:100],
            'type': 'internal' if is_internal else 'external',
            'title': link.get('title', '')
        })
    
    return links

def main():
    url = os.getenv('ARG_URL')
    filter_type = os.getenv('ARG_FILTER', 'all').lower()
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    unique = os.getenv('ARG_UNIQUE', 'true').lower() == 'true'
    
    if not url:
        print("Error: url required")
        sys.exit(1)
    
    print(f"[*] Link Harvester")
    print(f"[*] Target: {url}")
    print(f"[*] Filter: {filter_type}, Unique: {unique}")
    
    try:
        print(f"[*] Fetching page...")
        response = requests.get(url, timeout=timeout)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"[!] Error: {e}")
        sys.exit(1)
    
    print(f"[*] Harvesting links...")
    links = harvest_links(response.text, url, filter_type, unique)
    
    if not links:
        print(f"[!] No links found")
        sys.exit(0)
    
    # Categorize
    internal = [l for l in links if l['type'] == 'internal']
    external = [l for l in links if l['type'] == 'external']
    
    print(f"\n[+] Found {len(links)} link(s):")
    print(f"    Internal: {len(internal)}")
    print(f"    External: {len(external)}")
    
    print(f"\n{'#':<3} {'Type':<10} {'URL':<60}")
    print("-" * 75)
    
    for i, link in enumerate(links, 1):
        url_short = link['url'][:57] if len(link['url']) > 57 else link['url']
        print(f"{i:<3} {link['type']:<10} {url_short:<60}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
