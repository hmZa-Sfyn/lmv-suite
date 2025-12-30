#!/usr/bin/env python3
"""
Email Scraper Module
Extract email addresses from text, HTML, or web pages
"""

import os
import sys
import re

try:
    from bs4 import BeautifulSoup
    import requests
except ImportError:
    print("Error: beautifulsoup4 and requests required")
    sys.exit(1)

EMAIL_REGEX = r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b'

def extract_emails(text, validate=True):
    """Extract and validate emails from text"""
    emails = set()
    
    matches = re.findall(EMAIL_REGEX, text)
    
    for email in matches:
        if validate:
            # Basic validation
            if len(email) < 5 or len(email) > 254:
                continue
            if email.startswith('.') or email.endswith('.'):
                continue
            if '..' in email:
                continue
        
        emails.add(email)
    
    return sorted(list(emails))

def main():
    source = os.getenv('ARG_URL')
    source_type = os.getenv('ARG_SOURCE_TYPE', 'url').lower()
    validate = os.getenv('ARG_VALIDATE', 'true').lower() == 'true'
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    
    if not source:
        print("Error: url required")
        sys.exit(1)
    
    print(f"[*] Email Scraper")
    print(f"[*] Source: {source}")
    print(f"[*] Type: {source_type}, Validate: {validate}")
    
    # Get content based on source type
    if source_type == "url":
        try:
            print(f"[*] Fetching page...")
            response = requests.get(source, timeout=timeout)
            response.raise_for_status()
            content = response.text
        except requests.exceptions.RequestException as e:
            print(f"[!] Error: {e}")
            sys.exit(1)
    
    elif source_type == "file":
        try:
            print(f"[*] Reading file...")
            with open(source, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()
        except Exception as e:
            print(f"[!] Error: {e}")
            sys.exit(1)
    
    elif source_type == "text":
        content = source
    
    else:
        print(f"Error: Unknown source type: {source_type}")
        sys.exit(1)
    
    # Extract HTML text if URL
    if source_type == "url":
        soup = BeautifulSoup(content, 'html.parser')
        # Remove script and style elements
        for script in soup(["script", "style"]):
            script.decompose()
        content = soup.get_text()
    
    print(f"[*] Extracting emails...")
    emails = extract_emails(content, validate)
    
    if not emails:
        print(f"[!] No emails found")
        sys.exit(0)
    
    print(f"\n[+] Found {len(emails)} email(s):")
    print()
    
    for email in emails:
        print(f"    {email}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
