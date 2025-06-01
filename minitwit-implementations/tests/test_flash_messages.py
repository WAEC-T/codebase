import pytest
import requests
import psycopg2

from helpers import test_helper
from helpers.postgres_helper import DATABASE_URL
from test_api_endpoints import BASE_URL


def register_user(username, email, password, user_session=None):
    """
    Register a user with the given username, email, and password.

    :param username: Username string
    :param email: Email string
    :param password: Password string
    :param user_session: Optional requests session
    :return: Response object from the request
    """
    register_url = f"{BASE_URL}/register"
    data = {
        "username": username,
        "email": email,
        "password": password,
        "password2": password,
    }
    return user_session.post(register_url, data=data, allow_redirects=True)


def login_user(username, password, user_session=None):
    """
    Log in a user with the given username and password.

    :param username: Username string
    :param password: Password string
    :param user_session: Optional requests session
    :return: Response object from the request
    """
    login_url = f"{BASE_URL}/login"
    data = {
        "username": username,
        "password": password,
    }
    return user_session.post(login_url, data=data, allow_redirects=True)


def follow_user(follow_username, user_session=None):
    """
    Follow another user.

    :param follow_username: Username of the user to follow
    :param user_session: Optional requests session
    :return: Response object from the request
    """
    follow_url = f"{BASE_URL}/{follow_username}/follow"
    return user_session.get(follow_url, allow_redirects=True)


def unfollow_user(unfollow_username, user_session=None):
    """
    Unfollow another user.

    :param unfollow_username: Username of the user to unfollow
    :param user_session: Optional requests session
    :return: Response object from the request
    """
    unfollow_url = f"{BASE_URL}/{unfollow_username}/unfollow"
    return user_session.get(unfollow_url, allow_redirects=True)


def post_message(message_text, user_session=None):
    """
    Post a message.

    :param message_text: Message text to post
    :param user_session: Optional requests session
    :return: Response object from the request
    """
    post_message_url = f"{BASE_URL}/add_message"
    data = {"text": message_text}
    return user_session.post(post_message_url, data=data, allow_redirects=True)


def logout_user(user_session=None):
    """
    Log out the current user.

    :param user_session: Optional requests session
    :return: Response object from the request
    """
    logout_url = f"{BASE_URL}/logout"
    return user_session.get(logout_url, allow_redirects=True)


def clean_database():
    """
    Truncate all tables in the database to clean the state.
    """
    with psycopg2.connect(DATABASE_URL) as conn:
        with conn.cursor() as cur:
            cur.execute("TRUNCATE TABLE users CASCADE;")
            cur.execute("TRUNCATE TABLE followers CASCADE;")
            cur.execute("TRUNCATE TABLE messages CASCADE;")
            conn.commit()


@pytest.fixture(scope="module")
def user1_session():
    """
    Create and return a session for user1. This session persists for all tests in the module.
    """
    session = requests.Session()
    return session


@pytest.fixture(scope="module", autouse=True)
def fetch_public_page():
    """
    Fetches the public page before all tests as health check
    """
    public_page_url = f"{BASE_URL}/public"
    response = requests.get(public_page_url)
    assert response.status_code == 200
    return response


@pytest.fixture(scope="session", autouse=True)
def cleanup_db_after_tests():
    """
    Automatically clean the database after all tests in the session.
    """
    yield  # Let all tests execute
    clean_database()


### Test Cases ########################################################################################################


def test_register_flash(user1_session):
    """
    Test flash message for user registration.
    """
    expected_flash = "You were successfully registered and can login now"
    response = register_user("user1", "user1@waect.com", "waect", user1_session)
    assert response.status_code == 200, "register failed"
    assert expected_flash in response.text


def test_login_flash(user1_session):
    """
    Test flash message for user login.
    """
    expected_flash = "You were logged in"
    response = login_user("user1", "waect", user1_session)
    assert response.status_code == 200, "login failed"
    assert expected_flash in response.text


def test_user1_follow_user2_flash(user1_session):
    """
    Test flash message for user1 following user2.
    """
    user2_session = requests.Session()
    register_user("user2", "user2@waect.com", "waect", user2_session)
    expected_flash = f"You are now following user2"
    response = follow_user("user2", user1_session)
    assert response.status_code == 200, "follow failed"
    assert expected_flash in response.text


def test_user1_unfollow_user2_flash(user1_session):
    """
    Test flash message for user1 unfollowing user2.
    """
    expected_flash = "You are no longer following user2"
    response = unfollow_user("user2", user1_session)
    assert response.status_code == 200, "unfollow failed"
    assert expected_flash in response.text


def test_user1_post_message_flash(user1_session):
    """
    Test flash message for user1 posting a message.
    """
    expected_flash = "Your message was recorded"
    response = post_message("Hello, world!", user1_session)
    assert response.status_code == 200, "post message failed"
    assert expected_flash in response.text


def test_logout_flash(user1_session):
    """
    Test flash message for user logout.
    """
    expected_flash = "You were logged out"
    response = logout_user(user1_session)
    assert response.status_code == 200, "logout failed"
    assert expected_flash in response.text
