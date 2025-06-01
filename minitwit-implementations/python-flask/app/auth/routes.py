from flask import Blueprint, redirect, url_for, session, flash, render_template, request
from app.extensions import db
from app.models.user import User

admin_bp = Blueprint("admin", __name__)


@admin_bp.route("/login", methods=["GET", "POST"])
def login():
    """Logs the user in."""
    if "user_id" in session:
        return redirect(url_for("main.timeline"))
    error = None
    if request.method == "POST":
        username = request.form.get("username")
        password = request.form.get("password")

        user = User.query.filter_by(username=username).first()

        if user is None:
            error = "Invalid username"
        elif user.pw_hash != password:
            error = "Invalid password"
        else:
            # user authenticated
            session["user_id"] = user.user_id
            flash("You were logged in")
            return redirect(url_for("main.timeline"))

    return render_template("login.html", error=error)


@admin_bp.route("/register", methods=["GET", "POST"])
def register():
    """Register user"""
    if "user_id" in session:
        return redirect(url_for("main.timeline"))
    error = None
    if request.method == "POST":
        username = request.form.get("username")
        email = request.form.get("email")
        password = request.form.get("password")
        password2 = request.form.get("password2")

        if not username:
            error = "You have to enter a username"
        elif not email or "@" not in email:
            error = "You have to enter a valid email address"
        elif not password:
            error = "You have to enter a password"
        elif password != password2:
            error = "The two passwords do not match"
        elif User.query.filter_by(username=username).first():
            error = "The username is already taken"
        else:
            new_user = User(
                username=username, email=email, pw_hash=password
            )
            db.session.add(new_user)
            db.session.commit()
            flash("You were successfully registered and can login now")
            return redirect(url_for("admin.login"))
    return render_template("register.html", error=error)


@admin_bp.route("/logout")
def logout():
    """logs user out"""
    session.pop("user_id", None)
    flash("You were logged out")
    return redirect(url_for("main.public_timeline"))
