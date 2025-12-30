#!/usr/bin/env python3
"""
PDF to Text Module
Convert PDF files to plain text or other formats
"""

import os
import sys

try:
    import PyPDF2
except ImportError:
    print("Error: PyPDF2 required. Install with: pip install PyPDF2")
    sys.exit(1)

def extract_text_from_pdf(pdf_path, page_range=None):
    """Extract text from PDF file"""
    try:
        with open(pdf_path, 'rb') as file:
            reader = PyPDF2.PdfReader(file)
            total_pages = len(reader.pages)
            
            # Parse page range
            start_page = 0
            end_page = total_pages
            
            if page_range:
                parts = page_range.split('-')
                if len(parts) == 2:
                    start_page = max(0, int(parts[0].strip()) - 1)
                    end_page = min(total_pages, int(parts[1].strip()))
            
            text = ""
            for page_num in range(start_page, end_page):
                page = reader.pages[page_num]
                text += f"\n--- Page {page_num + 1} ---\n"
                text += page.extract_text()
            
            return text, total_pages
    except Exception as e:
        return None, None

def convert_to_html(text, title="PDF Content"):
    """Convert text to HTML format"""
    html = f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{title}</title>
    <style>
        body {{ font-family: Arial, sans-serif; margin: 20px; }}
        .page {{ border: 1px solid #ccc; padding: 20px; margin: 20px 0; }}
        h1 {{ color: #333; }}
    </style>
</head>
<body>
    <h1>{title}</h1>
    <div class="content">
        <pre>{text}</pre>
    </div>
</body>
</html>"""
    return html

def convert_to_markdown(text):
    """Convert text to markdown format"""
    # Simple conversion
    lines = text.split('\n')
    markdown = ""
    
    for line in lines:
        if '---' in line:
            markdown += f"\n### {line}\n"
        else:
            markdown += line + "\n"
    
    return markdown

def main():
    input_file = os.getenv('ARG_INPUT_FILE')
    output_format = os.getenv('ARG_OUTPUT_FORMAT', 'text').lower()
    output_file = os.getenv('ARG_OUTPUT_FILE')
    page_range = os.getenv('ARG_PAGE_RANGE')
    
    if not input_file:
        print("Error: input_file required")
        sys.exit(1)
    
    if not os.path.exists(input_file):
        print(f"Error: File not found: {input_file}")
        sys.exit(1)
    
    print(f"[*] PDF to {output_format.upper()} Converter")
    print(f"[*] Input: {input_file}")
    
    if page_range:
        print(f"[*] Page range: {page_range}")
    
    print(f"[*] Extracting text...")
    
    text, total_pages = extract_text_from_pdf(input_file, page_range)
    
    if text is None:
        print(f"[!] Error extracting text from PDF")
        sys.exit(1)
    
    print(f"[+] Extracted from {total_pages} total pages")
    
    # Convert format
    if output_format == "html":
        output = convert_to_html(text, os.path.basename(input_file))
    elif output_format == "markdown":
        output = convert_to_markdown(text)
    else:
        output = text
    
    # Save or print
    if output_file:
        try:
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(output)
            print(f"[+] Saved to: {output_file}")
        except Exception as e:
            print(f"[!] Error saving file: {e}")
            sys.exit(1)
    else:
        print(f"\n[+] Extracted Content:")
        print("-" * 70)
        print(output[:1000])
        if len(output) > 1000:
            print(f"\n... ({len(output)} total characters)")
    
    print(f"\n[+] Complete")

if __name__ == "__main__":
    main()
