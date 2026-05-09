import subprocess, time, re, sys, urllib.request, urllib.parse

proc = subprocess.Popen(
    [sys.executable, 'main.py'],
    stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
    text=True
)
time.sleep(3)

try:
    resp = urllib.request.urlopen('http://localhost:59248/forgot-password')
    html = resp.read().decode()
    m = re.search(r'name="csrf_token" value="([^"]+)"', html)
    csrf = m.group(1) if m else 'NONE'

    data = urllib.parse.urlencode({'email': 'mdriyadkhan585@gmail.com', 'csrf_token': csrf}).encode()
    req = urllib.request.Request('http://localhost:59248/forgot-password', data=data)
    resp2 = urllib.request.urlopen(req)
    print(f'POST status: {resp2.status}')
    print(f'Redirected to: {resp2.url}')

    time.sleep(1)
    proc.terminate()
    out, _ = proc.communicate(timeout=5)
    print('\n--- Server logs ---')
    for line in out.splitlines():
        if any(w in line.lower() for w in ['email', 'reset', 'smtp', 'error', 'fail']):
            print(line)
finally:
    if proc.poll() is None:
        proc.terminate()
        proc.wait(timeout=5)
