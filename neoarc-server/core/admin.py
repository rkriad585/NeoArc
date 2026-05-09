from flask import Blueprint, render_template, request, redirect, url_for, session, flash
from werkzeug.security import check_password_hash
import sqlite3
import config
from core.db import get_db

admin_bp = Blueprint('admin', __name__)

@admin_bp.before_request
def require_admin():
    allowed_routes = ['admin.login']
    if request.endpoint in allowed_routes:
        return
    if not session.get('is_admin'):
        return redirect(url_for('admin.login'))

@admin_bp.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        email = request.form.get('email', '')
        password = request.form.get('password', '')

        if email == config.ADMIN_EMAIL and check_password_hash(config.ADMIN_PASSWORD_HASH, password):
            session['is_admin'] = True
            return redirect(url_for('admin.dashboard'))
        else:
            flash('Invalid Admin Credentials', 'error')

    return render_template('admin/login.html')

@admin_bp.route('/')
def dashboard():
    conn = get_db()
    users = conn.execute('''
        SELECT u.*, COUNT(a.id) as alias_count
        FROM users u
        LEFT JOIN aliases a ON u.id = a.user_id
        GROUP BY u.id
        ORDER BY u.id DESC
    ''').fetchall()
    return render_template('admin/dashboard.html', users=users)

@admin_bp.route('/user/<int:user_id>')
def view_user(user_id):
    conn = get_db()
    user = conn.execute("SELECT * FROM users WHERE id=?", (user_id,)).fetchone()
    aliases = conn.execute("SELECT * FROM aliases WHERE user_id=?", (user_id,)).fetchall()

    if not user:
        return "User not found", 404

    return render_template('admin/view_user.html', user=user, aliases=aliases)

@admin_bp.route('/block/<int:user_id>')
def block_user(user_id):
    conn = get_db()
    user = conn.execute("SELECT is_blocked FROM users WHERE id=?", (user_id,)).fetchone()
    if user:
        new_status = 0 if user['is_blocked'] else 1
        conn.execute("UPDATE users SET is_blocked=? WHERE id=?", (new_status, user_id))
        conn.commit()
    return redirect(url_for('admin.dashboard'))

@admin_bp.route('/delete/<int:user_id>')
def delete_user(user_id):
    conn = get_db()
    conn.execute("DELETE FROM aliases WHERE user_id=?", (user_id,))
    conn.execute("DELETE FROM users WHERE id=?", (user_id,))
    conn.commit()
    flash('User purged from database.', 'success')
    return redirect(url_for('admin.dashboard'))

@admin_bp.route('/delete_alias/<int:alias_id>/<int:user_id>')
def delete_user_alias(alias_id, user_id):
    conn = get_db()
    conn.execute("DELETE FROM aliases WHERE id=?", (alias_id,))
    conn.commit()
    return redirect(url_for('admin.view_user', user_id=user_id))

@admin_bp.route('/logout')
def logout():
    session.pop('is_admin', None)
    return redirect(url_for('admin.login'))
