import logging
import smtplib
import sqlite3
from email.mime.text import MIMEText

from flask import g

from core.db import get_db

logger = logging.getLogger("neoarc")

_RESET_HTML_TEMPLATE = """\
<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<body style="margin:0;padding:0;background:#000;font-family:Inter,system-ui,sans-serif">
<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:#000;padding:40px 20px">
<tr><td align="center">
<table role="presentation" width="480" cellpadding="0" cellspacing="0" style="background:rgba(255,255,255,0.03);border:1px solid rgba(255,255,255,0.08);border-radius:24px;padding:40px">
<tr><td style="text-align:center;padding-bottom:8px">
<span style="font-family:DotGothic16,sans-serif;font-size:28px;color:#fff;letter-spacing:2px">NEO<span style="color:#ea2b2b">ARC</span></span>
</td></tr>
<tr><td style="text-align:center;padding-bottom:24px">
<span style="color:#9ca3af;font-size:12px;font-family:monospace">PASSWORD RESET NODE</span>
</td></tr>
<tr><td style="height:1px;background:linear-gradient(to right,transparent,rgba(255,255,255,0.08),transparent);"></td></tr>
<tr><td style="padding:24px 0;color:#d1d5db;font-size:14px;line-height:1.6">
Hello <strong style="color:#fff">%s</strong>,
<br><br>
A password reset request was received for your NeoArc account.
<br><br>
Use the gateway below to reset your credentials. This link expires in <strong style="color:#ea2b2b">1 hour</strong>.
</td></tr>
<tr><td style="padding:12px 0" align="center">
<a href="%s" style="display:inline-block;background:#fff;color:#000;font-family:DotGothic16,sans-serif;font-size:16px;padding:14px 40px;border-radius:999px;text-decoration:none;letter-spacing:1px">RESET PASSWORD</a>
</td></tr>
<tr><td style="padding:24px 0 0 0;color:#6b7280;font-size:11px;font-family:monospace;word-break:break-all">
If the button does not work, copy and paste this link into your browser:<br>
<span style="color:#ea2b2b">%s</span>
</td></tr>
<tr><td style="padding:24px 0 0 0;color:#6b7280;font-size:11px;font-family:monospace;border-top:1px dashed rgba(255,255,255,0.08)">
If you did not request this reset, you can safely ignore this email.
</td></tr>
</table>
</td></tr>
</table>
</body>
</html>"""


def get_conn():
    if '_db' not in g:
        g._db = get_db()
    return g._db


def db_execute(query, params=()):
    try:
        return get_conn().execute(query, params)
    except sqlite3.OperationalError as e:
        logger.error("Database error: %s\nQuery: %s\nParams: %s", e, query, params)
        raise


def db_commit():
    try:
        get_conn().commit()
    except sqlite3.OperationalError as e:
        logger.error("Database commit error: %s", e)
        raise


def send_password_reset_email(to_email: str, to_name: str, reset_link: str) -> bool:
    import config
    if not config.SMTP_ENABLED:
        logger.info("SMTP not configured. Password reset link for %s: %s", to_email, reset_link)
        return False

    html = _RESET_HTML_TEMPLATE % (to_name, reset_link, reset_link)
    msg = MIMEText(html, "html")
    msg["Subject"] = "NeoArc — Password Reset"
    msg["From"] = config.SMTP_FROM
    msg["To"] = to_email

    try:
        server = smtplib.SMTP(config.SMTP_HOST, config.SMTP_PORT, timeout=15)
        server.starttls()
        server.login(config.SMTP_USER, config.SMTP_PASSWORD)
        server.send_message(msg)
        server.quit()
        logger.info("Password reset email sent to %s", to_email)
        return True
    except smtplib.SMTPAuthenticationError:
        logger.error("SMTP authentication failed for %s", config.SMTP_USER)
        return False
    except smtplib.SMTPException as e:
        logger.error("SMTP error sending to %s: %s", to_email, e)
        return False
    except OSError as e:
        logger.error("SMTP connection failed for %s: %s", to_email, e)
        return False
