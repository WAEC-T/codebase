import requests
import time
import pandas as pd 
from utils import clean_database

BASE_URL = "http://localhost:5000"
BASE_DELAY = 0.2

register_data = {
    "username": "vegeta",
    "email": "vegeta@waect.com",
    "password": "waect123",
    "password2": "waect123"
}

register_data_kakaroto = {
    "username": "goku",
    "email": "goku@waect.com",
    "password": "waect123",
    "password2": "waect123"
}

login_data = {
    "username": "vegeta",
    "password": "waect123"
}
message_data = {
    "text": "O melhor guerreiro não é aquele que sempre ganha, mas o que mantém o seu orgulho mesmo na derrota."
}

session = requests.Session()

def request_endpoint(path, method="get", data=None):
    url = f"{BASE_URL}{path}"
    start = time.time()
    response = session.request(method=method, url=url, data=data)
    end = time.time()
    print(f"Request to {url} | Status: {response.status_code} | start: {start} | end: {end}")
    return {"endpoint": path, "response": response.status_code, "text":response.text, "start": start, "end": end, "delta": end - start}


def sequential_interval_scenario(service, start, iter):
    # 1. Access Public Timeline
    public_page = request_endpoint("/public")
    time.sleep(BASE_DELAY)

    # 2. Register 2 new super saiyajin users
    register_1 = request_endpoint("/register", method="post", data=register_data)
    time.sleep(BASE_DELAY)
    register_2 = request_endpoint("/register", method="post", data=register_data_kakaroto)
    time.sleep(BASE_DELAY)

    # 3. Login with vegeta credentials
    login = request_endpoint("/login", method="post", data=login_data)
    time.sleep(BASE_DELAY)

    # 4. Access user timeline
    user_timeline = request_endpoint("/vegeta")
    time.sleep(BASE_DELAY)

    # 5. Follow kakaroto - YIKES
    follow = request_endpoint("/goku/follow")
    time.sleep(BASE_DELAY)

    # 6. Post a new message
    add_message = request_endpoint("/add_message", method="post", data=message_data)
    time.sleep(BASE_DELAY)

    # 7. Unfollow kakaroto - YAY
    unfollow = request_endpoint("/goku/unfollow")
    time.sleep(BASE_DELAY)

    # 8. Access Timeline again (logged-in user)
    user_timeline_redirect = request_endpoint("/")
    time.sleep(BASE_DELAY)

    # 9. Logout
    logout = request_endpoint("/logout")

    print(f"Finished page sequence for service {service} - iteration {iter}!", flush=True)

    return [f"{service}-page-{start}-{iter}", public_page, register_1, register_2, login, user_timeline, follow, add_message, unfollow, user_timeline_redirect, logout]

def run_page_seq_scenario(service, start):
    data = []
    for i in range(20):
        iteration_data = sequential_interval_scenario(service, start, i)
        data.append(iteration_data)
        clean_database()

    df = pd.DataFrame(data, columns=[
        "Experiment ID", 
        "Public Page", 
        "Register 1", 
        "Register 2", 
        "Login", 
        "User Timeline", 
        "Follow", 
        "Add Message", 
        "Unfollow", 
        "User Timeline Redirect", 
        "Logout"
    ])
    return df