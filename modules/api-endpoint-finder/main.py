#!/usr/bin/env python3
"""
API Endpoint Finder Module
Discover API endpoints from web applications
"""

import os
import sys
import re
from urllib.parse import urljoin, urlparse

try:
    from bs4 import BeautifulSoup
    import requests
except ImportError:
    print("Error: beautifulsoup4 and requests required")
    sys.exit(1)

def find_endpoints_in_js(js_code):
    """Find API endpoints in JavaScript code"""
    endpoints = set()
    
    # Common API patterns
    patterns = [
        r'["\']\/api\/[^"\']+["\']',
        r'["\']\/v\d+\/[^"\']+["\']',
        r'["\']\/rest\/[^"\']+["\']',
        r'["\']\/service\/[^"\']+["\']',
        r'fetch\(["\']([^"\']+)["\']',
        r'axios\.(get|post|put|delete)\(["\']([^"\']+)["\']',
    ]
    
    for pattern in patterns:
        matches = re.findall(pattern, js_code)
        for match in matches:
            if isinstance(match, tuple):
                endpoint = match[-1] if match[-1] else match[0]
            else:
                endpoint = match
            
            # Clean up
            endpoint = endpoint.strip('"\'')
            if endpoint.startswith('/'):
                endpoints.add(endpoint)
    
    return endpoints

def extract_js_from_page(html):
    """Extract JavaScript from HTML"""
    soup = BeautifulSoup(html, 'html.parser')
    js_code = ""
    
    # Inline scripts
    for script in soup.find_all('script'):
        if script.string:
            js_code += script.string + "\n"
    
    return js_code

def find_endpoints(url, timeout=10):
    """Find API endpoints in web app"""
    try:
        response = requests.get(url, timeout=timeout)
        response.raise_for_status()
    except:
        return set()
    
    endpoints = set()
    
    # Extract from JavaScript
    js_code = extract_js_from_page(response.text)
    endpoints.update(find_endpoints_in_js(js_code))
    
    # Look for links that might be API endpoints
    soup = BeautifulSoup(response.text, 'html.parser')
    for link in soup.find_all('a', href=True):
        href = link.get('href')
        if '/api' in href or '/v1' in href or '/v2' in href:
            endpoints.add(href)
    
    return endpoints

def main():
    url = os.getenv('ARG_URL')
    depth = int(os.getenv('ARG_DEPTH', '1'))
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    api_patterns = os.getenv('ARG_API_PATTERNS', 'true').lower() == 'true'
    
    if not url:
        print("Error: url required")
        sys.exit(1)
    
    print(f"[*] API Endpoint Finder")
    print(f"[*] Target: {url}")
    print(f"[*] Depth: {depth}, Timeout: {timeout}s")
    
    all_endpoints = set()
    visited = set()
    to_visit = [url]
    
    for d in range(depth):
        new_to_visit = []
        
        for current_url in to_visit:
            if current_url in visited:
                continue
            
            visited.add(current_url)
            print(f"[*] Scanning: {current_url}")
            
            endpoints = find_endpoints(current_url, timeout)
            all_endpoints.update(endpoints)
            
            # Find links for next level
            if d < depth - 1:
                try:
                    response = requests.get(current_url, timeout=timeout)
                    soup = BeautifulSoup(response.text, 'html.parser')
                    
                    for link in soup.find_all('a', href=True):
                        href = urljoin(current_url, link.get('href'))
                        base_domain = urlparse(url).netloc
                        href_domain = urlparse(href).netloc
                        
                        if base_domain == href_domain and href not in visited:
                            new_to_visit.append(href)
                except:
                    pass
        
        to_visit = new_to_visit[:5]  # Limit to avoid explosion
    
    if not all_endpoints:
        print(f"[!] No API endpoints found")
        sys.exit(0)
    
    print(f"\n[+] Found {len(all_endpoints)} endpoint(s):")
    print()
    
    for i, endpoint in enumerate(sorted(all_endpoints), 1):
        print(f"{i:3}. {endpoint}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
