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
import sqlite3
from hashlib import md5
from datetime import datetime
from contextlib import closing
from flask import Flask, request, session, url_for, redirect, \
    render_template, abort, g, flash
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from werkzeug.security import check_password_hash, generate_password_hash
from pathlib import Path
from dotenv import load_dotenv

load_dotenv()

# configuration
username = os.getenv('POSTGRES_USER', 'user')
password = os.getenv('POSTGRES_PASSWORD', 'pass')
host = os.getenv('POSTGRES_HOST', '192.168.8.175')
port = os.getenv('POSTGRES_PORT', '5432')
database = os.getenv('POSTGRES_DB', 'mydatabase')

DATABASE_URL = f"postgresql://{username}:{password}@{host}:{port}/{database}"
PER_PAGE = 30
DEBUG = True
SECRET_KEY = os.getenv('SECRET_KEY')

# create our little application :)
app = Flask(__name__)
app.config[
    'SQLALCHEMY_DATABASE_URI'] = f"postgresql://{os.getenv('POSTGRES_USER', 'user')}:{os.getenv('POSTGRES_PASSWORD', 'pass')}@{os.getenv('POSTGRES_HOST', '192.168.8.175')}:{os.getenv('POSTGRES_PORT', '5432')}/{os.getenv('POSTGRES_DB', 'mydatabase')}"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
app.secret_key = os.getenv('SECRET_KEY', 'development key')
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
    return datetime.utcfromtimestamp(timestamp).strftime('%Y-%m-%d @ %H:%M')


def gravatar_url(email, size=80):
    """Return the gravatar image for the given email address."""
    return 'http://www.gravatar.com/avatar/%s?d=identicon&s=%d' % \
        (md5(email.strip().lower().encode('utf-8')).hexdigest(), size)


@app.before_request
def before_request():
    """Get the current user from the session."""
    g.user = None
    if 'user_id' in session:
        g.user = User.query.filter_by(user_id=session['user_id']).first()


@app.after_request
def after_request(response):
    """Closes the database again at the end of the request."""
    g.db.close()
    return response


@app.route('/')
def timeline():
    """Shows a user's timeline or the public timeline if no user is logged in."""
    if not g.user:
        return redirect(url_for('public_timeline'))

    messages = Message.query.filter(
        (Message.author_id == g.user.user_id) |
        (Message.author_id.in_(
            Follower.query.with_entities(Follower.whom_id).filter_by(who_id=g.user.user_id)
        ))
    ).order_by(Message.pub_date.desc()).limit(PER_PAGE).all()

    return render_template('timeline.html', messages=messages)


@app.route('/public')
def public_timeline():
    """Displays the latest messages of all users."""
    messages = Message.query.order_by(Message.pub_date.desc()).limit(PER_PAGE).all()
    return render_template('timeline.html', messages=messages)


@app.route('/<username>')
def user_timeline(username):
    """Displays a user's timeline."""
    profile_user = User.query.filter_by(username=username).first_or_404()

    followed = False
    if g.user:
        followed = Follower.query.filter_by(
            who_id=g.user.user_id, whom_id=profile_user.user_id
        ).first() is not None

    messages = Message.query.filter_by(author_id=profile_user.user_id).order_by(Message.pub_date.desc()).limit(
        PER_PAGE).all()

    return render_template('timeline.html', messages=messages, followed=followed, profile_user=profile_user)


@app.route('/<username>/follow')
def follow_user(username):
    """Adds the current user to the follower of the given user"""
    if not g.user:
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
    if not g.user:
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


@app.route('/<username>/add_message', methods=['POST'])
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
            pub_date=str(int(time.time())),
            flagged=0
        )

        db.session.add(new_message)
        db.session.commit()
        flash(f'Message was recorded')

    return redirect(url_for('timeline'))


@app.route('/login', methods=['GET', 'POST'])
def login():
    """Logs the user in."""
    if g.user:
        return redirect(url_for('timeline'))
    error = None
    if request.method == 'POST':
        user = query_db('''select * from user where
            username = ?''', [request.form['username']], one=True)
        if user is None:
            error = 'Invalid username'
        elif not check_password_hash(user['pw_hash'],
                                     request.form['password']):
            error = 'Invalid password'
        else:
            flash('You were logged in')
            session['user_id'] = user['user_id']
            return redirect(url_for('timeline'))
    return render_template('login.html', error=error)


@app.route('/register', methods=['GET', 'POST'])
def register():
    """Registers the user."""
    if g.user:
        return redirect(url_for('timeline'))

    error = None
    if request.method == 'POST':
        if not request.form['username']:
            error = 'You have to enter a username'
        elif not request.form['email'] or '@' not in request.form['email']:
            error = 'You have to enter a valid email address'
        elif not request.form['password']:
            error = 'You have to enter a password'
        elif request.form['password'] != request.form['password2']:
            error = 'The two passwords do not match'
        elif User.query.filter_by(username=request.form['username']).first():
            error = 'The username is already taken'
        else:
            user = User(
                username=request.form['username'],
                email=request.form['email'],
                pw_hash=generate_password_hash(request.form['password'])
            )
            db.session.add(user)
            db.session.commit()
            flash('You were successfully registered and can login now')
            return redirect(url_for('login'))

    return render_template('register.html', error=error)


@app.route('/logout')
def logout():
    """Logs the user out"""
    flash('You were logged out')
    session.pop('user_id', None)
    return redirect(url_for('public_timeline'))


@app.route('/cleardb')
def clean_up():
    print("Cleaning database from last run...")
    Path(DATABASE).unlink()

    def fill_db():
        with closing(connect_db()) as db:
            with open('dump.sql') as f:
                db.cursor().executescript(f.read())
            db.commit()

    print("Setting up new version of DB...")
    init_db()
    fill_db()
    return redirect(url_for('public_timeline'))


# add some filters to jinja and set the secret key and debug mode
# from the configuration.
app.jinja_env.filters['datetimeformat'] = format_datetime
app.jinja_env.filters['gravatar'] = gravatar_url
app.secret_key = SECRET_KEY
app.debug = DEBUG

if __name__ == '__main__':
    app.run(host="0.0.0.0")
