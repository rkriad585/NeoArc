import sys
sys.path.insert(0, '.')

from main import app
import config
from core.helpers import db_execute

# Check users in DB using the actual connection
user = db_execute('SELECT id, username, email FROM users WHERE email=?', ('mdriyadkhan585@gmail.com',)).fetchone()
print(f'Direct query result: {user}')
if user:
    print(f'  id={user["id"]} username={user["username"]} email={user["email"]}')
else:
    # Try all users
    all_users = db_execute('SELECT id, username, email FROM users').fetchall()
    print(f'All users ({len(all_users)}):')
    for u in all_users:
        print(f'  id={u["id"]} username={u["username"]} email=[{u["email"]}]')

# Now test with the app context
with app.test_request_context('/'):
    user2 = db_execute('SELECT id, username, email FROM users WHERE email=?', ('mdriyadkhan585@gmail.com',)).fetchone()
    print(f'In-request query: {user2}')
