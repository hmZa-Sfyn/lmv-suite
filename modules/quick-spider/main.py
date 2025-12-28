#!/usr/bin/env python3
"""
Quick-Spider - Simple & Reliable Web Crawler
- Single-threaded for clean output (no duplicates)
- Discovers pages, files, directories
- Shows type + URL
- Saves clean list to TOML after every new discovery
- Ctrl+C safe
"""

import os
import signal
import sys
from urllib.parse import urljoin, urlparse, urlsplit
import requests
from bs4 import BeautifulSoup

# File types
FILE_TYPES = {
    '.js': '[JS]', '.css': '[CSS]', '.php': '[PHP]', '.asp': '[ASP]', '.aspx': '[ASPX]',
    '.html': '[HTML]', '.htm': '[HTML]', '.jpg': '[IMG]', '.jpeg': '[IMG]', '.png': '[IMG]',
    '.gif': '[IMG]', '.pdf': '[PDF]', '.zip': '[ZIP]', '.xml': '[XML]', '.json': '[JSON]',
    '.env': '[ENV]', '.bak': '[BAK]', '.sql': '[SQL]', '.git': '[GIT]',
}

def get_type(url):
    path = urlsplit(url).path.lower()
    for ext, tag in FILE_TYPES.items():
        if path.endswith(ext):
            return tag
    return '[DIR]' if url.endswith('/') else '[PAGE]'

# State
discovered = set()
to_crawl = []
output_file = "spider.toml"
start_url = ""
session = None

def save_results():
    path = os.path.expanduser(output_file)
    os.makedirs(os.path.dirname(path) or '.', exist_ok=True)
    with open(path, 'w') as f:
        f.write('# Quick-Spider Results\n')
        f.write(f'# Target: {start_url}\n')
        f.write(f'# Discovered: {len(discovered)}\n\n')
        f.write('urls = [\n')
        for url in sorted(discovered):
            f.write(f'  {{url = "{url}", type = "{get_type(url)}"}},\n')
        f.write(']\n')
    print(f"[+] Saved {len(discovered)} items → {path}")

def signal_handler(sig, frame):
    print("\n[!] Ctrl+C - Saving and exiting...")
    save_results()
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)

def extract_links(current_url):
    try:
        resp = session.get(current_url, timeout=10, allow_redirects=True)
        if 'text/html' not in resp.headers.get('content-type', ''):
            return []
        soup = BeautifulSoup(resp.text, 'html.parser')
        links = set()
        for tag in soup.find_all(True):
            for attr in ['href', 'src']:
                if tag.has_attr(attr):
                    raw = tag[attr].strip()
                    if raw:
                        full = urljoin(current_url, raw)
                        parsed = urlparse(full)
                        clean = parsed._replace(fragment='', query='').geturl()
                        if parsed.netloc == urlparse(start_url).netloc:
                            links.add(clean)
        return list(links)
    except Exception as e:
        print(f"    Error fetching {current_url}: {e}")
        return []

def crawl(target, depth_limit, out_file):
    global start_url, output_file, session, to_crawl
    start_url = target if target.startswith(('http://', 'https://')) else 'http://' + target
    output_file = out_file
    session = requests.Session()
    session.headers['User-Agent'] = 'QuickSpider/3.0'

    print(f"Starting simple crawl: {start_url}")
    print(f"Max depth: {depth_limit}\n")

    discovered.add(start_url)
    to_crawl = [(start_url, 0)]
    save_results()

    while to_crawl:
        url, depth = to_crawl.pop(0)
        print(f"Depth {depth} | {get_type(url)} {url}")

        if depth >= depth_limit:
            continue

        links = extract_links(url)
        for link in links:
            if link not in discovered:
                discovered.add(link)
                to_crawl.append((link, depth + 1))
                save_results()

    print("\nCrawl complete!")
    print(f"Total discovered: {len(discovered)}")
    save_results()

def main():
    url = os.getenv('ARG_URL')
    depth = int(os.getenv('ARG_DEPTH', '2'))
    output = os.getenv('ARG_OUTPUT', 'spider.toml')

    if not url:
        print("[!] ARG_URL required")
        sys.exit(1)

    crawl(url, depth, output)

if __name__ == '__main__':
    main()