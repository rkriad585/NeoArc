# config.py

PORT = 59248
HOST = "0.0.0.0"
DEBUG = True

UPLOAD_FOLDER = 'static/uploads'
ALLOWED_EXTENSIONS = {'png', 'jpg', 'jpeg', 'gif', 'webp'}

DB_PATH = 'neoarc.db'

ADMIN_ROUTE = '/admin'

SECRET_KEY = 'neoarc_super_secret_key_change_me'

ADMIN_CREDENTIALS = {
    "email": "rkriad585@gamil.com",
    "password": "riad"
}