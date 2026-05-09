import os
import re
import json

os.environ['NEOARC_API_TOKEN'] = 'test-api-token-32chars!!'
from tests.helpers import register_user, login, create_alias, extract_csrf
from main import app
from core.cache import alias_cache
import config


API_TOKEN = 'test-api-token-32chars!!'


def _api_headers():
    return {'Authorization': f'Bearer {API_TOKEN}'}


def test_api_missing_token(client):
    resp = client.get('/api/alias/test')
    assert resp.status_code == 401
    data = resp.get_json()
    assert data['success'] is False


def test_api_invalid_token(client):
    resp = client.get('/api/alias/test', headers={'Authorization': 'Bearer wrongtoken'})
    assert resp.status_code == 403


def test_api_invalid_auth_scheme(client):
    resp = client.get('/api/alias/test', headers={'Authorization': 'Basic xyz'})
    assert resp.status_code == 401


def test_api_alias_not_found(client):
    resp = client.get('/api/alias/nonexistent', headers=_api_headers())
    assert resp.status_code == 404
    data = resp.get_json()
    assert data['success'] is False


def test_api_alias_found(client):
    register_user(client)
    login(client)
    create_alias(client, name="api_test_alias", command="echo api_ok", exec_type="bash")

    resp = client.get('/api/alias/api_test_alias', headers=_api_headers())
    assert resp.status_code == 200
    data = resp.get_json()
    assert data['success'] is True
    assert data['alias'] == 'api_test_alias'
    assert data['command'] == 'echo api_ok'
    assert data['exec_type'] == 'bash'


def test_api_etag_caching(client):
    register_user(client)
    login(client)
    create_alias(client, name="etag_alias", command="echo etag")

    resp = client.get('/api/alias/etag_alias', headers=_api_headers())
    assert resp.status_code == 200
    etag = resp.headers.get('ETag')
    assert etag is not None

    resp2 = client.get('/api/alias/etag_alias', headers={
        **_api_headers(), 'If-None-Match': etag
    })
    assert resp2.status_code == 304


def test_api_cache_invalidation_on_create(client):
    register_user(client)
    login(client)

    resp = client.get('/api/alias/cache_test', headers=_api_headers())
    assert resp.status_code == 404

    create_alias(client, name="cache_test", command="echo cached")

    resp = client.get('/api/alias/cache_test', headers=_api_headers())
    assert resp.status_code == 200


def test_api_rate_limiting(client):
    register_user(client)
    login(client)

    for _ in range(61):
        resp = client.get('/api/alias/nonexistent', headers=_api_headers())
        if resp.status_code == 429:
            break
    else:
        assert False, "Rate limiter did not trigger after 61 requests"


def test_api_trim_endpoint(client):
    register_user(client)
    login(client)
    create_alias(client, name="trimmed", command="echo trimmed")

    resp = client.get('/api/alias/  trimmed  ', headers=_api_headers())
    assert resp.status_code == 200
    data = resp.get_json()
    assert data['alias'] == 'trimmed'
