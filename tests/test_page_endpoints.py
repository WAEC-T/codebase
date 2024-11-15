import base64
import psycopg2
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.firefox.options import Options

# Updated to match API test file settings
USERNAME = 'simulator'
PWD = 'super_safe!'
CREDENTIALS = ':'.join([USERNAME, PWD]).encode('ascii')
ENCODED_CREDENTIALS = base64.b64encode(CREDENTIALS).decode()
HEADERS = {
    'Connection': 'close',
    'Content-Type': 'application/json',
    'Authorization': f'Basic {ENCODED_CREDENTIALS}'
}

# Get the database URL from the environment variable
DATABASE_URL = "postgresql://user:pass@localhost:5432/waect"

def get_text_from_first_li(driver):
    try:
        flashes_ul = driver.find_element(By.CLASS_NAME, "flashes")
        li_elements = flashes_ul.find_elements(By.TAG_NAME, "li")
        if li_elements and li_elements[0].text.strip():
            return li_elements[0].text.strip()
    except:
        return None

def _register_user_via_gui(driver, data):
    register_url = "http://localhost:5000/register"
    driver.get(register_url)
    
    wait = WebDriverWait(driver, 30)
    input_fields = driver.find_elements(By.TAG_NAME, "input")

    for idx, str_content in enumerate(data):
        input_fields[idx].send_keys(str_content)
    input_fields[4].send_keys(Keys.RETURN)

    get_text_from_first_li(driver)

    wait = WebDriverWait(driver, 30)
    li_text = wait.until(get_text_from_first_li)
    return li_text

def _login_user_via_gui(driver, username, password):
    login_url = "http://localhost:5000/login"
    driver.get(login_url)

    wait = WebDriverWait(driver, 30)
    input_fields = driver.find_elements(By.TAG_NAME, "input")
    print(input_fields)

    # Assuming the login form has two input fields: username and password
    input_fields[0].send_keys(username)
    input_fields[1].send_keys(password)
    input_fields[2].send_keys(Keys.RETURN)

    get_text_from_first_li(driver)

    wait = WebDriverWait(driver, 30)
    li_text = wait.until(get_text_from_first_li)
    return li_text

def _logout_user_via_gui(driver):
    logout_url = "http://localhost:5000/logout"
    driver.get(logout_url)

    get_text_from_first_li(driver)

    wait = WebDriverWait(driver, 30)
    li_text = wait.until(get_text_from_first_li)
    return li_text

# def _post_message_via_gui(driver, user, message_text):
#     post_message_url = f"http://localhost:5001/add_message"
#     driver.get(post_message_url)
#
#     wait = WebDriverWait(driver, 5)
#     wait.until(EC.presence_of_element_located((By.NAME, "text")))
#
#     try:
#         input_field = wait.until(EC.presence_of_element_located((By.NAME, "text")))
#     except Exception as e:
#         print("Failed to find the input field for posting message.")
#         print(driver.page_source)  # Log page source for debugging
#         raise e
#
#     input_field.send_keys(message_text)
#
#     # Click the submit button to post the message
#     try:
#         submit_button = driver.find_element(By.CSS_SELECTOR, "input[type='submit']")
#         submit_button.click()
#     except Exception as e:
#         print("Failed to find or click the submit button.")
#         print(driver.page_source)  # Log page source for debugging
#         raise e
#
#     # Enter the message text
#     #input_field = driver.find_element(By.NAME, "text")
#     # input_field.send_keys(message_text)
#
#     # Click the submit button to post the message
#     submit_button = driver.find_element(By.CSS_SELECTOR, "input[type='submit']")
#     submit_button.click()
#
#     get_text_from_first_li(driver)
#
#     wait = WebDriverWait(driver, 5)
#     li_text = wait.until(EC.visibility_of_element_located((By.CLASS_NAME, "flashes")))
#     return li_text

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
        generated_msg = _register_user_via_gui(driver, ["user1", "user1@some.where", "waect", "waect"])
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

        generated_msg = _register_user_via_gui(driver, ["user2", "user2@some.where", "waect", "waect"])
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

# def test_post_message_flash_message():
#     """
#     This is an end-to-end test. It checks that the flash message "Your message was recorded" is displayed after a user posts a message.
#     """
#     firefox_options = Options()
#     firefox_options.add_argument("--headless")  # for visibility
#     with webdriver.Firefox(options=firefox_options) as driver:
#         _register_user_via_gui(driver, ["user2", "user1@some.where", "waect123", "waect123"])
#         _login_user_via_gui(driver, "user2", "waect123")
#         generated_msg = _post_message_via_gui(driver, "user2", "Die größte Befriedigung erfährt man, wenn man einer Sache über eine lange Zeit sein Herz und seine Seele schenkt – und sie es wert ist. — Steve Jobs")
#         expected_msg = "Your message was recorded"
#         assert generated_msg == expected_msg