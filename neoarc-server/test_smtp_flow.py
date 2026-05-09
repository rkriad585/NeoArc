import re
import sys
sys.path.insert(0, '.')

from main import app
import config

print(f'SMTP_ENABLED: {config.SMTP_ENABLED}')
print(f'SMTP_USER: {config.SMTP_USER}')

with app.test_client() as c:
    resp = c.get('/forgot-password')
    m = re.search(rb'name="csrf_token" value="([^"]+)"', resp.data)
    csrf = m.group(1).decode() if m else 'none'
    print(f'CSRF token: {csrf}')

    resp = c.post('/forgot-password', data={
        'email': 'mdriyadkhan585@gmail.com',
        'csrf_token': csrf,
    }, follow_redirects=True)

    print(f'POST status: {resp.status_code}')
    print(f'Flash message present: {"reset link" in resp.data.lower().decode()}')
    print(f'Redirected to login: {"LOGIN" in resp.data.decode().upper() or "REGISTER" in resp.data.decode().upper()}')

    from core.helpers import db_execute
    row = db_execute('SELECT COUNT(*) FROM reset_tokens').fetchone()
    print(f'Reset tokens in DB: {row[0]}')
