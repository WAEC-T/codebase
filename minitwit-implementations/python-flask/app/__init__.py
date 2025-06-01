from flask import Flask
from app.app_setup import prepare_application


def create_app():
    """Flask application factory."""
    app = Flask(__name__)
    prepare_application(app)
    return app
