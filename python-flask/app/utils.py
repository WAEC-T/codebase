from datetime import datetime
import hashlib
from flask import jsonify
from app.models.user import User
from app.models.latest import Latest
from app.extensions import db


def get_user_id(username):
    """Convenience method to look up the id for a username."""
    user = User.query.filter_by(username=username).first()
    return user.user_id if user else None


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


def not_req_from_simulator(request):
    from_simulator = request.headers.get("Authorization")
    if from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh":
        error = "You are not authorized to use this resource!"
        return jsonify({"status": 403, "error_msg": error}), 403


def update_latest(request):
    """Update the latest processed command id"""
    parsed_command_id = request.args.get("latest", type=int)
    if parsed_command_id is not None:
        latest_entry = Latest.query.first()
        if latest_entry:
            latest_entry.value = parsed_command_id
        else:
            latest_entry = Latest(id=1, value=parsed_command_id)
            db.session.add(latest_entry)
        db.session.commit()  # Save changes