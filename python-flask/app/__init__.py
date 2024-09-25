from distutils.command.config import config

from flask import Flask, session, g
from app.extensions import db, migrate
from app.main.routes import main_bp
from app.posts.routes import posts_bp
from app.auth.routes import admin_bp
from app.simulator.routes import sim_bp
from app.utils import format_datetime, gravatar
from app.models.user import User

def before_request():
    """Runs before every request to check if a user is logged in."""
    g.user = None
    if 'user_id' in session:
        g.user = User.query.filter_by(user_id=session['user_id']).first()

def create_app():
    app = Flask(__name__)

    # Configurations
    app.config.from_pyfile('../config.py')

    # Initialize extensions
    db.init_app(app)
    migrate.init_app(app, db)

    # Register Blueprints
    app.register_blueprint(main_bp)
    app.register_blueprint(posts_bp)
    app.register_blueprint(admin_bp)
    app.register_blueprint(sim_bp)

    app.jinja_env.filters['format_datetime'] = format_datetime
    app.jinja_env.filters['gravatar'] = gravatar

    app.secret_key = app.config['SECRET_KEY']
    app.debug = app.config['DEBUG']

    app.before_request(before_request)

    return app
