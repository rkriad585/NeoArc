# db.py

import config
import sqlite3


def get_db():
    conn = sqlite3.connect(config.DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn

def init_db():
    conn = get_db()
    conn.execute('''CREATE TABLE IF NOT EXISTS users (
                        id INTEGER PRIMARY KEY AUTOINCREMENT, 
                        username TEXT UNIQUE, 
                        email TEXT UNIQUE,
                        full_name TEXT,
                        password TEXT,
                        profile_pic TEXT,
                        is_blocked INTEGER DEFAULT 0)
                 ''')
    conn.execute('''CREATE TABLE IF NOT EXISTS aliases (
                        id INTEGER PRIMARY KEY AUTOINCREMENT, 
                        alias_name TEXT UNIQUE, 
                        command TEXT, 
                        exec_type TEXT, 
                        user_id INTEGER)''')
    conn.commit()
    conn.close()
