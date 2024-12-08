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
    if timestamp is None:
        return "Unknown date"
    if isinstance(timestamp, datetime):
        return timestamp.strftime("%Y-%m-%d @ %H:%M")
    return str(timestamp)


def gravatar(email, size=80):
    """Return the gravatar image for the given email address."""
    email_hash = hashlib.md5(email.strip().lower().encode("utf-8")).hexdigest()
    return f"http://www.gravatar.com/avatar/{email_hash}?d=identicon&s={size}"


def not_req_from_simulator(request):
    """Verifies if the request is authorized by checking the Authorization header."""
    from_simulator = request.headers.get("Authorization")
    if from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh":
        error = "You are not authorized to use this resource!"
        return jsonify({"status": 403, "error_msg": error}), 403

    return None


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
