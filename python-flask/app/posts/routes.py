from flask import Blueprint, request, redirect, url_for, flash, session, abort, render_template
from app.models.post import Message
from datetime import datetime
from app.extensions import db

posts_bp = Blueprint('posts', __name__)

@posts_bp.route('/add_message', methods=['POST'])
def add_message():
    if 'user_id' not in session:
        abort(401)

    message_text = request.form.get('text')

    if message_text:
        user_id = session['user_id']
        new_message = Message(
            author_id=user_id,
            text=message_text,
            pub_date=datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
            flagged=0
        )

        db.session.add(new_message)
        db.session.commit()
        flash(f'Message was recorded')

    return redirect(url_for('main.timeline'))

