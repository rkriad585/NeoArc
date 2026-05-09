import os
import secrets
import sqlite3

from flask import Blueprint, render_template, request, redirect, url_for, session, flash

from core.helpers import db_execute, db_commit, logger
from core.validation import validate_email, validate_full_name, validate_image

profile_bp = Blueprint('profile', __name__)


@profile_bp.route('/profile', methods=['GET', 'POST'])
def profile():
    if 'user_id' not in session:
        return redirect(url_for('index'))

    if request.method == 'POST':
        token = request.form.get('csrf_token') or ''
        stored = session.get('_csrf_token')
        if not stored or not secrets.compare_digest(stored, token):
            flash('Session expired or invalid form token. Please try again.', 'error')
            return redirect(url_for('profile.profile'))

        full_name = request.form.get('full_name', '').strip()
        email = request.form.get('email', '').strip()

        if not validate_full_name(full_name):
            flash('Full name is required (max 200 characters).', 'error')
            return redirect(url_for('profile.profile'))

        if not validate_email(email):
            flash('A valid email address is required.', 'error')
            return redirect(url_for('profile.profile'))

        from flask import current_app
        file = request.files.get('profile_pic')
        img_err = validate_image(file) if file and file.filename else None
        if img_err:
            flash(img_err, 'error')
            return redirect(url_for('profile.profile'))

        try:
            filename = None
            if file and file.filename:
                ext = file.filename.rsplit('.', 1)[1].lower()
                filename = f'user_{session["user_id"]}.{ext}'
                upload_folder = current_app.config.get('UPLOAD_FOLDER', 'static/uploads')
                file.save(os.path.join(upload_folder, filename))
                session['profile_pic'] = filename
                db_execute(
                    'UPDATE users SET full_name=?, email=?, profile_pic=? WHERE id=?',
                    (full_name, email, filename, session['user_id'])
                )
            else:
                db_execute(
                    'UPDATE users SET full_name=?, email=? WHERE id=?',
                    (full_name, email, session['user_id'])
                )
            db_commit()
            flash('Identity profile updated successfully.', 'success')
        except sqlite3.IntegrityError:
            flash('Email is already bound to another node.', 'error')
        except sqlite3.Error as e:
            logger.error("Profile update error: %s", e)
            flash('A system error occurred. Please try again.', 'error')

        return redirect(url_for('profile.profile'))

    try:
        user = db_execute(
            'SELECT * FROM users WHERE id=?', (session['user_id'],)
        ).fetchone()
        node_count = db_execute(
            'SELECT COUNT(*) FROM aliases WHERE user_id=?', (session['user_id'],)
        ).fetchone()[0]
    except sqlite3.Error as e:
        logger.error("Profile fetch error: %s", e)
        flash('A system error occurred.', 'error')
        return redirect(url_for('alias.dashboard'))

    return render_template('profile.html', user=user, node_count=node_count)
