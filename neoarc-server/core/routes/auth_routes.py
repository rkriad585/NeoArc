import hashlib
import secrets
import sqlite3

from datetime import datetime, timedelta

from flask import Blueprint, render_template, request, redirect, url_for, session, flash
from werkzeug.security import check_password_hash, generate_password_hash

from core.auth import csrf_required, regenerate_session
from core.helpers import db_execute, db_commit, logger, send_password_reset_email
from core.rate_limiter import rate_limiter
from core.validation import validate_email, validate_alias_name, validate_full_name, validate_password

import config

RESET_TOKEN_EXPIRY_HOURS = 1


def _hash_token(token: str) -> str:
    return hashlib.sha256(token.encode()).hexdigest()


auth_bp = Blueprint('auth', __name__)


@auth_bp.route('/login', methods=['POST'])
@csrf_required
def login():
    if rate_limiter.hit(f'login:{request.remote_addr}', max_attempts=5, window_seconds=60):
        flash('Too many login attempts. Please wait before trying again.', 'error')
        return redirect(url_for('index'))

    username = request.form.get('username', '').strip()
    password = request.form.get('password', '')

    if not username or not password:
        flash('Username and password are required.', 'error')
        return redirect(url_for('index'))

    try:
        user = db_execute(
            'SELECT * FROM users WHERE username=?', (username,)
        ).fetchone()
    except sqlite3.Error as e:
        from core.helpers import logger
        logger.error("Login DB error for user '%s': %s", username, e)
        flash('A system error occurred. Please try again.', 'error')
        return redirect(url_for('index'))

    if user and check_password_hash(user['password'], password):
        if user['is_blocked'] == 1:
            flash('ACCESS DENIED: Your account has been blocked by an Administrator.', 'error')
            return redirect(url_for('index'))

        regenerate_session()
        session['user_id'] = user['id']
        session['username'] = user['username']
        session['profile_pic'] = user['profile_pic'] if user['profile_pic'] else None
        flash('Authentication successful. Welcome to the NeoArc.', 'success')
        return redirect(url_for('alias.dashboard'))

    flash('Invalid credentials, matrix denied.', 'error')
    return redirect(url_for('index'))


@auth_bp.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        if rate_limiter.hit(f'register:{request.remote_addr}', max_attempts=3, window_seconds=300):
            flash('Too many registration attempts. Please wait before trying again.', 'error')
            return render_template('register.html')

        full_name = request.form.get('full_name', '').strip()
        email = request.form.get('email', '').strip()
        username = request.form.get('username', '').strip()
        password = request.form.get('password', '')

        errors = []
        if not validate_full_name(full_name):
            errors.append('Full name is required (max 200 characters).')
        if not validate_email(email):
            errors.append('A valid email address is required.')
        if not validate_alias_name(username):
            errors.append('Username must be 1-100 characters: letters, numbers, hyphens, underscores.')
        errors.extend(validate_password(password, username, email))

        if errors:
            for err in errors:
                flash(err, 'error')
            return render_template('register.html')

        hashed = generate_password_hash(password)
        try:
            db_execute(
                'INSERT INTO users (username, email, full_name, password) VALUES (?, ?, ?, ?)',
                (username, email, full_name, hashed)
            )
            db_commit()
            flash('Identity established. You may now login.', 'success')
            return redirect(url_for('index'))
        except sqlite3.IntegrityError:
            flash('Username or Email already exists in the grid.', 'error')
        except sqlite3.Error as e:
            from core.helpers import logger
            logger.error("Registration error: %s", e)
            flash('A system error occurred. Please try again.', 'error')

    return render_template('register.html')


@auth_bp.route('/logout')
def logout():
    session.clear()
    flash('Session terminated.', 'success')
    return redirect(url_for('index'))


@auth_bp.route('/forgot-password', methods=['GET', 'POST'])
def forgot_password():
    if request.method == 'POST':
        if rate_limiter.hit(f'forgot:{request.remote_addr}', max_attempts=3, window_seconds=300):
            flash('Too many attempts. Please wait before trying again.', 'error')
            return render_template('forgot_password.html')

        email = request.form.get('email', '').strip()
        if not validate_email(email):
            flash('A valid email address is required.', 'error')
            return render_template('forgot_password.html')

        try:
            user = db_execute(
                'SELECT id FROM users WHERE email=?', (email,)
            ).fetchone()
        except sqlite3.Error as e:
            logger.error("Forgot-password DB error: %s", e)
            flash('A system error occurred. Please try again.', 'error')
            return render_template('forgot_password.html')

        if user:
            token = secrets.token_urlsafe(32)
            token_hash = _hash_token(token)
            expires_at = (datetime.utcnow() + timedelta(hours=RESET_TOKEN_EXPIRY_HOURS)).isoformat()
            db_execute(
                'INSERT INTO reset_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)',
                (user['id'], token_hash, expires_at)
            )
            db_commit()

            reset_link = url_for('auth.reset_password', token=token, _external=True)
            user_row = db_execute(
                'SELECT full_name FROM users WHERE id=?', (user['id'],)
            ).fetchone()
            to_name = user_row['full_name'] if user_row and user_row['full_name'] else email

            ok = send_password_reset_email(email, to_name, reset_link)
            logger.info("Password reset email to %s: %s", email, "sent" if ok else "FAILED")

        flash('If that email is registered, a password reset link has been sent.', 'success')
        return redirect(url_for('index'))

    return render_template('forgot_password.html')


@auth_bp.route('/reset-password/<token>', methods=['GET', 'POST'])
def reset_password(token):
    if request.method == 'POST':
        password = request.form.get('password', '')
        confirm = request.form.get('confirm_password', '')

        if password != confirm:
            flash('Passwords do not match.', 'error')
            return render_template('reset_password.html', token=token)

        token_hash = _hash_token(token)
        try:
            row = db_execute(
                'SELECT id, user_id, expires_at FROM reset_tokens WHERE token_hash=? AND used=0',
                (token_hash,)
            ).fetchone()
        except sqlite3.Error as e:
            logger.error("Reset-password DB error: %s", e)
            flash('A system error occurred. Please try again.', 'error')
            return render_template('reset_password.html', token=token)

        if not row:
            flash('Invalid or expired reset token.', 'error')
            return redirect(url_for('index'))

        expires_at = datetime.fromisoformat(row['expires_at'])
        if datetime.utcnow() > expires_at:
            flash('Reset token has expired. Please request a new one.', 'error')
            return redirect(url_for('index'))

        user = db_execute(
            'SELECT username, email FROM users WHERE id=?', (row['user_id'],)
        ).fetchone()
        if not user:
            flash('User not found.', 'error')
            return redirect(url_for('index'))

        errors = validate_password(password, user['username'], user['email'])
        if errors:
            for err in errors:
                flash(err, 'error')
            return render_template('reset_password.html', token=token)

        hashed = generate_password_hash(password)
        db_execute('UPDATE users SET password=? WHERE id=?', (hashed, row['user_id']))
        db_execute('UPDATE reset_tokens SET used=1 WHERE id=?', (row['id'],))
        db_commit()

        flash('Password has been reset. You may now login with your new password.', 'success')
        return redirect(url_for('index'))

    return render_template('reset_password.html', token=token)
