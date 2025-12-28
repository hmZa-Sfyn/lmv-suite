#!/usr/bin/env python3
"""
js-enum - Simple & Effective JS Secrets Crawler
- Starts from any URL (page or direct .js)
- Crawls pages, follows links, finds all JS files
- Downloads & scans every JS file for secrets
- Shows matches with line numbers and context
- Single clean run (no duplicate output)
- Loads patterns from api_patterns.json
"""

import os
import signal
import sys
import re
import random
import json
from urllib.parse import urljoin, urlparse
from queue import Queue
from threading import Thread, Lock
import requests
from bs4 import BeautifulSoup

USER_AGENTS = [
    "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0 Safari/537.36",
]

# Global
discovered = set()
queue = Queue()
found_secrets = []
session = requests.Session()
lock = Lock()
max_depth = 3
threads = 10
filter_set = set()
should_stop = False

PATTERNS = {}

def load_patterns():
    global PATTERNS
    print("Loading patterns from api_patterns.json...")
    try:
        with open('api_patterns.json', 'r') as f:
            data = json.load(f)
        for name, regex in data.items():
            PATTERNS[name] = re.compile(regex)
        print(f"Loaded {len(PATTERNS)} patterns.")
    except FileNotFoundError:
        print("api_patterns.json not found. Using built-in.")
        PATTERNS.update({
            "Gemini API Key": re.compile(r"AIzaSyD[a-zA-Z0-9\\-_]{33}"),
            "DeepSeek API Key": re.compile(r"sk-[a-z0-9]{32}"),
            "Leonardo AI Key": re.compile(r"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"),
            "OpenAI Key": re.compile(r"sk-[a-zA-Z0-9]{48}"),
            "Supabase Anon Key": re.compile(r"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9\\.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6I[a-zA-Z0-9_-]{15,}"),
            "Hard-coded Key": re.compile(r"(GEMINI_API_KEY|DEEPSEEK_API_KEY|LEONARDO_API_KEY|OPENAI_API_KEY)[^\"']*[\"']([a-zA-Z0-9\\-_]{20,})[\"']", re.I),
            "Any sk- key": re.compile(r"sk-[a-zA-Z0-9]{20,}"),
            "Any AIza key": re.compile(r"AIzaSy[a-zA-Z0-9\\-_]{33}"),
        })

def signal_handler(sig, frame):
    global should_stop
    print("\n[!] Stopping...")
    should_stop = True

signal.signal(signal.SIGINT, signal_handler)

def extract_links(url, html):
    soup = BeautifulSoup(html, 'html.parser')
    links = set()
    for tag in soup.find_all(['a', 'script', 'link']):
        attr = tag.get('href') or tag.get('src')
        if attr:
            full = urljoin(url, attr)
            parsed = urlparse(full)
            clean = parsed._replace(fragment='', query='').geturl()
            if parsed.netloc == urlparse(url).netloc:
                links.add(clean)
    return links

def scan_js(url, content):
    lines = content.splitlines()
    for name, pat in PATTERNS.items():
        for match in pat.finditer(content):
            secret = match.group(1) if match.groups() else match.group(0)
            if len(secret) < 20: continue
            start_pos = match.start()
            line_no = content[:start_pos].count('\n') + 1
            line = lines[line_no - 1].strip()
            print(f"   \\_ ⚠️ {name}")
            print(f"      → {secret}")
            print(f"      Line {line_no}: {line[:150]}{'...' if len(line) > 150 else ''}")
            found_secrets.append({"service": name, "secret": secret, "url": url, "line": line_no})

def worker():
    while not should_stop:
        try:
            url, depth = queue.get(timeout=1)
        except:
            break

        tag = "JS" if url.lower().endswith('.js') else "PAGE"

        if filter_set and tag not in filter_set: 
            queue.task_done()
            continue

        print(f"{'|_ ' if depth > 0 else ''}{tag} → {url}")

        with lock:
            if url in discovered:
                queue.task_done()
                continue
            discovered.add(url)

        session.headers['User-Agent'] = random.choice(USER_AGENTS)

        try:
            r = session.get(url, timeout=10)
            if not r.ok:
                queue.task_done()
                continue

            if url.lower().endswith('.js'):
                scan_js(url, r.text)

            if 'text/html' in r.headers.get('content-type', '') and depth < max_depth:
                links = extract_links(url, r.text)
                with lock:
                    for link in links:
                        if link not in discovered:
                            queue.put((link, depth + 1))
        except Exception as e:
            pass

        queue.task_done()

def main():
    global max_depth, threads, filter_set

    url = os.getenv('ARG_URL')
    if not url:
        print("[!] ARG_URL required")
        sys.exit(1)

    url = url if url.startswith(('http://', 'https://')) else 'https://' + url
    max_depth = int(os.getenv('ARG_DEPTH', '3'))
    threads = int(os.getenv('ARG_THREADS', '10'))
    filter_str = os.getenv('ARG_FILTER', '')
    if filter_str:
        filter_set = {f.upper() for f in filter_str.split(',')}

    load_patterns()

    print(f"js-enum → {url}")
    print(f"Depth {max_depth} | Threads {threads}")

    discovered.add(url)
    queue.put((url, 0))

    workers = [Thread(target=worker, daemon=True) for _ in range(threads)]
    for w in workers: w.start()

    queue.join()

    print("\n[+] Done")
    print(f"Assets crawled: {len(discovered)} | Secrets found: {len(found_secrets)}")

if __name__ == '__main__':
    main()