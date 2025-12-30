#!/usr/bin/env python3
"""
HTTP Interceptor Module
Intercept and modify HTTP traffic in real-time
"""

import os
import sys
import socket
import threading
from http.server import HTTPServer, BaseHTTPRequestHandler

class InterceptHandler(BaseHTTPRequestHandler):
    """HTTP request handler"""
    
    intercepted_requests = []
    
    def do_GET(self):
        """Handle GET requests"""
        self.log_request()
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b"<h1>Intercepted</h1><p>Request logged</p>")
    
    def do_POST(self):
        """Handle POST requests"""
        content_length = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(content_length)
        
        self.log_request(body)
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b"<h1>Intercepted</h1><p>POST logged</p>")
    
    def log_request(self, body=None):
        """Log intercepted request"""
        request_info = {
            'method': self.command,
            'path': self.path,
            'headers': dict(self.headers),
            'remote_addr': self.client_address[0]
        }
        
        if body:
            request_info['body'] = body.decode('utf-8', errors='ignore')
        
        InterceptHandler.intercepted_requests.append(request_info)
        
        print(f"\n[+] Intercepted {self.command} request:")
        print(f"    Path: {self.path}")
        print(f"    From: {self.client_address[0]}")
        print(f"    Headers: {len(self.headers)} headers")
        
        if body:
            print(f"    Body Length: {len(body)} bytes")
    
    def log_message(self, format, *args):
        """Suppress default logging"""
        pass

def run_server(port):
    """Run HTTP intercept server"""
    try:
        server = HTTPServer(('0.0.0.0', port), InterceptHandler)
        print(f"[*] HTTP Interceptor listening on port {port}")
        print(f"[*] Waiting for requests...\n")
        server.serve_forever()
    except KeyboardInterrupt:
        print(f"\n[*] Shutting down...")
        server.shutdown()
    except PermissionError:
        print(f"Error: Permission denied on port {port}")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

def main():
    listen_port = int(os.getenv('ARG_LISTEN_PORT', '8080'))
    target_host = os.getenv('ARG_TARGET_HOST')
    mode = os.getenv('ARG_MODE', 'logger').lower()
    
    print(f"[!] WARNING: This tool intercepts HTTP traffic")
    print(f"[!] Only use on networks you own or have permission to monitor")
    print(f"\n[*] HTTP Interceptor Module")
    print(f"[*] Mode: {mode}")
    print(f"[*] Listen Port: {listen_port}")
    
    if target_host:
        print(f"[*] Target Host: {target_host}")
    
    print(f"\n[*] To use this interceptor:")
    print(f"    1. Configure your client to use proxy: localhost:{listen_port}")
    print(f"    2. All HTTP requests will be logged")
    
    run_server(listen_port)

if __name__ == "__main__":
    main()
