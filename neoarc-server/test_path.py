import sys
sys.path.insert(0, '.')
import config
print(f'Config DB_PATH: {config.DB_PATH}')

import sqlite3
conn = sqlite3.connect(config.DB_PATH)
users = conn.execute('SELECT id, username, email FROM users').fetchall()
print(f'Direct SQLite3 users: {len(users)}')
for u in users:
    print(f'  id={u[0]} username={u[1]} email=[{u[2]}]')

# Now check what the app sees
from main import app
with app.test_request_context('/'):
    from core.helpers import db_execute
    user = db_execute('SELECT id, username, email FROM users WHERE email=?', ('mdriyadkhan585@gmail.com',)).fetchone()
    print(f'db_execute found user: {user is not None}')
    if user:
        print(f'  id={user["id"]} username={user["username"]}')
    else:
        all_u = db_execute('SELECT id, username, email FROM users').fetchall()
        print(f'db_execute all users: {len(all_u)}')
        for u in all_u:
            print(f'  id={u["id"]} username={u["username"]} email=[{u["email"]}]')

# Check the DB file from the app's perspective
with app.test_request_context('/'):
    from core.db import get_db
    conn2 = get_db()
    print(f'\nget_db() path derived from config')
    users2 = conn2.execute('SELECT id, username, email FROM users').fetchall()
    print(f'get_db users: {len(users2)}')
