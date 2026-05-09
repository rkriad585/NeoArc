import re

import config


_EMAIL_RE = re.compile(r'^[^@\s]+@[^@\s]+\.[^@\s]+$')
_ALIAS_NAME_RE = re.compile(r'^[a-zA-Z0-9][a-zA-Z0-9_-]{0,99}$')
_FULL_NAME_RE = re.compile(r'^.{1,200}$')
_PASSWORD_MIN = 8

_IMAGE_SIGNATURES = {
    b'\x89PNG\r\n\x1a\n': 'png',
    b'\xff\xd8\xff': 'jpg',
    b'GIF8': 'gif',
    b'RIFF': 'webp',
}


def validate_email(email: str) -> bool:
    return bool(email and _EMAIL_RE.match(email))


def validate_alias_name(name: str) -> bool:
    return bool(name and _ALIAS_NAME_RE.match(name))


def validate_full_name(name: str) -> bool:
    return bool(name and _FULL_NAME_RE.match(name))


def validate_password(password: str, username: str = '', email: str = '') -> list[str]:
    errors = []
    if len(password) < _PASSWORD_MIN:
        errors.append(f'Password must be at least {_PASSWORD_MIN} characters.')
    if password and (password == username or password == email):
        errors.append('Password must not match your username or email.')
    return errors


def _detect_image_mime(data: bytes) -> str | None:
    for signature, fmt in _IMAGE_SIGNATURES.items():
        if data.startswith(signature):
            return fmt
    return None


def allowed_file(filename):
    if '.' not in filename:
        return False
    ext = filename.rsplit('.', 1)[1].lower()
    return ext in config.ALLOWED_EXTENSIONS


def validate_image(file) -> str | None:
    if not file or not file.filename:
        return None
    if not allowed_file(file.filename):
        return 'Invalid image format. Allowed: png, jpg, jpeg, gif, webp.'
    file.seek(0)
    magic = file.read(8)
    file.seek(0)
    detected = _detect_image_mime(magic)
    if not detected:
        return 'File content does not match an allowed image type.'
    return None
