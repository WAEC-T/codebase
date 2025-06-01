from app.extensions import db


class Latest(db.Model):
    __tablename__ = "latest"
    id = db.Column(db.Integer, primary_key=True)
    value = db.Column(db.Integer, nullable=False)
