from datetime import datetime
from flask import Blueprint, request, redirect, url_for, flash, session, abort
from app.models.message import Message
from app.main.routes import is_user_logged
from app.extensions import db

posts_bp = Blueprint("posts", __name__)


@posts_bp.route("/add_message", methods=["POST"])
def add_message():
    """Add a new message to the database."""
    is_user_logged()
    if "user_id" not in session:
        abort(401)

    message_text = request.form.get("text")

    if message_text:
        user_id = session["user_id"]
        new_message = Message(
            author_id=user_id, text=message_text, pub_date=datetime.now(), flagged=0
        )

        db.session.add(new_message)
        db.session.commit()
        flash("Your message was recorded")

    return redirect(request.referrer or url_for("main.timeline"))
