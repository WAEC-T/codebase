import base64
import psycopg2
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.firefox.options import Options

from helpers.postgres_helper import DATABASE_URL
from helpers.test_helper import BASE_URL

# Updated to match API test file settings
USERNAME = "simulator"
PWD = "super_safe!"
CREDENTIALS = ":".join([USERNAME, PWD]).encode("ascii")
ENCODED_CREDENTIALS = base64.b64encode(CREDENTIALS).decode()
HEADERS = {
    "Connection": "close",
    "Content-Type": "application/json",
    "Authorization": f"Basic {ENCODED_CREDENTIALS}",
}


def get_text_from_first_li(driver):
    try:
        flashes_ul = driver.find_element(By.CLASS_NAME, "flashes")
        li_elements = flashes_ul.find_elements(By.TAG_NAME, "li")
        if li_elements and li_elements[0].text.strip():
            return li_elements[0].text.strip()
    except:
        return None


def _register_user_via_gui(driver, data):
    register_url = f"{BASE_URL}/register"
    driver.get(register_url)

    wait = WebDriverWait(driver, 15)
    input_fields = driver.find_elements(By.TAG_NAME, "input")

    for idx, str_content in enumerate(data):
        input_fields[idx].send_keys(str_content)
    input_fields[4].send_keys(Keys.RETURN)

    get_text_from_first_li(driver)

    wait = WebDriverWait(driver, 15)
    li_text = wait.until(get_text_from_first_li)
    return li_text


def _login_user_via_gui(driver, username, password):
    login_url = f"{BASE_URL}/login"
    driver.get(login_url)

    wait = WebDriverWait(driver, 15)
    input_fields = driver.find_elements(By.TAG_NAME, "input")
    print(input_fields)

    # Assuming the login form has two input fields: username and password
    input_fields[0].send_keys(username)
    input_fields[1].send_keys(password)
    input_fields[2].send_keys(Keys.RETURN)

    get_text_from_first_li(driver)

    wait = WebDriverWait(driver, 15)
    li_text = wait.until(get_text_from_first_li)
    return li_text


def _logout_user_via_gui(driver):
    logout_url = f"{BASE_URL}/logout"
    driver.get(logout_url)

    get_text_from_first_li(driver)

    wait = WebDriverWait(driver, 15)
    li_text = wait.until(get_text_from_first_li)
    return li_text


def _get_user_by_name(name):
    with psycopg2.connect(DATABASE_URL) as conn:
        with conn.cursor() as cur:
            cur.execute(f"SELECT username FROM users WHERE username='{name}';")
            return cur.fetchone()


def test_register_user_via_gui():
    """
    This is a UI test. It only interacts with the UI that is rendered in the browser and checks that visual
    responses that users observe are displayed.
    """
    firefox_options = Options()
    firefox_options.add_argument("--headless")  # for visibility
    with webdriver.Firefox(options=firefox_options) as driver:
        generated_msg = _register_user_via_gui(
            driver, ["user1", "user1@some.where", "waect", "waect"]
        )
        expected_msg = "You were successfully registered and can login now"
        assert generated_msg == expected_msg


def test_register_user_via_gui_and_check_db_entry():
    """
    This is an end-to-end test. Before registering a user via the UI, it checks that no such user exists in the
    database yet. After registering a user, it checks that the respective user appears in the database.
    """
    firefox_options = Options()
    firefox_options.add_argument("--headless")  # for visibility
    with webdriver.Firefox(options=firefox_options) as driver:
        # Check that user does not exist before registration
        assert _get_user_by_name("user2") is None

        generated_msg = _register_user_via_gui(
            driver, ["user2", "user2@some.where", "waect", "waect"]
        )
        expected_msg = "You were successfully registered and can login now"
        assert generated_msg == expected_msg

        # Check that user now exists in the database
        assert _get_user_by_name("user2")[0] == "user2"


def test_login_flash_message():
    """
    This is an end-to-end test. It checks that the flash message "You were logged in" is displayed after a user logs in.
    """
    firefox_options = Options()
    firefox_options.add_argument("--headless")  # for visibility
    with webdriver.Firefox(options=firefox_options) as driver:
        _register_user_via_gui(driver, ["user3", "user3@some.where", "waect", "waect"])
        generated_msg = _login_user_via_gui(driver, "user3", "waect")
        expected_msg = "You were logged in"
        assert generated_msg == expected_msg


def test_logout_flash_message():
    """
    This is an end-to-end test. It checks that the flash message "You were logged out" is displayed after a user logs out.
    """
    firefox_options = Options()
    firefox_options.add_argument("--headless")  # for visibility
    with webdriver.Firefox(options=firefox_options) as driver:
        _register_user_via_gui(driver, ["user4", "user1@some.where", "waect", "waect"])
        _login_user_via_gui(driver, "user4", "waect")
        generated_msg = _logout_user_via_gui(driver)
        expected_msg = "You were logged out"
        assert generated_msg == expected_msg
