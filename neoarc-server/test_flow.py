import re, sys
sys.path.insert(0, '.')
from main import app

with app.test_request_context('/'):
    from core.helpers import db_execute

    user = db_execute('SELECT id, username, email FROM users WHERE email=?', ('mdriyadkhan585@gmail.com',)).fetchone()
    print(f'User found: {user is not None}')
    if user:
        print(f'  id={user["id"]} username={user["username"]} email=[{user["email"]}]')

with app.test_client() as c:
    resp = c.get('/forgot-password')
    m = re.search(rb'name="csrf_token" value="([^"]+)"', resp.data)
    csrf = m.group(1).decode()

    print(f'\nSubmitting forgot-password with email=mdriyadkhan585@gmail.com')
    resp = c.post('/forgot-password', data={
        'email': 'mdriyadkhan585@gmail.com',
        'csrf_token': csrf,
    }, follow_redirects=True)

    print(f'Response status: {resp.status_code}')
    text = resp.data.decode().lower()
    print(f'Reset link in flash: {"reset link" in text}')
    print(f'First 200 chars: {text[:200]}')

# Use another request context to query DB
with app.test_request_context('/'):
    from core.helpers import db_execute
    tokens = db_execute('SELECT id, user_id, used FROM reset_tokens').fetchall()
    print(f'\nTokens in DB: {len(tokens)}')
    for t in tokens:
        print(f'  id={t["id"]} user_id={t["user_id"]} used={t["used"]}')
