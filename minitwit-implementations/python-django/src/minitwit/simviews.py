import json
import logging

from django.contrib.auth.models import User
from django.db import connection
from django.http import HttpResponse, JsonResponse, HttpResponseForbidden
from django.views.decorators.csrf import csrf_exempt

from . import models

LATEST = 0

MAXMESSAGES = 30


# Get latest
def latest(request):
    global LATEST
    auth_check = sim_middleware(request)
    if auth_check:
        return auth_check  # Returns 403 response if unauthorized
    return JsonResponse({"latest": int(LATEST)})


# Register user
@csrf_exempt
def register(request):
    auth_check = sim_middleware(request)
    if auth_check:
        return auth_check  # Returns 403 response if unauthorized

    username = json.loads(request.body)["username"]
    email = json.loads(request.body)["email"]
    pwd = json.loads(request.body)["pwd"]

    update_latest(request)
    error = None
    if request.method == "POST":
        if username is None:
            error = "You have to enter a username"
        elif email is None or "@" not in email:
            error = "You have to enter a valid email address"
        elif pwd is None:
            error = "You have to enter a password"
        else:
            try:
                User.objects.get(username=username)
                error = "The username is already taken"
            except:
                User.objects.create_user(username, email, pwd)
    if error:
        return JsonResponse({"status": 400, "error_msg": error}, status=400)

    return HttpResponse("", status=204)


# Get messages
@csrf_exempt
def all_msgs(request):
    auth_check = sim_middleware(request)
    if auth_check:
        return auth_check  # Returns 403 response if unauthorized

    update_latest(request)

    amount = int(request.GET.get("no", MAXMESSAGES))

    messages = models.Message.objects.all().order_by(
        "-pub_date").values()[:amount]
    messages = [dict(message) for message in list(messages)]

    for message in messages:
        message["content"] = message["text"]
        message["user"] = User.objects.get(id=message["user_id"]).username

    return JsonResponse(messages, safe=False)


# Get user messages / Post message
@csrf_exempt
def user_msgs(request, username):
    auth_check = sim_middleware(request)
    if auth_check:
        return auth_check  # Returns 403 response if unauthorized

    update_latest(request)

    if request.method == "GET":
        amount = int(request.GET.get("no", MAXMESSAGES))
        user = User.objects.get(username=username)

        messages = (
            models.Message.objects.filter(
                user=user).order_by("-pub_date").values()[:amount]
        )
        messages = [dict(message) for message in list(messages)]

        for message in messages:
            message["content"] = message["text"]
            message["user"] = User.objects.get(id=message["user_id"]).username

        return JsonResponse(messages, safe=False)

    elif request.method == "POST":
        try:
            user = User.objects.get(username=username)
        except:
            return JsonResponse(
                {"status": 404, "error_msg": "User not found"}, status=404
            )
        message = models.Message.objects.create(
            text=json.loads(request.body)["content"], user=user, flagged=0
        )

    return HttpResponse(None, status=204)


# Follow/unfollow user
@csrf_exempt
def follow_user(request, username):
    auth_check = sim_middleware(request)
    if auth_check:
        return auth_check  # Returns 403 response if unauthorized

    update_latest(request)

    if request.method == "POST":
        try:
            user = User.objects.get(username=username)
        except:
            return JsonResponse(
                {"status": 404, "error_msg": "User not found"}, status=404
            )
        try:
            body = json.loads(request.body)
            if "follow" in body:
                user_follow = body["follow"]
                follower = User.objects.get(username=user_follow)
                models.Follower.objects.create(who_id=user, whom_id=follower)
                return JsonResponse({"status": 200, "msg": "Followed successfully"})

            elif "unfollow" in body:
                user_unfollow = body["unfollow"]
                follower = User.objects.get(username=user_unfollow)
                models.Follower.objects.filter(
                    who_id=user, whom_id=follower).delete()
                return JsonResponse({"status": 200, "msg": "Unfollowed successfully"})

            else:
                return JsonResponse(
                    {"status": 400, "error_msg": "Request must contain either 'follow' or 'unfollow'"},
                    status=400
                )

        except User.DoesNotExist:
            return JsonResponse(
                {"status": 404, "error_msg": "User not found"}, status=404
            )

        except Exception as e:
            logging.error(f"Unexpected error: {e}")
            return JsonResponse(
                {"status": 500, "error_msg": "Internal server error"}, status=500
            )
    elif request.method == "GET":
        amount = int(request.GET.get("no", MAXMESSAGES))
        try:
            user = User.objects.get(username=username)
        except:
            return JsonResponse(
                {"status": 404, "error_msg": "User not found"}, status=404
            )
        followers = models.Follower.objects.filter(
            who_id=user).values()[:amount]
        followers = [dict(follow) for follow in list(followers)]

        followers_usernames = [
            User.objects.get(id=fo["whom_id_id"]).username for fo in followers
        ]

        return JsonResponse({"follows": followers_usernames})
    return HttpResponse(None, status=204)


# Helper functions
def query_db(query, args=()):
    """Queries the database and returns a list of dictionaries."""
    with connection.cursor() as cursor:
        cursor.execute(query, args)
        rv = [
            dict((cursor.description[idx][0], value)
                 for idx, value in enumerate(row))
            for row in cursor.fetchall()
        ]
    return (rv[0] if rv else None) if one else rv


def sim_middleware(request):
    from_simulator = request.headers.get("Authorization")
    if from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh":
        return HttpResponseForbidden("You are not authorized to use this resource!")


# Update latest
def update_latest(request):
    global LATEST
    try_latest = request.GET["latest"]
    LATEST = try_latest if try_latest != -1 else LATEST
    LATEST = int(LATEST)
