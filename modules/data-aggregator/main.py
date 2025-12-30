#!/usr/bin/env python3
"""
Data Aggregator Module
Aggregate harvested data from multiple sources
"""

import os
import sys
import json
import csv
from pathlib import Path

try:
    import requests
except ImportError:
    print("Error: requests required")
    sys.exit(1)

def load_source(source):
    """Load data from file or URL"""
    data = []
    
    try:
        if source.startswith('http://') or source.startswith('https://'):
            # URL source
            response = requests.get(source, timeout=5)
            
            if source.endswith('.json'):
                data = response.json()
            elif source.endswith('.csv'):
                import io
                reader = csv.DictReader(io.StringIO(response.text))
                data = list(reader)
            else:
                data = [{'content': response.text}]
        else:
            # File source
            if not os.path.exists(source):
                print(f"[!] Source not found: {source}")
                return []
            
            with open(source, 'r', encoding='utf-8', errors='ignore') as f:
                if source.endswith('.json'):
                    data = json.load(f)
                elif source.endswith('.csv'):
                    reader = csv.DictReader(f)
                    data = list(reader)
                elif source.endswith('.txt'):
                    data = [{'content': line.strip()} for line in f if line.strip()]
                else:
                    data = [{'content': f.read()}]
        
        return data if isinstance(data, list) else [data]
    except Exception as e:
        print(f"[!] Error loading {source}: {e}")
        return []

def merge_data(all_data, strategy='append'):
    """Merge data from multiple sources"""
    if strategy == 'unique':
        # Remove duplicates
        seen = set()
        unique_data = []
        
        for item in all_data:
            item_str = json.dumps(item, sort_keys=True)
            if item_str not in seen:
                seen.add(item_str)
                unique_data.append(item)
        
        return unique_data
    
    elif strategy == 'merge':
        # Merge dictionaries
        merged = {}
        for item in all_data:
            if isinstance(item, dict):
                merged.update(item)
        return [merged] if merged else all_data
    
    else:  # append
        return all_data

def format_output(data, format_type='json'):
    """Format data for output"""
    if format_type == 'json':
        return json.dumps(data, indent=2)
    
    elif format_type == 'csv':
        if not data or not isinstance(data[0], dict):
            return "Cannot convert to CSV: data must be list of dicts"
        
        import io
        output = io.StringIO()
        writer = csv.DictWriter(output, fieldnames=data[0].keys())
        writer.writeheader()
        writer.writerows(data)
        return output.getvalue()
    
    elif format_type == 'html':
        html = "<html><body><table border='1'>\n"
        
        if data and isinstance(data[0], dict):
            # Header
            html += "<tr>"
            for key in data[0].keys():
                html += f"<th>{key}</th>"
            html += "</tr>\n"
            
            # Rows
            for item in data:
                html += "<tr>"
                for value in item.values():
                    html += f"<td>{value}</td>"
                html += "</tr>\n"
        
        html += "</table></body></html>"
        return html
    
    else:  # text
        result = ""
        for item in data:
            result += str(item) + "\n"
        return result

def main():
    sources = os.getenv('ARG_SOURCES')
    output_format = os.getenv('ARG_OUTPUT_FORMAT', 'json').lower()
    output_file = os.getenv('ARG_OUTPUT_FILE')
    merge_strategy = os.getenv('ARG_MERGE_STRATEGY', 'append').lower()
    
    if not sources:
        print("Error: sources required (comma-separated)")
        sys.exit(1)
    
    source_list = [s.strip() for s in sources.split(',')]
    
    print(f"[*] Data Aggregator")
    print(f"[*] Sources: {len(source_list)}")
    print(f"[*] Merge strategy: {merge_strategy}")
    print(f"[*] Output format: {output_format}")
    
    # Load all sources
    all_data = []
    for source in source_list:
        print(f"[*] Loading: {source}")
        data = load_source(source)
        all_data.extend(data)
    
    if not all_data:
        print(f"[!] No data loaded from sources")
        sys.exit(1)
    
    print(f"[+] Loaded {len(all_data)} items")
    
    # Merge
    print(f"[*] Merging data...")
    merged_data = merge_data(all_data, merge_strategy)
    
    print(f"[+] Final data: {len(merged_data)} items")
    
    # Format
    output = format_output(merged_data, output_format)
    
    # Save or print
    if output_file:
        try:
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(output)
            print(f"[+] Saved to: {output_file}")
            print(f"[+] File size: {os.path.getsize(output_file)} bytes")
        except Exception as e:
            print(f"[!] Error saving: {e}")
            sys.exit(1)
    else:
        print(f"\n[+] Output:")
        print(output[:1000])
        if len(output) > 1000:
            print(f"\n... ({len(output)} total characters)")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
