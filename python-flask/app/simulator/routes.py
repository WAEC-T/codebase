from flask import Blueprint
from flask import jsonify, request, g, abort
from datetime import datetime
from app.models.user import User, Follower
from app.models.post import Message
from app.models.latest import Latest
from werkzeug.security import generate_password_hash
from app.utils import not_req_from_simulator, update_latest, get_user_id
from app.extensions import db

sim_bp = Blueprint('sim', __name__)

@sim_bp.route("/latest", methods=["GET"])
def get_latest():
    latest_entry = Latest.query.first()  # Fetch the first and only entry
    latest_processed_command_id = latest_entry.value if latest_entry else -1
    return jsonify({"latest": latest_processed_command_id})

@sim_bp.route("/register", methods=["POST"])
def register():
    update_latest(request)
    not_from_sim_response = not_req_from_simulator(request)
    if not_from_sim_response:
        return not_from_sim_response
    request_data = request.json
    error = None
    if not request_data.get("username"):
        error = "You have to enter a username"
    elif not request_data.get("email") or "@" not in request_data["email"]:
        error = "You have to enter a valid email address"
    elif not request_data.get("pwd"):
        error = "You have to enter a password"
    elif get_user_id(request_data["username"]) is not None:
        error = "The username is already taken"
    else:
        user = User(username=request_data["username"], email=request_data["email"], pw_hash=generate_password_hash(request_data["pwd"]))
        db.session.add(user)
        db.session.commit()

    if error:
        return jsonify({"status": 400, "error_msg": error}), 400
    else:
        return "", 204


@sim_bp.route("/msgs", methods=["GET"])
def messages():
    """Return all latest messages."""
    update_latest(request)
    not_from_sim_response = not_req_from_simulator(request)
    if not_from_sim_response:
        return not_from_sim_response

    no_msgs = request.args.get("no", type=int, default=100)
    messages = db.session.query(Message, User).join(User, Message.author_id == User.user_id) \
        .filter(Message.flagged == 0) \
        .order_by(Message.pub_date.desc()) \
        .limit(no_msgs).all()

    filtered_msgs = [{
        "content": message.text,
        "pub_date": message.pub_date,
        "user": user.username
    } for message, user in messages]
    return jsonify(filtered_msgs)


@sim_bp.route("/msgs/<username>", methods=["GET", "POST"])
def messages_per_user(username):
    """Returns all messages for a specific user or adds a new message for specified user."""
    update_latest(request)
    not_from_sim_response = not_req_from_simulator(request)
    if not_from_sim_response:
        return not_from_sim_response

    user = User.query.filter_by(username=username).first()
    if not user:
        abort(404)

    if request.method == "GET":
        no_msgs = request.args.get("no", type=int, default=100)
        messages = Message.query.filter_by(author_id=user.user_id, flagged=0).order_by(Message.pub_date.desc()).limit(no_msgs).all()
        filtered_msgs = [{"content": msg.text, "pub_date": msg.pub_date, "user": username} for msg in messages]
        return jsonify(filtered_msgs)

    elif request.method == "POST":
        request_data = request.json
        message = Message(author_id=user.user_id, text=request_data["content"], pub_date=datetime.now().strftime('%Y-%m-%d %H:%M:%S'), flagged=0)
        db.session.add(message)
        db.session.commit()
        return "", 204


@sim_bp.route("/fllws/<username>", methods=["GET", "POST"])
def follow(username):
    update_latest(request)
    not_from_sim_response = not_req_from_simulator(request)
    if not_from_sim_response:
        return not_from_sim_response

    user = User.query.filter_by(username=username).first()
    if not user:
        abort(404)

    if request.method == "POST":
        if "follow" in request.json:
            follows_username = request.json["follow"]
            follows_user = User.query.filter_by(username=follows_username).first()
            if not follows_user:
                abort(404)
            follower = Follower(who_id=user.user_id, whom_id=follows_user.user_id)
            db.session.add(follower)
            db.session.commit()
            return "", 204
        elif "unfollow" in request.json:
            unfollows_username = request.json["unfollow"]
            unfollows_user = User.query.filter_by(username=unfollows_username).first()
            if not unfollows_user:
                abort(404)
            follower = Follower.query.filter_by(who_id=user.user_id, whom_id=unfollows_user.user_id).first()
            if follower:
                db.session.delete(follower)
                db.session.commit()
            return "", 204

