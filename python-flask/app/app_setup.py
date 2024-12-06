from dotenv import load_dotenv

from app.main.routes import main_bp
from app.posts.routes import posts_bp
from app.auth.routes import admin_bp
from app.api.routes import api_bp
from app.extensions import db
from app.utils import format_datetime, gravatar

# Load environment variables from .env
load_dotenv()


def register_blueprints(app):
    """Register all blueprints to the app."""
    app.register_blueprint(main_bp)
    app.register_blueprint(posts_bp)
    app.register_blueprint(admin_bp)
    app.register_blueprint(api_bp, url_prefix="/api")


def configure_server(app):
    """Configure server settings and database initialization."""
    app.config.from_pyfile("../config.py")
    app.debug = app.config["DEBUG"]
    db.init_app(app)
    # app.before_request(is_user_logged)


def set_environment_filters(app):
    """Set environment-specific Jinja filters."""
    app.jinja_env.filters["format_datetime"] = format_datetime
    app.jinja_env.filters["gravatar"] = gravatar


def prepare_application(app):
    """Prepare the application by configuring and setting up blueprints and filters."""
    configure_server(app)
    register_blueprints(app)
    set_environment_filters(app)
