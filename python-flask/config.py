import os
# test push cd in python-flask 3
SQLALCHEMY_DATABASE_URI = os.getenv('DATABASE_URL')
SECRET_KEY = os.getenv('SECRET_KEY')
SQLALCHEMY_TRACK_MODIFICATIONS = False
PER_PAGE = 30