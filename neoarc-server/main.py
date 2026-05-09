import logging
import os
import traceback

from datetime import timedelta

from flask import (
    Flask, g, request, jsonify, render_template,
    redirect, url_for, session, flash
)

from core.admin import admin_bp
from core.auth import generate_csrf_token
from core.db import init_db
from core.routes import register_route_blueprints

import config

logging.basicConfig(
    level=logging.INFO,
    format="[%(asctime)s] %(levelname)s %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
logger = logging.getLogger("neoarc")

app = Flask(__name__)
app.register_blueprint(admin_bp, url_prefix=config.ADMIN_ROUTE)
app.secret_key = config.SECRET_KEY

app.config['UPLOAD_FOLDER'] = config.UPLOAD_FOLDER
app.config['MAX_CONTENT_LENGTH'] = config.MAX_CONTENT_LENGTH
app.config['PERMANENT_SESSION_LIFETIME'] = timedelta(hours=24)

app.jinja_env.globals['csrf_token'] = generate_csrf_token

register_route_blueprints(app)

if not os.path.exists(config.UPLOAD_FOLDER):
    os.makedirs(config.UPLOAD_FOLDER)

init_db()


@app.before_request
def make_session_permanent():
    session.permanent = True


@app.teardown_appcontext
def close_db(exception):
    db = g.pop('_db', None)
    if db is not None:
        db.close()


@app.errorhandler(Exception)
def unhandled_exception(e):
    logger.error("Unhandled exception: %s\n%s", e, traceback.format_exc())
    if request.path.startswith('/api/'):
        return jsonify({'success': False, 'message': 'Internal server error'}), 500
    return render_template(
        'error.html', code='500', message='SYSTEM_FAILURE',
        subtext='A critical error occurred within the mainframe. Engineers have been notified.'
    ), 500


@app.errorhandler(404)
def page_not_found(e):
    return render_template(
        'error.html', code='404', message='NODE_NOT_FOUND',
        subtext='The requested sector has been purged or does not exist in the matrix.'
    ), 404


@app.errorhandler(413)
def request_entity_too_large(e):
    flash('Uploaded file exceeds the maximum allowed size (5 MB).', 'error')
    return redirect(request.referrer or url_for('index'))


@app.route('/')
def index():
    if 'user_id' in session:
        return redirect(url_for('alias.dashboard'))
    return render_template('login.html')


if __name__ == '__main__':
    logger.info("Starting NeoArc server on %s:%s", config.HOST, config.PORT)
    app.run(host=config.HOST, port=config.PORT, debug=config.DEBUG)
