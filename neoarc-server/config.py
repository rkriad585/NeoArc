import os
import secrets
import warnings
from dotenv import load_dotenv
from werkzeug.security import generate_password_hash

_BASE_DIR = os.path.dirname(os.path.abspath(__file__))

load_dotenv()

PORT = int(os.environ.get("NEOARC_PORT", "59248"))
HOST = os.environ.get("NEOARC_HOST", "0.0.0.0")
DEBUG = os.environ.get("NEOARC_DEBUG", "False").lower() == "true"

UPLOAD_FOLDER = os.path.join(_BASE_DIR, os.environ.get("NEOARC_UPLOAD_FOLDER", "static/uploads"))
ALLOWED_EXTENSIONS = {"png", "jpg", "jpeg", "gif", "webp"}
MAX_CONTENT_LENGTH = int(os.environ.get("NEOARC_MAX_UPLOAD_MB", "5")) * 1024 * 1024

DB_PATH = os.environ.get("NEOARC_DB_PATH", "neoarc.db")
if not os.path.isabs(DB_PATH):
    DB_PATH = os.path.join(_BASE_DIR, DB_PATH)

ADMIN_ROUTE = os.environ.get("NEOARC_ADMIN_ROUTE", "/admin")

SECRET_KEY = os.environ.get("NEOARC_SECRET_KEY")
if SECRET_KEY is None:
    SECRET_KEY = secrets.token_hex(32)
    warnings.warn("Auto-generated SECRET_KEY. Sessions invalidated on restart. Set NEOARC_SECRET_KEY env var.")

ADMIN_EMAIL = os.environ.get("NEOARC_ADMIN_EMAIL", "admin@neoarc.local")
ADMIN_PASSWORD_HASH = os.environ.get("NEOARC_ADMIN_PASSWORD_HASH")
if ADMIN_PASSWORD_HASH is None:
    ADMIN_PASSWORD_HASH = generate_password_hash(
        os.environ.get("NEOARC_ADMIN_PASSWORD", "change_me")
    )
    warnings.warn("Using default admin password. Set NEOARC_ADMIN_PASSWORD env var.")

API_TOKEN = os.environ.get("NEOARC_API_TOKEN", "").strip()
if not API_TOKEN:
    API_TOKEN = secrets.token_hex(32)
    warnings.warn("Auto-generated API token for CLI. Set NEOARC_API_TOKEN env var for persistence.")

SMTP_HOST = os.environ.get("NEOARC_SMTP_HOST", "smtp.gmail.com")
SMTP_PORT = int(os.environ.get("NEOARC_SMTP_PORT", "587"))
SMTP_USER = os.environ.get("NEOARC_SMTP_USER", "").strip()
SMTP_PASSWORD = os.environ.get("NEOARC_SMTP_PASSWORD", "").strip()
SMTP_FROM = os.environ.get("NEOARC_SMTP_FROM", SMTP_USER)
SMTP_ENABLED = bool(SMTP_USER and SMTP_PASSWORD)
