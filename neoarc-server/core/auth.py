import functools
import secrets

from flask import session, request, jsonify, flash, redirect, url_for

import config


def generate_csrf_token():
    if '_csrf_token' not in session:
        session['_csrf_token'] = secrets.token_hex(32)
    return session['_csrf_token']


def csrf_required(f):
    @functools.wraps(f)
    def wrapper(*args, **kwargs):
        if request.method in ('POST', 'PUT', 'DELETE'):
            token = request.form.get('csrf_token') or request.headers.get('X-CSRF-Token')
            stored = session.get('_csrf_token')
            if not stored or not secrets.compare_digest(stored, token or ''):
                flash('Session expired or invalid form token. Please try again.', 'error')
                return redirect(request.referrer or url_for('index'))
        return f(*args, **kwargs)
    return wrapper


def api_token_required(f):
    @functools.wraps(f)
    def wrapper(*args, **kwargs):
        auth = request.headers.get('Authorization', '')
        if not auth.startswith('Bearer '):
            return jsonify({'success': False, 'message': 'Missing or invalid Authorization header'}), 401
        token = auth[len('Bearer '):]
        if not secrets.compare_digest(token, config.API_TOKEN):
            return jsonify({'success': False, 'message': 'Invalid API token'}), 403
        return f(*args, **kwargs)
    return wrapper


def regenerate_session():
    session_data = dict(session)
    session.clear()
    for k, v in session_data.items():
        session[k] = v
    session.permanent = True
