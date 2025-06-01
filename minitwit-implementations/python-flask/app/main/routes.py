from flask import (
    Blueprint,
    render_template,
    redirect,
    url_for,
    session,
    abort,
    flash,
    g,
)
from sqlalchemy.exc import SQLAlchemyError
from app.models.message import Message
from app.models.user import User, Follower
from app.extensions import db
from config import PER_PAGE

main_bp = Blueprint("main", __name__)


def is_user_logged():
    """Checks if a user is logged in before each request."""
    g.user = None
    if "user_id" in session:
        g.user = User.query.filter_by(user_id=session["user_id"]).first()


@main_bp.route("/")
def timeline():
    """Shows a users timeline or if no user is logged in it will
    redirect to the public timeline.  This timeline shows the user's
    messages as well as all the messages of followed users."""
    is_user_logged()
    if "user_id" not in session:
        return redirect(url_for("main.public_timeline"))
    user_id = session["user_id"]
    messages = (
        db.session.query(Message, User)
        .join(User, Message.author_id == User.user_id)
        .filter(
            (Message.author_id == user_id)
            | (
                Message.author_id.in_(
                    db.session.query(Follower.whom_id).filter(
                        Follower.who_id == user_id
                    )
                )
            )
        )
        .order_by(Message.pub_date.desc())
        .limit(PER_PAGE)
        .all()
    )

    messages_with_users = [
        {"message": message, "user": user} for message, user in messages
    ]

    return render_template("timeline.html", messages=messages_with_users)


@main_bp.route("/public")
def public_timeline():
    """Displays the latest messages of all users."""
    is_user_logged()
    messages = (
        db.session.query(Message, User)
        .join(User, Message.author_id == User.user_id)
        .order_by(Message.pub_date.desc())
        .limit(PER_PAGE)
        .all()
    )

    messages_with_users = [
        {"message": message, "user": user} for message, user in messages
    ]

    return render_template("timeline.html", messages=messages_with_users)


@main_bp.route("/user/<username>")
def user_timeline(username):
    """Displays a user's timeline."""
    is_user_logged()
    profile_user = User.query.filter_by(username=username).first_or_404()

    followed = False
    if g.user:
        followed = (
            Follower.query.filter_by(
                who_id=g.user.user_id, whom_id=profile_user.user_id
            ).first()
            is not None
        )

    messages = (
        Message.query.filter_by(author_id=profile_user.user_id, flagged=0)
        .order_by(Message.pub_date.desc())
        .limit(PER_PAGE)
        .all()
    )

    messages_with_users = [
        {"message": message, "user": profile_user} for message in messages
    ]

    return render_template(
        "timeline.html",
        messages=messages_with_users,
        followed=followed,
        profile_user=profile_user,
    )


@main_bp.route("/<username>/follow")
def follow_user(username):
    """Adds the current user to the follower of the given user"""
    is_user_logged()
    if "user_id" not in session:
        abort(401)

    whom = User.query.filter_by(username=username).first()

    if whom is None:
        abort(404)

    new_follower = Follower(who_id=g.user.user_id, whom_id=whom.user_id)
    db.session.add(new_follower)
    db.session.commit()
    flash(f"You are now following {username}")

    return redirect(url_for("main.user_timeline", username=username))


@main_bp.route("/<username>/unfollow")
def unfollow_user(username):
    """Removes the curent user as follower of the given user"""
    is_user_logged()
    if "user_id" not in session:
        abort(401)

    whom = User.query.filter_by(username=username).first()
    if whom is None:
        abort(404)

    existing_follower = Follower.query.filter_by(
        who_id=g.user.user_id, whom_id=whom.user_id
    ).first()
    if existing_follower is None:
        abort(404)
    else:
        db.session.delete(existing_follower)
        db.session.commit()
        flash(f"You are no longer following {username}")
    return redirect(url_for("main.user_timeline", username=username))


@main_bp.route("/cleardb")
def clean_up():
    """Clears the current database and reinitializes it."""
    db.drop_all()
    db.create_all()

    return redirect(url_for("public_timeline"))


@main_bp.route("/check_db")
def check_db_connection():
    """Checks if the connection to the database is established."""
    try:
        # Execute a simple query to check the connection
        result = db.session.execute("SELECT 1").scalar()
        if result == 1:
            return "Database connection is successful!"
        return "Database connection failed!", 500
    except SQLAlchemyError as e:
        # Catch only database-related errors
        return f"Database connection failed! Error: {str(e)}", 500
