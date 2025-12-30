#!/usr/bin/env python3
"""
HTTP Fuzzer Module
Basic HTTP request fuzzer to test web servers for vulnerabilities
"""

import os
import sys
import requests
from threading import Thread, Lock
from queue import Queue
from urllib.parse import urljoin

session = requests.Session()
results_lock = Lock()
found_resources = []

def fuzz_path(base_url, path, timeout, status_codes=[200, 201, 301, 302, 401, 403]):
    """Test a single path"""
    try:
        url = urljoin(base_url, path)
        response = session.get(url, timeout=timeout, allow_redirects=False)
        
        if response.status_code in status_codes:
            with results_lock:
                found_resources.append({
                    'url': url,
                    'status': response.status_code,
                    'length': len(response.content)
                })
            return True
    except requests.exceptions.RequestException:
        pass
    return False

def worker(queue, base_url, timeout):
    """Worker thread"""
    while True:
        path = queue.get()
        if path is None:
            break
        fuzz_path(base_url, path, timeout)
        queue.task_done()

def main():
    url = os.getenv('ARG_URL')
    wordlist = os.getenv('ARG_WORDLIST', '/usr/share/wordlists/common.txt')
    threads = int(os.getenv('ARG_THREADS', '5'))
    timeout = int(os.getenv('ARG_TIMEOUT', '5'))
    
    if not url:
        print("Error: url required")
        sys.exit(1)
    
    # Ensure URL ends with /
    if not url.endswith('/'):
        url += '/'
    
    print(f"[*] Fuzzing: {url}")
    print(f"[*] Wordlist: {wordlist}, Threads: {threads}")
    
    if not os.path.exists(wordlist):
        print(f"[!] Wordlist not found: {wordlist}")
        print("[*] Creating sample wordlist...")
        wordlist_data = ["admin", "test", "api", "config", "backup", "upload", "users", "products"]
    else:
        with open(wordlist, 'r') as f:
            wordlist_data = [line.strip() for line in f if line.strip()]
    
    queue = Queue()
    thread_list = []
    
    # Start workers
    for i in range(threads):
        t = Thread(target=worker, args=(queue, url, timeout))
        t.start()
        thread_list.append(t)
    
    # Queue paths
    for path in wordlist_data:
        queue.put(path)
    
    # Wait for completion
    queue.join()
    
    # Stop workers
    for i in range(threads):
        queue.put(None)
    for t in thread_list:
        t.join()
    
    print(f"\n[+] Found {len(found_resources)} resources:")
    for res in found_resources:
        print(f"    [{res['status']}] {res['url']} ({res['length']}B)")
    
    print("\n[+] Complete")

if __name__ == "__main__":
    main()
