#!/usr/bin/env python3
"""
CSV to JSON Module
Convert CSV files to JSON format with data cleaning
"""

import os
import sys
import csv
import json

def csv_to_json(csv_path, delimiter=','):
    """Convert CSV to JSON"""
    try:
        data = []
        
        with open(csv_path, 'r', encoding='utf-8', errors='ignore') as f:
            reader = csv.DictReader(f, delimiter=delimiter)
            
            if not reader.fieldnames:
                return None
            
            for row in reader:
                # Clean data
                cleaned_row = {}
                for key, value in row.items():
                    if key:  # Skip empty keys
                        # Try to convert to appropriate type
                        if value:
                            if value.lower() in ['true', 'yes']:
                                cleaned_row[key] = True
                            elif value.lower() in ['false', 'no']:
                                cleaned_row[key] = False
                            elif value.isdigit():
                                cleaned_row[key] = int(value)
                            elif is_float(value):
                                cleaned_row[key] = float(value)
                            else:
                                cleaned_row[key] = value
                        else:
                            cleaned_row[key] = None
                
                data.append(cleaned_row)
        
        return data
    except Exception as e:
        print(f"Error: {e}")
        return None

def is_float(value):
    """Check if value is float"""
    try:
        float(value)
        return '.' in value
    except ValueError:
        return False

def main():
    input_file = os.getenv('ARG_INPUT_FILE')
    output_file = os.getenv('ARG_OUTPUT_FILE')
    delimiter = os.getenv('ARG_DELIMITER', ',')
    pretty_print = os.getenv('ARG_PRETTY_PRINT', 'true').lower() == 'true'
    
    if not input_file:
        print("Error: input_file required")
        sys.exit(1)
    
    if not os.path.exists(input_file):
        print(f"Error: File not found: {input_file}")
        sys.exit(1)
    
    print(f"[*] CSV to JSON Converter")
    print(f"[*] Input: {input_file}")
    print(f"[*] Delimiter: '{delimiter}'")
    
    print(f"[*] Converting...")
    
    data = csv_to_json(input_file, delimiter)
    
    if data is None or len(data) == 0:
        print(f"[!] Error converting file or no data found")
        sys.exit(1)
    
    print(f"[+] Converted {len(data)} records")
    
    # Generate JSON
    if pretty_print:
        json_output = json.dumps(data, indent=2)
    else:
        json_output = json.dumps(data)
    
    # Save or print
    if output_file:
        try:
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(json_output)
            print(f"[+] Saved to: {output_file}")
            print(f"[+] File size: {os.path.getsize(output_file)} bytes")
        except Exception as e:
            print(f"[!] Error saving file: {e}")
            sys.exit(1)
    else:
        print(f"\n[+] JSON Output:")
        print(json_output[:1000])
        if len(json_output) > 1000:
            print(f"\n... ({len(json_output)} total characters)")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
