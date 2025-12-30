#!/usr/bin/env python3
"""
Tag Matcher Module
Find and collect HTML tags matching specific patterns
"""

import os
import sys

try:
    from bs4 import BeautifulSoup
    import requests
except ImportError:
    print("Error: beautifulsoup4 and requests required")
    sys.exit(1)

def match_tags(html, tag, attribute=None):
    """Find matching tags in HTML"""
    soup = BeautifulSoup(html, 'html.parser')
    results = []
    
    try:
        # Try as CSS selector first
        tags = soup.select(tag)
    except:
        # Fall back to simple tag search
        tags = soup.find_all(tag)
    
    for elem in tags:
        if attribute:
            # Extract specific attribute
            attr_value = elem.get(attribute)
            if attr_value:
                results.append({
                    'tag': elem.name,
                    'attribute': attribute,
                    'value': attr_value,
                    'text': elem.get_text(strip=True)[:100]
                })
        else:
            # Get all attributes and text
            results.append({
                'tag': elem.name,
                'text': elem.get_text(strip=True)[:100],
                'attributes': dict(elem.attrs),
                'html': str(elem)[:200]
            })
    
    return results

def main():
    url = os.getenv('ARG_URL')
    tag = os.getenv('ARG_TAG')
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    attribute = os.getenv('ARG_ATTRIBUTE')
    
    if not url or not tag:
        print("Error: url and tag required")
        sys.exit(1)
    
    print(f"[*] Tag Matcher")
    print(f"[*] Target: {url}")
    print(f"[*] Searching for: {tag}")
    
    if attribute:
        print(f"[*] Extracting attribute: {attribute}")
    
    try:
        print(f"[*] Fetching page...")
        response = requests.get(url, timeout=timeout)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"[!] Error: {e}")
        sys.exit(1)
    
    print(f"[*] Matching tags...")
    results = match_tags(response.text, tag, attribute)
    
    if not results:
        print(f"[!] No matches found")
        sys.exit(0)
    
    print(f"\n[+] Found {len(results)} matching tag(s):")
    
    if attribute:
        print(f"\n{'#':<3} {'Attribute':<40} {'Text':<30}")
        print("-" * 75)
        
        for i, result in enumerate(results, 1):
            attr_val = result['value'][:37] if len(result['value']) > 37 else result['value']
            text = result['text'][:27] if len(result['text']) > 27 else result['text']
            print(f"{i:<3} {attr_val:<40} {text:<30}")
    else:
        print(f"\n{'#':<3} {'Tag':<10} {'Text':<50}")
        print("-" * 70)
        
        for i, result in enumerate(results, 1):
            text = result['text'][:47] if len(result['text']) > 47 else result['text']
            print(f"{i:<3} {result['tag']:<10} {text:<50}")
    
    # Show detailed info for first few
    if len(results) <= 5:
        print(f"\n[+] Detailed Information:")
        for i, result in enumerate(results, 1):
            print(f"\n[+] Match #{i}:")
            for key, value in result.items():
                if key not in ['html']:
                    if isinstance(value, dict):
                        print(f"    {key}: {len(value)} attributes")
                    else:
                        print(f"    {key}: {str(value)[:70]}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
