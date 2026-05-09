import secrets
import sqlite3

from flask import Blueprint, render_template, request, redirect, url_for, session, flash

from core.cache import alias_cache
from core.helpers import db_execute, db_commit, logger
from core.validation import validate_alias_name

alias_bp = Blueprint('alias', __name__)


def _csrf_check() -> bool:
    token = request.form.get('csrf_token') or ''
    stored = session.get('_csrf_token')
    if not stored or not secrets.compare_digest(stored, token):
        flash('Session expired or invalid form token. Please try again.', 'error')
        return False
    return True


@alias_bp.route('/dashboard', methods=['GET', 'POST'])
def dashboard():
    if 'user_id' not in session:
        return redirect(url_for('index'))

    if request.method == 'POST':
        if not _csrf_check():
            return redirect(url_for('alias.dashboard'))

        alias_name = request.form.get('alias_name', '').strip()
        command = request.form.get('command', '').strip()
        exec_type = request.form.get('exec_type', '').strip()

        if not alias_name or not validate_alias_name(alias_name):
            flash('Alias name must be 1-100 characters: letters, numbers, hyphens, underscores.', 'error')
            return redirect(url_for('alias.dashboard'))

        if not command:
            flash('Command / code cannot be empty.', 'error')
            return redirect(url_for('alias.dashboard'))

        if len(command) > 100000:
            flash('Command exceeds maximum length (100 000 characters).', 'error')
            return redirect(url_for('alias.dashboard'))

        try:
            db_execute(
                'INSERT INTO aliases (alias_name, command, exec_type, user_id) VALUES (?, ?, ?, ?)',
                (alias_name, command, exec_type, session['user_id'])
            )
            db_commit()
            alias_cache.invalidate(alias_name)
            flash('Alias successfully compiled into NeoArc.', 'success')
        except sqlite3.IntegrityError:
            flash('Alias name already exists!', 'error')
        except sqlite3.Error as e:
            logger.error("Dashboard create alias error: %s", e)
            flash('A system error occurred. Please try again.', 'error')

    try:
        aliases = db_execute(
            'SELECT * FROM aliases WHERE user_id=? ORDER BY id DESC', (session['user_id'],)
        ).fetchall()
    except sqlite3.Error as e:
        logger.error("Dashboard fetch aliases error: %s", e)
        aliases = []

    return render_template('dashboard.html', aliases=aliases, username=session['username'])


@alias_bp.route('/edit/<int:id>', methods=['GET', 'POST'])
def edit_alias(id):
    if 'user_id' not in session:
        return redirect(url_for('index'))

    try:
        alias = db_execute(
            'SELECT * FROM aliases WHERE id=? AND user_id=?', (id, session['user_id'])
        ).fetchone()
    except sqlite3.Error as e:
        logger.error("Edit alias fetch error: %s", e)
        flash('A system error occurred.', 'error')
        return redirect(url_for('alias.dashboard'))

    if not alias:
        flash('Alias not found.', 'error')
        return redirect(url_for('alias.dashboard'))

    if request.method == 'POST':
        if not _csrf_check():
            return redirect(url_for('alias.dashboard'))

        alias_name = request.form.get('alias_name', '').strip()
        command = request.form.get('command', '').strip()
        exec_type = request.form.get('exec_type', '').strip()

        if not alias_name or not validate_alias_name(alias_name):
            flash('Alias name must be 1-100 characters: letters, numbers, hyphens, underscores.', 'error')
            return render_template('edit_alias.html', alias=alias)

        if not command:
            flash('Command / code cannot be empty.', 'error')
            return render_template('edit_alias.html', alias=alias)

        if len(command) > 100000:
            flash('Command exceeds maximum length (100 000 characters).', 'error')
            return render_template('edit_alias.html', alias=alias)

        try:
            db_execute(
                'UPDATE aliases SET alias_name=?, command=?, exec_type=? WHERE id=? AND user_id=?',
                (alias_name, command, exec_type, id, session['user_id'])
            )
            db_commit()
            alias_cache.invalidate(alias_name)
            alias_cache.invalidate(alias['alias_name'])
            flash('Alias updated successfully.', 'success')
            return redirect(url_for('alias.dashboard'))
        except sqlite3.IntegrityError:
            flash('Alias name already exists!', 'error')
        except sqlite3.Error as e:
            logger.error("Edit alias update error: %s", e)
            flash('A system error occurred. Please try again.', 'error')

    return render_template('edit_alias.html', alias=alias)


@alias_bp.route('/delete/<int:id>')
def delete_alias(id):
    if 'user_id' not in session:
        return redirect(url_for('index'))

    try:
        alias = db_execute(
            'SELECT alias_name FROM aliases WHERE id=? AND user_id=?', (id, session['user_id'])
        ).fetchone()
        if alias:
            alias_cache.invalidate(alias['alias_name'])
        db_execute('DELETE FROM aliases WHERE id=? AND user_id=?', (id, session['user_id']))
        db_commit()
        flash('Alias purged from the system.', 'success')
    except sqlite3.Error as e:
        logger.error("Delete alias error: %s", e)
        flash('A system error occurred. Please try again.', 'error')

    return redirect(url_for('alias.dashboard'))


@alias_bp.route('/search')
def search():
    if 'user_id' not in session:
        return redirect(url_for('index'))

    query = request.args.get('q', '').strip()
    page = request.args.get('page', 1, type=int)
    per_page = 18
    results = []
    total = 0

    if query:
        if len(query) > 200:
            query = query[:200]
        try:
            total = db_execute('''
                SELECT COUNT(*) FROM aliases a
                JOIN users u ON a.user_id = u.id
                WHERE a.alias_name LIKE ? OR a.command LIKE ?
            ''', ('%' + query + '%', '%' + query + '%')).fetchone()[0]

            offset = (page - 1) * per_page
            results = db_execute('''
                SELECT a.*, u.username, u.profile_pic
                FROM aliases a
                JOIN users u ON a.user_id = u.id
                WHERE a.alias_name LIKE ? OR a.command LIKE ?
                ORDER BY a.id DESC
                LIMIT ? OFFSET ?
            ''', ('%' + query + '%', '%' + query + '%', per_page, offset)).fetchall()
        except sqlite3.Error as e:
            logger.error("Search error: %s", e)

    total_pages = max(1, (total + per_page - 1) // per_page)

    return render_template(
        'search.html',
        results=results,
        query=query,
        page=page,
        total_pages=total_pages,
        total=total
    )


@alias_bp.route('/view/<alias_name>')
def view_alias(alias_name):
    if 'user_id' not in session:
        return redirect(url_for('index'))

    try:
        alias = db_execute('''
            SELECT a.*, u.username, u.profile_pic, u.email
            FROM aliases a
            JOIN users u ON a.user_id = u.id
            WHERE a.alias_name = ?
        ''', (alias_name,)).fetchone()
    except sqlite3.Error as e:
        logger.error("View alias error: %s", e)
        return render_template(
            'error.html', code='500', message='DB_ERROR',
            subtext='Could not query the matrix.'
        ), 500

    if not alias:
        return render_template(
            'error.html', code='404', message='ALIAS_MISSING',
            subtext='This node does not exist in the public matrix.'
        ), 404

    return render_template('view.html', alias=alias)
