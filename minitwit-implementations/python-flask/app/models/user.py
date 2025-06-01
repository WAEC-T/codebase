from app.extensions import db


class User(db.Model):
    __tablename__ = "users"
    user_id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.Text, nullable=False, unique=True)
    email = db.Column(db.Text, nullable=False, unique=True)
    pw_hash = db.Column(db.Text, nullable=False)


class Follower(db.Model):
    __tablename__ = "followers"
    who_id = db.Column(db.Integer, db.ForeignKey("users.user_id"), primary_key=True)
    whom_id = db.Column(db.Integer, db.ForeignKey("users.user_id"), primary_key=True)
