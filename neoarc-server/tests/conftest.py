import os
import sys
import tempfile
import shutil
import pytest

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

# Use a session-scoped temp DB (created once, cleaned before each test)
_tmp_db = tempfile.mktemp(suffix='.db')
_tmp_upload = tempfile.mkdtemp()


def _init_env():
    os.environ['NEOARC_DB_PATH'] = _tmp_db
    os.environ['NEOARC_SECRET_KEY'] = 'test-secret-key-32chars!!ok'
    os.environ['NEOARC_API_TOKEN'] = 'test-api-token-32chars!!'
    os.environ['NEOARC_ADMIN_PASSWORD'] = 'adminpass123'
    os.environ['NEOARC_UPLOAD_FOLDER'] = _tmp_upload
    os.environ['NEOARC_DEBUG'] = 'False'


_init_env()

import config
assert config.DB_PATH == _tmp_db

from core.db import init_db
init_db()

# Clean slate after import
from core.db import get_db
conn = get_db()
for table in ('reset_tokens', 'aliases', 'users'):
    conn.execute(f"DELETE FROM {table}")
conn.commit()
conn.close()

from main import app


@pytest.fixture
def client():
    app.config['TESTING'] = True
    from core.db import get_db
    conn = get_db()
    for table in ('reset_tokens', 'aliases', 'users'):
        conn.execute(f"DELETE FROM {table}")
    conn.execute("DELETE FROM sqlite_sequence")
    conn.commit()
    conn.close()
    from core.rate_limiter import rate_limiter
    rate_limiter.clear()
    with app.test_client() as c:
        yield c


def pytest_sessionfinish():
    try:
        os.remove(_tmp_db)
        for f in os.listdir('.'):
            if f.startswith(os.path.basename(_tmp_db)):
                os.remove(f)
        shutil.rmtree(_tmp_upload, ignore_errors=True)
    except OSError:
        pass
