import re


def register_user(client, username="testuser", password="password123",
                  email="test@example.com", full_name="Test User"):
    resp = client.get('/register')
    csrf = extract_csrf(resp.data)
    return client.post('/register', data={
        'full_name': full_name,
        'email': email,
        'username': username,
        'password': password,
        'csrf_token': csrf,
    }, follow_redirects=True)


def login(client, username="testuser", password="password123"):
    resp = client.get('/')
    csrf = extract_csrf(resp.data)
    return client.post('/login', data={
        'username': username,
        'password': password,
        'csrf_token': csrf,
    }, follow_redirects=True)


def extract_csrf(html):
    m = re.search(rb'name="csrf_token" value="([^"]+)"', html)
    if m:
        return m.group(1).decode()
    return None


def create_alias(client, name="test_alias", command="echo hello",
                 exec_type="bash"):
    resp = client.get('/dashboard')
    csrf = extract_csrf(resp.data)
    return client.post('/dashboard', data={
        'alias_name': name,
        'command': command,
        'exec_type': exec_type,
        'csrf_token': csrf,
    }, follow_redirects=True)
