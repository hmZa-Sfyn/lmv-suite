#!/usr/bin/env python3
"""
Form Extractor Module
Parse HTML from web pages to extract and analyze form tags
"""

import os
import sys
import json

try:
    from bs4 import BeautifulSoup
    import requests
except ImportError:
    print("Error: beautifulsoup4 and requests required. Install with: pip install beautifulsoup4 requests")
    sys.exit(1)

def extract_forms(html):
    """Extract all forms from HTML"""
    soup = BeautifulSoup(html, 'html.parser')
    forms = []
    
    for i, form in enumerate(soup.find_all('form'), 1):
        form_data = {
            'id': i,
            'action': form.get('action', 'N/A'),
            'method': form.get('method', 'GET').upper(),
            'name': form.get('name', 'N/A'),
            'enctype': form.get('enctype', 'application/x-www-form-urlencoded'),
            'inputs': []
        }
        
        # Extract all inputs
        for inp in form.find_all('input'):
            input_data = {
                'type': inp.get('type', 'text'),
                'name': inp.get('name', 'N/A'),
                'value': inp.get('value', ''),
                'required': inp.has_attr('required')
            }
            form_data['inputs'].append(input_data)
        
        # Extract textareas
        for textarea in form.find_all('textarea'):
            form_data['inputs'].append({
                'type': 'textarea',
                'name': textarea.get('name', 'N/A'),
                'value': textarea.string or '',
                'required': textarea.has_attr('required')
            })
        
        # Extract selects
        for select in form.find_all('select'):
            options = [opt.get_text() for opt in select.find_all('option')]
            form_data['inputs'].append({
                'type': 'select',
                'name': select.get('name', 'N/A'),
                'options': options,
                'required': select.has_attr('required')
            })
        
        forms.append(form_data)
    
    return forms

def main():
    url = os.getenv('ARG_URL')
    output_format = os.getenv('ARG_OUTPUT_FORMAT', 'text').lower()
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    
    if not url:
        print("Error: url required")
        sys.exit(1)
    
    print(f"[*] Form Extractor")
    print(f"[*] Target: {url}")
    
    try:
        print(f"[*] Fetching page...")
        response = requests.get(url, timeout=timeout)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"[!] Error fetching URL: {e}")
        sys.exit(1)
    
    print(f"[*] Parsing HTML...")
    forms = extract_forms(response.text)
    
    if not forms:
        print(f"[!] No forms found on page")
        sys.exit(0)
    
    print(f"\n[+] Found {len(forms)} form(s):")
    
    if output_format == "json":
        print(json.dumps(forms, indent=2))
    elif output_format == "csv":
        print("form_id,action,method,name,enctype,input_count")
        for form in forms:
            print(f"{form['id']},{form['action']},{form['method']},{form['name']},{form['enctype']},{len(form['inputs'])}")
    else:  # text
        for form in forms:
            print(f"\n[+] Form #{form['id']}:")
            print(f"    Action: {form['action']}")
            print(f"    Method: {form['method']}")
            print(f"    Name: {form['name']}")
            print(f"    Enctype: {form['enctype']}")
            print(f"    Inputs ({len(form['inputs'])}):")
            
            for inp in form['inputs']:
                if inp['type'] == 'select':
                    print(f"      - {inp['name']} ({inp['type']}) -> {inp['options']}")
                else:
                    print(f"      - {inp['name']} ({inp['type']}) = {inp['value']}")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
