# 1. Classic port scan on one target
for $p in 1..1024 -> nmap -p$p -sV 10.10.10.50

# 2. Most common - sweep last octet
for $ip in 192.168.1.1..192.168.1.254 -> ping -c1 -W1 $ip

# 3. Small subnet sweep
for $ip in 10.10.50.100..10.10.50.150 -> nmap -sS -T4 $ip

# 4. Username / vhost / subdomain brute
for $u in admin\|root\|test\|backup\|dev -> hydra -l $u -P pass.txt ssh://10.10.10.10

# 5. Basic charset for password generation / testing
for $c in a..z+0..9 -> echo "trying: pass$c"

# 6. Hex digits (very useful)
for $h in 0..9+A..F -> echo "hex digit: $h"

# 7. Multiple ranges combined
for $c in a..z+A..Z+0..9 -> touch "file_$c.txt"