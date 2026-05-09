import sqlite3

from flask import Blueprint, request, jsonify

from core.auth import api_token_required
from core.cache import alias_cache
from core.helpers import db_execute, logger
from core.rate_limiter import rate_limiter

api_bp = Blueprint('api', __name__)


@api_bp.route('/api/alias/<alias_name>', methods=['GET'])
@api_token_required
def get_alias(alias_name):
    if rate_limiter.hit(f'api:{request.remote_addr}', max_attempts=60, window_seconds=60):
        return jsonify({'success': False, 'message': 'Rate limit exceeded'}), 429

    alias_name_normalized = alias_name.strip()

    cached, etag = alias_cache.get(alias_name_normalized)
    if cached:
        if_none_match = request.headers.get('If-None-Match', '')
        if if_none_match == etag:
            return '', 304
        resp = jsonify(cached)
        resp.headers['ETag'] = etag
        return resp

    try:
        result = db_execute(
            'SELECT command, exec_type FROM aliases WHERE alias_name=?', (alias_name_normalized,)
        ).fetchone()
    except sqlite3.Error as e:
        logger.error("API alias fetch error: %s", e)
        return jsonify({'success': False, 'message': 'Database error'}), 500

    if result:
        data = {
            'success': True,
            'alias': alias_name_normalized,
            'command': result['command'],
            'exec_type': result['exec_type']
        }
        etag = alias_cache.set(alias_name_normalized, data)
        resp = jsonify(data)
        resp.headers['ETag'] = etag
        return resp, 200

    return jsonify({'success': False, 'message': 'Alias not found'}), 404


@api_bp.route('/api/aliases', methods=['GET'])
@api_token_required
def list_aliases():
    if rate_limiter.hit(f'api:{request.remote_addr}', max_attempts=60, window_seconds=60):
        return jsonify({'success': False, 'message': 'Rate limit exceeded'}), 429

    try:
        results = db_execute(
            'SELECT alias_name FROM aliases ORDER BY alias_name'
        ).fetchall()
    except sqlite3.Error as e:
        logger.error("API aliases list error: %s", e)
        return jsonify({'success': False, 'message': 'Database error'}), 500

    return jsonify({
        'success': True,
        'aliases': [row['alias_name'] for row in results]
    })
