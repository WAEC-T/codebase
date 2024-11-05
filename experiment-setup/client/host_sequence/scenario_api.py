import requests
import time
import pandas as pd
from utils import clean_database, print_info_call
from const_data import api_register_data, message_data, api_fllw_data, api_unfllw_data
import sys

BASE_URL = "http://localhost:5000/api"
AUTH_HEADER = {"Authorization": "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh"}
BASE_DELAY = 1.8
ITER_NUM = 20 # iteration number for each endpoint call has to be <= 400

user_sessions = {}
api_latest_query = {"latest": 1}
api_message_amount = {"no": 5}

session = requests.Session()

def request_endpoint(path, method="get", data=None, params=None, user_session=None):
    url = f"{BASE_URL}{path}"
    start = time.time()
    response = (user_session or session).request(method=method, url=url, json=data, params=params, headers=AUTH_HEADER)
    end = time.time()
    print(f"Request to {url:<40} | Status: {response.status_code:<3} | start: {start:<20.6f} | end: {end:<20.6f}")
    return {"endpoint": path, "response": response.status_code, "text": response.text, "start": start, "end": end,
            "delta": end - start}


def sequential_interval_scenario(service, start, iter):
    # set main user
    main_user = api_register_data[0]["username"]

    # 1. Register all users
    print_info_call("API", service, "Register", ITER_NUM)
    api_register_response = []
    for user_register_data in api_register_data[:ITER_NUM]:
        response = request_endpoint("/register", method="post", data=user_register_data, params=api_latest_query)
        api_register_response.append(response)
        user_sessions[user_register_data["username"]] = requests.Session()
    time.sleep(BASE_DELAY)

    # 2. User 1 post messages
    print_info_call("API", service, "Post messages", ITER_NUM)
    main_user_session = user_sessions[main_user]
    api_message_response = [
        request_endpoint(f"/msgs/{main_user}", method="post", data=message_data, params=api_latest_query) for _ in range(ITER_NUM)
    ]
    time.sleep(BASE_DELAY)

    # 3. Retrieve public messages
    print_info_call("API", service, "Retrieve public messages", ITER_NUM)
    api_public_msgs = [request_endpoint("/msgs", method="get", params={**api_message_amount, **api_latest_query}) for _
                       in range(ITER_NUM)]
    time.sleep(BASE_DELAY)

    # 4. Retrieve user messages
    print_info_call("API", service, "Retrieve user messages", ITER_NUM)
    api_user_msgs = [
        request_endpoint(f"/msgs/{main_user}", method="get", params={**api_message_amount, **api_latest_query}) for _ in
        range(ITER_NUM)]
    time.sleep(BASE_DELAY)

    # 5. All users follow user1
    print_info_call("API", service, "Follow users", ITER_NUM)
    api_follow_user_1 = []
    for follow in api_fllw_data[1:ITER_NUM]:
        user_session = user_sessions[follow["follow"]]
        response = request_endpoint(f"/fllws/{follow['follow']}", method="post", data={"follow": "user1"},
                                    params=api_latest_query, user_session=user_session)
        api_follow_user_1.append(response)
    time.sleep(BASE_DELAY)

    # 6. Get followers for user one
    print_info_call("API", service, "Get followers", ITER_NUM)
    api_user_followers = [request_endpoint("/fllws/user1", method="get",
                                           params={**api_message_amount, **api_latest_query},
                                           user_session=main_user_session) for _ in range(ITER_NUM)]
    time.sleep(BASE_DELAY)

    # 7. All users unfollow user1
    print_info_call("API", service, "Unfollow users", ITER_NUM)
    api_unfollow_user_1 = []
    for unfollow in api_unfllw_data[1:ITER_NUM]:
        user_session = user_sessions[unfollow["unfollow"]]
        response = request_endpoint(f"/fllws/{unfollow['unfollow']}", method="post", data={"unfollow": "user1"},
                                    params=api_latest_query, user_session=user_session)
        api_unfollow_user_1.append(response)
    time.sleep(BASE_DELAY)

    # 8. Retrieve the latest status
    print_info_call("API", service, "Retrieve latest status", ITER_NUM)
    api_latest = [request_endpoint("/latest", method="get") for _ in range(ITER_NUM)]

    print(f"Finished API sequence for service {service} - iteration {iter}!", flush=True)

    return [f"{service}-api-{start}-{iter}",
            api_register_response,
            api_message_response,
            api_public_msgs,
            api_user_msgs,
            api_follow_user_1,
            api_user_followers,
            api_unfollow_user_1,
            api_latest]


def run_api_seq_scenario(service, start):
    data = [sequential_interval_scenario(service, start, 0)]

    clean_database()

    df = pd.DataFrame(data, columns=["ExperimentID",
                                     "ApiRegisterResponse",
                                     "ApiMessageResponse",
                                     "ApiPublicMsgsResponse",
                                     "ApiUserMsgsResponse",
                                     "ApiFollowsUserResponse",
                                     "ApiUserFollowersResponse",
                                     "ApiUnfollowsUserResponse",
                                     "ApiLatestResponse"
                                     ])
    return df
