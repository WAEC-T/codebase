from flask import Blueprint, render_template, redirect, url_for, session, abort, flash, g
from app.models.post import Message
from app.models.user import User, Follower
from app.extensions import db
from config import PER_PAGE

main_bp = Blueprint('main', __name__)

@main_bp.route('/')
def timeline():
    """Shows a user's timeline or the public timeline if no user is logged in."""
    if 'user_id' not in session:
        return redirect(url_for('main.public_timeline'))

    messages = db.session.query(Message, User).join(User, Message.author_id == User.user_id).order_by(
        Message.pub_date.desc()).limit(PER_PAGE).all()

    messages_with_users = [{'message': message, 'user': user} for message, user in messages]

    return render_template('timeline.html', messages=messages_with_users )


@main_bp.route('/public')
def public_timeline():
    """Displays the latest messages of all users."""
    """Displays the latest messages of all users."""
    # Join the Message and User tables to get the message and user details
    messages = db.session.query(Message, User).join(User, Message.author_id == User.user_id).order_by(
        Message.pub_date.desc()).limit(PER_PAGE).all()

    messages_with_users = [{'message': message, 'user': user} for message, user in messages]

    return render_template('timeline.html', messages=messages_with_users)

@main_bp.route('/<username>')
def user_timeline(username):
    """Displays a user's timeline."""
    profile_user = User.query.filter_by(username=username).first_or_404()

    followed = False
    if g.user:
        followed = Follower.query.filter_by(
            who_id=g.user.user_id, whom_id=profile_user.user_id
        ).first() is not None

    messages = db.session.query(Message, User).join(User, Message.author_id == User.user_id).order_by(
        Message.pub_date.desc()).limit(PER_PAGE).all()

    messages_with_users = [{'message': message, 'user': user} for message, user in messages]

    return render_template('timeline.html', messages=messages_with_users, followed=followed, profile_user=profile_user)



@main_bp.route('/<username>/follow')
def follow_user(username):
    """Adds the current user to the follower of the given user"""
    if 'user_id' not in session:
        abort(401)

    # Fetch the user to follow by their username
    whom = User.query.filter_by(username=username).first()

    if whom is None:
        abort(404)  # User to follow does not exist

    # Check if the user is already following this person
    existing_follower = Follower.query.filter_by(who_id=g.user.user_id,
                                                 whom_id=whom.user_id).first()  # execute follower query

    if existing_follower is None:
        # Create a new follower relationship if it doesn't exist
        new_follower = Follower(who_id=g.user.user_id, whom_id=whom.user_id)
        db.session.add(new_follower)
        db.session.commit()
        flash(f'You are now following "{username}"')
    else:
        flash(f'You are already following "{username}"')

    return redirect(url_for('user_timeline', username=username))


@main_bp.route('/<username>/unfollow')
def unfollow_user(username):
    """Removes the curent user as follower of the given user"""
    if 'user_id' not in session:
        abort(401)

    whom = User.query.filter_by(username=username).first()
    if whom is None:
        abort(404)

    # Check if the user is following this person
    existing_follower = Follower.query.filter_by(who_id=g.user.user_id,
                                                 whom_id=whom.user_id).first()  # execute follower query
    if existing_follower is None:
        flash(f'You are no longer following "{username}"')
    else:
        db.session.delete(existing_follower)
        db.session.commit()
    return redirect(url_for('user_timeline', username=username))


@main_bp.route('/cleardb')
def clean_up():
    """Clears the current database and reinitializes it."""
    print("Cleaning database from last run...")

    # Drop all tables in the database
    db.drop_all()

    print("Setting up new version of DB...")

    db.create_all()

    return redirect(url_for('public_timeline'))


@main_bp.route('/check_db')
def check_db_connection():
    """Checks if the connection to the database is established."""
    try:
        # Execute a simple query to check the connection
        result = db.session.execute('SELECT 1').scalar()
        if result == 1:
            return 'Database connection is successful!'
        else:
            return 'Database connection failed!', 500
    except Exception as e:
        # Catch any exceptions (like connection errors) and print them
        return f"Database connection failed! Error: {str(e)}", 500