#!/usr/bin/env python3
import os, sys, requests
ip = os.getenv('ARG_IP', '')
if not ip:
    print('[!] IP required')
    sys.exit(1)
try:
    r = requests.get(f'http://ip-api.com/json/{ip}', timeout=5)
    data = r.json()
    if data['status'] == 'success':
        lat = data.get("lat")
        lon = data.get("lon")
        maps_url = f'https://www.google.com/maps/search/{lat},{lon}'
        
        print(f'[+] IP: {data.get("query")}')
        print(f'[+] Country: {data.get("country")}')
        print(f'[+] Region: {data.get("regionName")}')
        print(f'[+] City: {data.get("city")}')
        print(f'[+] ZIP: {data.get("zip")}')
        print(f'[+] Coordinates: {lat}, {lon}')
        print(f'[+] Google Maps: {maps_url}')
        print(f'[+] ISP: {data.get("isp")}')
        print(f'[+] Organization: {data.get("org")}')
        print(f'[+] AS: {data.get("as")}')
        print(f'[+] Timezone: {data.get("timezone")}')
        print(f'[+] Mobile: {data.get("mobile")}')
        print(f'[+] Proxy: {data.get("proxy")}')
        print(f'[+] Hosting: {data.get("hosting")}')
    else:
        print('[!] Could not resolve IP')
except Exception as e:
    print(f'[!] Error: {e}')
