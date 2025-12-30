#!/usr/bin/env python3
"""
Fetch Analyzer Module
Scan JavaScript code to detect Fetch API usage
"""

import os
import sys
import re
import json

try:
    from bs4 import BeautifulSoup
    import requests
except ImportError:
    print("Error: beautifulsoup4 and requests required")
    sys.exit(1)

def extract_js_code(html):
    """Extract JavaScript from HTML"""
    soup = BeautifulSoup(html, 'html.parser')
    js_code = []
    
    # Extract inline script tags
    for script in soup.find_all('script'):
        if script.string:
            js_code.append(script.string)
    
    return '\n'.join(js_code)

def find_fetch_calls(js_code):
    """Find fetch() calls in JavaScript"""
    fetch_calls = []
    
    # Pattern to match fetch(...) calls
    pattern = r'fetch\s*\(\s*["\']([^"\']+)["\']'
    
    matches = re.finditer(pattern, js_code)
    for match in matches:
        url = match.group(1)
        fetch_calls.append({
            'url': url,
            'type': 'GET',
            'context': js_code[max(0, match.start()-50):min(len(js_code), match.end()+50)]
        })
    
    # Pattern for fetch with options (POST, etc)
    pattern_options = r'fetch\s*\(\s*["\']([^"\']+)["\'\s]*,\s*\{([^}]+)\}'
    
    matches = re.finditer(pattern_options, js_code)
    for match in matches:
        url = match.group(1)
        options = match.group(2)
        
        # Try to detect method
        method_match = re.search(r'method\s*:\s*["\'](\w+)["\']', options)
        method = method_match.group(1) if method_match else 'POST'
        
        fetch_calls.append({
            'url': url,
            'type': method,
            'options': options[:100],
            'context': js_code[max(0, match.start()-50):min(len(js_code), match.end()+50)]
        })
    
    return fetch_calls

def main():
    url = os.getenv('ARG_URL')
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    extract_inline = os.getenv('ARG_EXTRACT_INLINE', 'true').lower() == 'true'
    
    if not url:
        print("Error: url required")
        sys.exit(1)
    
    print(f"[*] Fetch API Analyzer")
    print(f"[*] Target: {url}")
    
    try:
        print(f"[*] Fetching page...")
        response = requests.get(url, timeout=timeout)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"[!] Error: {e}")
        sys.exit(1)
    
    print(f"[*] Extracting JavaScript...")
    js_code = extract_js_code(response.text)
    
    if not js_code:
        print(f"[!] No JavaScript found")
        sys.exit(0)
    
    print(f"[*] Analyzing Fetch API calls...")
    fetch_calls = find_fetch_calls(js_code)
    
    if not fetch_calls:
        print(f"[!] No Fetch API calls detected")
        sys.exit(0)
    
    print(f"\n[+] Found {len(fetch_calls)} Fetch API call(s):")
    print(f"\n{'#':<3} {'Method':<8} {'URL':<50}")
    print("-" * 70)
    
    for i, call in enumerate(fetch_calls, 1):
        url_short = call['url'][:47] if len(call['url']) > 47 else call['url']
        print(f"{i:<3} {call['type']:<8} {url_short:<50}")
    
    print(f"\n[+] Detailed Analysis:")
    for i, call in enumerate(fetch_calls, 1):
        print(f"\n[+] Call #{i}:")
        print(f"    URL: {call['url']}")
        print(f"    Method: {call['type']}")
        if 'options' in call:
            print(f"    Options: {call['options'][:100]}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
