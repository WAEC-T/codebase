# -*- coding: utf-8 -*-
"""
    MiniTwit
    ~~~~~~~~

    A microblogging application written with Flask and sqlite3.

    :copyright: (c) 2010 by Armin Ronacher.
    :license: BSD, see LICENSE for more details.
"""
from __future__ import with_statement
import os
import time
import logging
from hashlib import md5
from datetime import datetime
from flask import Flask, request, session, url_for, redirect, \
    render_template, abort, g, flash
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import text
from flask_migrate import Migrate
from werkzeug.security import check_password_hash, generate_password_hash
from dotenv import load_dotenv

load_dotenv()

# configuration
username = os.getenv('POSTGRES_USER', 'user')
password = os.getenv('POSTGRES_PASSWORD', 'pass')
host = os.getenv('POSTGRES_HOST', '192.168.8.175')
port = os.getenv('POSTGRES_PORT', '5432')
database = os.getenv('POSTGRES_DB', 'mydatabase')
SECRET_KEY = os.getenv('SECRET_KEY')

DATABASE_URL = f"postgresql://{username}:{password}@{host}:{port}/{database}"
PER_PAGE = 30
DEBUG = True

# create our little application :)
app = Flask(__name__)
app.config[
    'SQLALCHEMY_DATABASE_URI'] = f"postgresql://{os.getenv('POSTGRES_USER', username)}:{os.getenv('POSTGRES_PASSWORD', password)}@{os.getenv('POSTGRES_HOST', host)}:{os.getenv('POSTGRES_PORT', port)}/{os.getenv('POSTGRES_DB', database)}"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
app.secret_key = SECRET_KEY
app.debug = True

# Initialize database
db = SQLAlchemy(app)
migrate = Migrate(app, db)


# Define your models (this replaces schema.sql)
class User(db.Model):
    __tablename__ = 'users'
    user_id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.Text, nullable=False, unique=True)
    email = db.Column(db.Text, nullable=False, unique=True)
    pw_hash = db.Column(db.Text, nullable=False)


class Follower(db.Model):
    __tablename__ = 'followers'
    who_id = db.Column(db.Integer, db.ForeignKey('users.user_id'), primary_key=True)
    whom_id = db.Column(db.Integer, db.ForeignKey('users.user_id'), primary_key=True)


class Message(db.Model):
    __tablename__ = 'messages'
    message_id = db.Column(db.Integer, primary_key=True)
    author_id = db.Column(db.Integer, db.ForeignKey('users.user_id'), nullable=False)
    text = db.Column(db.Text, nullable=False)
    pub_date = db.Column(db.Text)
    flagged = db.Column(db.Integer)


class Latest(db.Model):
    __tablename__ = 'latest'
    id = db.Column(db.Integer, primary_key=True)
    value = db.Column(db.Integer, nullable=False)


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


def gravatar_url(email, size=80):
    """Return the gravatar image for the given email address."""
    return 'http://www.gravatar.com/avatar/%s?d=identicon&s=%d' % \
        (md5(email.strip().lower().encode('utf-8')).hexdigest(), size)

@app.before_request
def before_request():
    g.user = None
    if 'user_id' in session:
        g.user = User.query.filter_by(user_id=session['user_id']).first() # g = global object

@app.route('/')
def timeline():
    print(request.endpoint)
    """Shows a user's timeline or the public timeline if no user is logged in."""
    if 'user_id' not in session:
        return redirect(url_for('public_timeline'))

    messages = db.session.query(Message, User).join(User, Message.author_id == User.user_id).order_by(
        Message.pub_date.desc()).limit(PER_PAGE).all()

    messages_with_users = [{'message': message, 'user': user} for message, user in messages]

    return render_template('timeline.html', messages=messages_with_users )


@app.route('/public')
def public_timeline():
    """Displays the latest messages of all users."""
    """Displays the latest messages of all users."""
    # Join the Message and User tables to get the message and user details
    messages = db.session.query(Message, User).join(User, Message.author_id == User.user_id).order_by(
        Message.pub_date.desc()).limit(PER_PAGE).all()

    messages_with_users = [{'message': message, 'user': user} for message, user in messages]

    return render_template('timeline.html', messages=messages_with_users)


@app.route('/<username>')
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


@app.route('/<username>/follow')
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


@app.route('/<username>/unfollow')
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


@app.route('/add_message', methods=['POST'])
def add_message():
    """Adds a new message to the timeline"""
    if 'user_id' not in session:
        abort(401)

    message_text = request.form.get('text')

    if message_text:
        user_id = session['user_id']

        # new message
        new_message = Message(
            author_id=user_id,
            text=message_text,
            pub_date=datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
            flagged=0
        )

        db.session.add(new_message)
        db.session.commit()
        flash(f'Message was recorded')

    return redirect(url_for('timeline'))


@app.route('/login', methods=['GET', 'POST'])
def login():
    """Logs the user in."""
    if 'user_id' in session:
        return redirect(url_for('timeline'))
    error = None
    if request.method == 'POST':
        username = request.form.get('username')
        password = request.form.get('password')

        user = User.query.filter_by(username=username).first()

        if user is None:
            error = 'Invalid username'
        elif not check_password_hash(user.pw_hash, password):
            error = 'Invalid password'
        else:
            # user authenticated
            session['user_id'] = user.user_id
            print(session['user_id'])
            flash('You were logged in')
            return redirect(url_for('timeline'))

    return render_template('login.html', error=error)

@app.route('/register', methods=['GET', 'POST'])
def register():
    """Register user"""
    if 'user_id' in session:
        return redirect(url_for('timeline'))
    error = None
    if request.method == 'POST':
        username = request.form.get('username')
        email = request.form.get('email')
        password = request.form.get('password')
        password2 = request.form.get('password2')

        if not username:
            error = 'You have to enter a username'
        elif not email or '@' not in email:
            error = 'You have to enter a valid email address'
        elif not password:
            error = 'You have to enter a password'
        elif password != password2:
            error = 'The two passwords do not match'
        elif User.query.filter_by(username=username).first():
            error = 'The username is already taken'
        else:
            new_user = User(
                username=username,
                email=email,
                pw_hash=generate_password_hash(password)
            )
            db.session.add(new_user)
            db.session.commit()
            flash('You were successfully registered and can login now')
            return redirect(url_for('login'))
    return render_template('register.html', error=error)


@app.route('/logout')
def logout():
    """logs user out"""
    session.pop('user_id', None)
    flash('You were logged out')
    return redirect(url_for('public_timeline'))


@app.route('/cleardb')
def clean_up():
    """Clears the current database and reinitializes it."""
    print("Cleaning database from last run...")

    # Drop all tables in the database
    db.drop_all()

    print("Setting up new version of DB...")

    db.create_all()

    return redirect(url_for('public_timeline'))


@app.route('/check_db')
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



# Register filters in Jinja
app.jinja_env.filters['format_datetime'] = format_datetime
app.jinja_env.filters['gravatar'] = gravatar_url

app.secret_key = SECRET_KEY
app.debug = DEBUG

if __name__ == '__main__':
    with app.app_context():
        try:
            # Attempt to connect to the database
            db.session.execute(text('SELECT 1'))
            logging.info("Database connection successful!")
        except Exception as e:
            logging.error(f"Database connection failed! Error: {str(e)}")
    app.run(host="0.0.0.0")
