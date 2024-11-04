import requests
import time
import pandas as pd 
from utils import clean_database

BASE_URL = "http://localhost:5000/api"
BASE_DELAY = 0.2

register_data = {
    "username": "yusuke",
    "email": "yusuke@waect.com",
    "pwd": "waect123",
    "pwd2": "waect123"
}

register_data_2 = {
    "username": "kuwabara",
    "email": "kuwabara@waect.com",
    "pwd": "waect123",
    "pwd2": "waect123"
}

message_data = {
    "content": "REIIIII GUNNNNNNN"
}

message_data_2 = {
    "content": "REIIIII KEEEEEEEEEEN"
}

follow_data = {
    "follow": "kuwabara"
}

unfollow_data = {
    "unfollow": "kuwabara"
}

latest_query = {"latest": 1}

message_amount = {"no": 5}

session = requests.Session()

def request_endpoint(path, method="get", data=None, params=None):
    url = f"{BASE_URL}{path}"
    start = time.time()
    response = session.request(method=method, url=url, json=data, params=params)
    end = time.time()
    print(f"Request to {url} | Status: {response.status_code} | Response: {response.text}")
    return {"endpoint": path, "response": response.status_code, "text":response.text, "start": start, "end": end, "delta": end - start}


def sequential_interval_scenario(service, start, iter):

    # 1. Register 2 spirit detectives
    register1 = request_endpoint("/register", method="post", data=register_data, params=latest_query)
    time.sleep(BASE_DELAY)
    register2 = request_endpoint("/register", method="post", data=register_data_2, params=latest_query)
    time.sleep(BASE_DELAY)

    # 2. Post 2 messages
    post_msg_1 = request_endpoint("/msgs/yusuke", method="post", data=message_data, params=latest_query)
    time.sleep(BASE_DELAY)
    post_msg_2 = request_endpoint("/msgs/kuwabara", method="post", data=message_data_2, params=latest_query)
    time.sleep(BASE_DELAY)

    # 3. Retrieve public messages
    public_msgs = request_endpoint("/msgs", method="get", params={**message_amount, **latest_query})
    time.sleep(BASE_DELAY)

    # 4. Retrieve yusuke messages
    user_msgs = request_endpoint("/msgs/yusuke", method="get", params={**message_amount, **latest_query})
    time.sleep(BASE_DELAY)

    # 7. Follow another user
    follows_user = request_endpoint("/fllws/yusuke", method="post", data=follow_data, params=latest_query)
    time.sleep(BASE_DELAY)

    # 6. Get followers for the user
    user_followers = request_endpoint("/fllws/yusuke", method="get", params={**message_amount, **latest_query})
    time.sleep(BASE_DELAY)

    # 7. Unfollow the user
    unfollows_user = request_endpoint("/fllws/yusuke", method="post", data=unfollow_data, params=latest_query)
    time.sleep(BASE_DELAY)

    # 8. Retrieve the latest status
    latest = request_endpoint("/latest", method="get")

    print(f"Finished API sequence for service {service} - iteration {iter}!", flush=True)

    return [f"{service}-api-{start}-{iter}", register1, register2, post_msg_1, post_msg_2, public_msgs, user_msgs, follows_user, user_followers, unfollows_user, latest]


def run_api_seq_scenario(service, start):
    data = []
    for i in range(10):
        iteration_data = sequential_interval_scenario(service, start, i)
        data.append(iteration_data)
        clean_database()

    df = pd.DataFrame(data, columns=["Experiment run ID", "Register 1", "Register 2", "Post Msg 1", 
                                         "Post Msg 2", "Public Messages", "User Messages", 
                                         "Follows User", "User Followers", "Unfollows User", 
                                         "Latest Status"])
    return df