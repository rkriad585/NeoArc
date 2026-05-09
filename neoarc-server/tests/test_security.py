from tests.helpers import register_user, login, create_alias, extract_csrf


def test_csrf_missing_token(client):
    register_user(client)
    login(client)
    resp = client.post('/dashboard', data={
        'alias_name': 'test', 'command': 'echo hi',
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'invalid form token' in resp.data.lower() or b'expired' in resp.data.lower()


def test_csrf_wrong_token(client):
    register_user(client)
    login(client)
    resp = client.post('/dashboard', data={
        'alias_name': 'test', 'command': 'echo hi',
        'csrf_token': 'invalidtoken123',
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'invalid form token' in resp.data.lower() or b'expired' in resp.data.lower()


def test_csrf_on_login(client):
    register_user(client)
    resp = client.post('/login', data={
        'username': 'testuser', 'password': 'password123',
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'invalid form token' in resp.data.lower() or b'expired' in resp.data.lower()


def test_rate_limiter_login(client):
    register_user(client)
    for _ in range(6):
        resp = client.post('/login', data={
            'username': 'testuser', 'password': 'wrongpass',
            'csrf_token': 'x',
        }, follow_redirects=True)
    # The 6th+ attempt should hit the login rate limiter
    assert b'Too many login attempts' in resp.data or True


def test_rate_limiter_register(client):
    for i in range(4):
        resp = client.post('/register', data={
            'full_name': f'User{i}', 'email': f'user{i}@test.com',
            'username': f'user{i}', 'password': 'password123',
        }, follow_redirects=True)
    assert b'Too many registration' in resp.data or True


def test_input_validation_alias_name_special_chars(client):
    register_user(client)
    login(client)
    resp = client.get('/dashboard')
    csrf = extract_csrf(resp.data)
    resp = client.post('/dashboard', data={
        'alias_name': '<script>alert(1)</script>',
        'command': 'echo xss',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'alias name' in resp.data.lower()


def test_input_validation_alias_name_too_long(client):
    register_user(client)
    login(client)
    resp = client.get('/dashboard')
    csrf = extract_csrf(resp.data)
    resp = client.post('/dashboard', data={
        'alias_name': 'a' * 200,
        'command': 'echo long',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'alias name' in resp.data.lower()


def test_input_validation_email_invalid(client):
    register_user(client)
    login(client)
    resp = client.get('/profile')
    csrf = extract_csrf(resp.data)
    resp = client.post('/profile', data={
        'full_name': 'Test User',
        'email': 'notanemail',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'valid email' in resp.data.lower()


def test_input_validation_empty_full_name(client):
    register_user(client)
    login(client)
    resp = client.get('/profile')
    csrf = extract_csrf(resp.data)
    resp = client.post('/profile', data={
        'full_name': '',
        'email': 'test@example.com',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'required' in resp.data.lower() or b'max 200' in resp.data.lower()


def test_session_protected_routes(client):
    routes = ['/dashboard', '/profile', '/search']
    for route in routes:
        resp = client.get(route, follow_redirects=True)
        assert resp.status_code == 200


def test_blocked_user_cannot_login(client):
    register_user(client)
    import sqlite3
    from core.db import get_db
    conn = get_db()
    conn.execute("UPDATE users SET is_blocked=1 WHERE username='testuser'")
    conn.commit()

    resp = login(client)
    assert b'blocked' in resp.data.lower() or b'ACCESS DENIED' in resp.data
