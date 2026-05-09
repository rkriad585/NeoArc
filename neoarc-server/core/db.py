import config
import sqlite3
import time
import functools


def get_db():
    conn = sqlite3.connect(config.DB_PATH, check_same_thread=False)
    conn.row_factory = sqlite3.Row
    conn.execute("PRAGMA journal_mode=WAL")
    conn.execute("PRAGMA foreign_keys=ON")
    conn.execute("PRAGMA busy_timeout=5000")
    return conn


def retry_on_locked(f):
    @functools.wraps(f)
    def wrapper(*args, **kwargs):
        for attempt in range(3):
            try:
                return f(*args, **kwargs)
            except sqlite3.OperationalError as e:
                if "database is locked" in str(e) and attempt < 2:
                    time.sleep(0.1 * (attempt + 1))
                    continue
                raise
        return None
    return wrapper


def init_db():
    conn = get_db()

    conn.execute("""CREATE TABLE IF NOT EXISTS users (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        username TEXT UNIQUE,
                        email TEXT UNIQUE,
                        full_name TEXT,
                        password TEXT,
                        profile_pic TEXT,
                        is_blocked INTEGER DEFAULT 0)
                 """)
    conn.execute("""CREATE TABLE IF NOT EXISTS aliases (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        alias_name TEXT UNIQUE,
                        command TEXT,
                        exec_type TEXT,
                        user_id INTEGER,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
                 """)

    conn.execute("CREATE INDEX IF NOT EXISTS idx_aliases_alias_name ON aliases(alias_name)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_aliases_user_id ON aliases(user_id)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)")
    conn.execute("""CREATE TABLE IF NOT EXISTS reset_tokens (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        user_id INTEGER NOT NULL,
                        token_hash TEXT NOT NULL,
                        expires_at TIMESTAMP NOT NULL,
                        used INTEGER DEFAULT 0,
                        FOREIGN KEY (user_id) REFERENCES users(id))
                 """)

    conn.execute("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)")

    columns = [row[1] for row in conn.execute("PRAGMA table_info(aliases)").fetchall()]
    if 'created_at' not in columns:
        conn.execute("ALTER TABLE aliases ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
        logger = __import__('logging').getLogger('neoarc')
        logger.info("Added created_at column to aliases table")

    conn.commit()
    conn.close()
