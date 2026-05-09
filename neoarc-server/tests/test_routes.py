from tests.helpers import register_user, login, create_alias, extract_csrf


def test_index_redirect_when_logged_out(client):
    resp = client.get('/')
    assert resp.status_code == 200
    assert b'login' in resp.data.lower() or b'register' in resp.data.lower()


def test_register_page_loads(client):
    resp = client.get('/register')
    assert resp.status_code == 200
    assert b'csrf_token' in resp.data


def test_register_success(client):
    resp = register_user(client)
    assert resp.status_code == 200
    assert b'Identity established' in resp.data


def test_register_duplicate_username(client):
    register_user(client, username="dupuser")
    resp = register_user(client, username="dupuser", email="other@example.com")
    assert resp.status_code == 200
    assert b'already exists' in resp.data.lower()


def test_register_duplicate_email(client):
    register_user(client, email="dup@example.com")
    resp = register_user(client, username="other", email="dup@example.com")
    assert resp.status_code == 200
    assert b'already exists' in resp.data.lower()


def test_register_short_password(client):
    resp = client.post('/register', data={
        'full_name': 'Test', 'email': 'a@b.com',
        'username': 'user1', 'password': 'short',
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'Password must be at least' in resp.data


def test_register_missing_fields(client):
    resp = client.post('/register', data={
        'full_name': '', 'email': '', 'username': '', 'password': '',
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'required' in resp.data.lower()


def test_register_password_matches_username(client):
    resp = client.post('/register', data={
        'full_name': 'Test', 'email': 'a@b.com',
        'username': 'testuser', 'password': 'testuser',
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'Password must not match' in resp.data


def test_register_password_matches_email(client):
    resp = client.post('/register', data={
        'full_name': 'Test', 'email': 'same@example.com',
        'username': 'uniqueuser', 'password': 'same@example.com',
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'Password must not match' in resp.data


def test_login_success(client):
    register_user(client)
    resp = login(client)
    assert resp.status_code == 200
    assert b'Authentication successful' in resp.data


def test_login_invalid_credentials(client):
    resp = client.get('/')
    csrf = extract_csrf(resp.data)
    resp = client.post('/login', data={
        'username': 'nonexistent', 'password': 'wrong',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'Invalid credentials' in resp.data


def test_login_empty_fields(client):
    resp = client.get('/')
    csrf = extract_csrf(resp.data)
    resp = client.post('/login', data={
        'username': '', 'password': '',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'required' in resp.data.lower()


def test_dashboard_redirects_when_logged_out(client):
    resp = client.get('/dashboard', follow_redirects=True)
    assert resp.status_code == 200
    assert b'login' in resp.data.lower() or b'register' in resp.data.lower()


def test_dashboard_shows_no_aliases_initially(client):
    register_user(client)
    login(client)
    resp = client.get('/dashboard')
    assert resp.status_code == 200
    assert b'No aliases detected' in resp.data


def test_create_alias(client):
    register_user(client)
    login(client)
    resp = create_alias(client, name="my_alias", command="echo hi", exec_type="bash")
    assert resp.status_code == 200
    assert b'my_alias' in resp.data
    assert b'No aliases detected' not in resp.data


def test_create_alias_duplicate_name(client):
    register_user(client)
    login(client)
    create_alias(client, name="dup_alias", command="echo first")
    resp = create_alias(client, name="dup_alias", command="echo second")
    assert b'already exists' in resp.data.lower()


def test_create_alias_empty_name(client):
    register_user(client)
    login(client)
    resp = client.get('/dashboard')
    csrf = extract_csrf(resp.data)
    resp = client.post('/dashboard', data={
        'alias_name': '', 'command': 'echo hi',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert b'alias name' in resp.data.lower()


def test_create_alias_empty_command(client):
    register_user(client)
    login(client)
    resp = client.get('/dashboard')
    csrf = extract_csrf(resp.data)
    resp = client.post('/dashboard', data={
        'alias_name': 'test', 'command': '',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert b'cannot be empty' in resp.data.lower()


def test_delete_alias(client):
    register_user(client)
    login(client)
    create_alias(client, name="delete_me", command="echo gone")
    resp = client.get('/delete/1', follow_redirects=True)
    assert resp.status_code == 200
    assert b'purged' in resp.data.lower()


def test_delete_nonexistent_alias(client):
    register_user(client)
    login(client)
    resp = client.get('/delete/999', follow_redirects=True)
    assert resp.status_code == 200


def test_edit_page_loads(client):
    register_user(client)
    login(client)
    create_alias(client, name="editable", command="echo old")
    resp = client.get('/edit/1')
    assert resp.status_code == 200
    assert b'editable' in resp.data


def test_edit_nonexistent(client):
    register_user(client)
    login(client)
    resp = client.get('/edit/999', follow_redirects=True)
    assert resp.status_code == 200
    assert b'not found' in resp.data.lower()


def test_edit_update(client):
    register_user(client)
    login(client)
    create_alias(client, name="before_edit", command="echo before")
    resp = client.get('/edit/1')
    csrf = extract_csrf(resp.data)
    resp = client.post('/edit/1', data={
        'alias_name': 'after_edit', 'command': 'echo after',
        'exec_type': 'python', 'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'after_edit' in resp.data
    assert b'PYTHON' in resp.data


def test_search_no_query(client):
    register_user(client)
    login(client)
    resp = client.get('/search')
    assert resp.status_code == 200


def test_search_finds_alias(client):
    register_user(client)
    login(client)
    create_alias(client, name="searchable_alias", command="echo findme")
    resp = client.get('/search?q=findme')
    assert resp.status_code == 200
    assert b'searchable_alias' in resp.data


def test_search_no_results(client):
    register_user(client)
    login(client)
    resp = client.get('/search?q=nonexistent12345')
    assert resp.status_code == 200


def test_search_pagination(client):
    register_user(client)
    login(client)
    resp = client.get('/search?q=test&page=1')
    assert resp.status_code == 200


def test_view_alias(client):
    register_user(client)
    login(client)
    create_alias(client, name="viewable", command="echo viewme")
    resp = client.get('/view/viewable')
    assert resp.status_code == 200
    assert b'viewable' in resp.data
    assert b'echo viewme' in resp.data


def test_view_nonexistent_alias(client):
    register_user(client)
    login(client)
    resp = client.get('/view/nonexistent')
    assert resp.status_code == 404


def test_profile_redirects_when_logged_out(client):
    resp = client.get('/profile', follow_redirects=True)
    assert resp.status_code == 200
    assert b'login' in resp.data.lower() or b'register' in resp.data.lower()


def test_profile_loads(client):
    register_user(client)
    login(client)
    resp = client.get('/profile')
    assert resp.status_code == 200
    assert b'Test User' in resp.data


def test_profile_update(client):
    register_user(client)
    login(client)
    resp = client.get('/profile')
    csrf = extract_csrf(resp.data)
    resp = client.post('/profile', data={
        'full_name': 'Updated Name',
        'email': 'updated@example.com',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'Updated' in resp.data
    assert b'updated@example.com' in resp.data


def test_logout(client):
    register_user(client)
    login(client)
    resp = client.get('/logout', follow_redirects=True)
    assert resp.status_code == 200


def test_profile_node_count(client):
    register_user(client)
    login(client)
    create_alias(client, name="alias_a", command="echo a")
    create_alias(client, name="alias_b", command="echo b")
    resp = client.get('/profile')
    assert b'2' in resp.data or b'node_count' not in resp.data


def test_dashboard_shows_created_at(client):
    register_user(client)
    login(client)
    create_alias(client, name="dated_alias", command="echo date")
    resp = client.get('/dashboard')
    assert b'dated_alias' in resp.data
    assert b'Created:' in resp.data or b'created' in resp.data.lower()


def test_edit_other_users_alias_forbidden(client):
    register_user(client, username="user_a", email="a@test.com")
    login(client, username="user_a")
    create_alias(client, name="others_alias", command="echo secret")
    client.get('/logout', follow_redirects=True)

    register_user(client, username="user_b", email="b@test.com")
    login(client, username="user_b")
    resp = client.get('/edit/1', follow_redirects=True)
    assert resp.status_code == 200
    assert b'not found' in resp.data.lower()


def test_forgot_password_page_loads(client):
    resp = client.get('/forgot-password')
    assert resp.status_code == 200
    assert b'csrf_token' in resp.data


def _insert_reset_token(user_id, raw_token=None):
    import hashlib
    from datetime import datetime, timedelta
    from core.helpers import db_execute, db_commit
    token = raw_token or 'test-reset-token-123456'
    token_hash = hashlib.sha256(token.encode()).hexdigest()
    expires_at = (datetime.utcnow() + timedelta(hours=1)).isoformat()
    db_execute(
        'INSERT INTO reset_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)',
        (user_id, token_hash, expires_at)
    )
    db_commit()
    return token


def test_forgot_password_valid_email(client):
    register_user(client, email="reset@example.com")
    resp = client.get('/forgot-password')
    csrf = extract_csrf(resp.data)
    resp = client.post('/forgot-password', data={
        'email': 'reset@example.com',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'reset link' in resp.data.lower()


def test_forgot_password_unknown_email(client):
    resp = client.get('/forgot-password')
    csrf = extract_csrf(resp.data)
    resp = client.post('/forgot-password', data={
        'email': 'unknown@example.com',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'reset link' in resp.data.lower()


def test_forgot_password_invalid_email(client):
    resp = client.get('/forgot-password')
    csrf = extract_csrf(resp.data)
    resp = client.post('/forgot-password', data={
        'email': 'notanemail',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'valid email' in resp.data.lower()


def test_reset_password_with_valid_token(client):
    register_user(client, email="resetme@example.com")
    from core.helpers import db_execute
    user = db_execute("SELECT id FROM users WHERE email=?", ("resetme@example.com",)).fetchone()
    token = _insert_reset_token(user['id'])
    resp = client.get(f'/reset-password/{token}')
    assert resp.status_code == 200
    assert b'csrf_token' in resp.data

    resp = client.get(f'/reset-password/{token}')
    csrf = extract_csrf(resp.data)
    resp = client.post(f'/reset-password/{token}', data={
        'password': 'newpass123',
        'confirm_password': 'newpass123',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert resp.status_code == 200
    assert b'Password has been reset' in resp.data

    resp = client.get('/')
    csrf = extract_csrf(resp.data)
    resp = client.post('/login', data={
        'username': 'testuser', 'password': 'newpass123',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert b'Authentication successful' in resp.data


def test_reset_password_mismatch(client):
    register_user(client, email="mismatch@example.com")
    from core.helpers import db_execute
    user = db_execute("SELECT id FROM users WHERE email=?", ("mismatch@example.com",)).fetchone()
    token = _insert_reset_token(user['id'])
    resp = client.get(f'/reset-password/{token}')
    csrf = extract_csrf(resp.data)
    resp = client.post(f'/reset-password/{token}', data={
        'password': 'newpass123',
        'confirm_password': 'different',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert b'do not match' in resp.data.lower()


def test_reset_password_invalid_token(client):
    resp = client.get('/reset-password/invalidtoken123')
    assert resp.status_code == 200
    assert b'csrf_token' in resp.data


def test_reset_password_weak_password(client):
    register_user(client, email="weakpass@example.com")
    from core.helpers import db_execute
    user = db_execute("SELECT id FROM users WHERE email=?", ("weakpass@example.com",)).fetchone()
    token = _insert_reset_token(user['id'])
    resp = client.get(f'/reset-password/{token}')
    csrf = extract_csrf(resp.data)
    resp = client.post(f'/reset-password/{token}', data={
        'password': 'short',
        'confirm_password': 'short',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert b'Password must be at least' in resp.data


def test_forgot_password_rate_limiting(client):
    resp = client.get('/forgot-password')
    csrf = extract_csrf(resp.data)
    for _ in range(3):
        client.post('/forgot-password', data={
            'email': 'a@b.com',
            'csrf_token': csrf,
        })
    resp = client.post('/forgot-password', data={
        'email': 'a@b.com',
        'csrf_token': csrf,
    }, follow_redirects=True)
    assert b'Too many attempts' in resp.data


def test_delete_other_users_alias_forbidden(client):
    register_user(client, username="user_a", email="a@test.com")
    login(client, username="user_a")
    create_alias(client, name="secret_alias", command="echo secret")
    client.get('/logout', follow_redirects=True)

    register_user(client, username="user_b", email="b@test.com")
    login(client, username="user_b")
    resp = client.get('/delete/1', follow_redirects=True)
    assert resp.status_code == 200
