import hashlib
from datetime import datetime
from flask import g

def get_user_id(username):
    """Convenience method to look up the id for a username."""
    rv = g.db.execute('select user_id from user where username = ?',
                      [username]).fetchone()
    return rv[0] if rv else None


def format_datetime(timestamp):
    """Format a timestamp for display."""
    print(timestamp)
    if timestamp is None:
        return "Unknown date"
    if isinstance(timestamp, str):
        dt = datetime.strptime(timestamp, '%Y-%m-%d %H:%M:%S')
        return dt.strftime('%Y-%m-%d @ %H:%M:%S')
    # If timestamp is an integer (Unix timestamp), convert it to datetime
    return datetime.utcfromtimestamp(int(timestamp)).strftime('%Y-%m-%d @ %H:%M')


def gravatar(email, size=80):
    """Return the gravatar image for the given email address."""
    return 'http://www.gravatar.com/avatar/%s?d=identicon&s=%d' % \
        (hashlib.md5(email.strip().lower().encode('utf-8')).hexdigest(), size)