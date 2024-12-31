import requests
import time
import pandas as pd
from utils import print_info_call
from host_sequence.const_data import api_register_data_dummie, api_follow_data_dummie, api_unfollow_data_dummie, api_message_data

BASE_URL = "YOUR_IP_HERE:PORT/api"
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
    return {"endpoint": f"/api{path}", "response": response.status_code, "start": start, "end": end,
            "delta": end - start}


def sequential_interval_scenario(service, start, iter):
    # set main user and dataframe
    main_user = api_register_data_dummie(0)["username"]
    results = []

    # 0. Health check call public endpoint
    print_info_call("API", service, "Health check", 1)
    request_endpoint("/public", method="get")

    # 1. Register all users
    print_info_call("API", service, "Register", ITER_NUM)
    for i in range(ITER_NUM):
        user_register_data = api_register_data_dummie(i)
        response = request_endpoint("/register", method="post", data=user_register_data, params=api_latest_query)
        results.append(response)
        user_sessions[user_register_data["username"]] = requests.Session()
    time.sleep(BASE_DELAY)

    # 2. User 1 post messages
    print_info_call("API", service, "Post messages", ITER_NUM)
    main_user_session = user_sessions[main_user]
    for _ in range(ITER_NUM):
        response = request_endpoint(f"/msgs/{main_user}", method="post", data=api_message_data, params=api_latest_query)
        results.append(response)
    time.sleep(BASE_DELAY)

    # 3. Retrieve public messages
    print_info_call("API", service, "Retrieve public messages", ITER_NUM)
    for _ in range(ITER_NUM):
        response = request_endpoint("/msgs", method="get", params={**api_message_amount, **api_latest_query})
        results.append(response)
    time.sleep(BASE_DELAY)

    # 4. Retrieve user messages
    print_info_call("API", service, "Retrieve user messages", ITER_NUM)
    for _ in range(ITER_NUM):
        response = request_endpoint(f"/msgs/{main_user}", method="get", params={**api_message_amount, **api_latest_query})
        results.append(response)
    time.sleep(BASE_DELAY)

    # 5. All users follow user1
    print_info_call("API", service, "Follow users", ITER_NUM)
    for i in range(1, ITER_NUM):
        follow = api_follow_data_dummie(i)
        user_session = user_sessions[follow["follow"]]
        response = request_endpoint(f"/fllws/{follow['follow']}", method="post", data={"follow": "user0"},
                                    params=api_latest_query, user_session=user_session)
        results.append(response)
    time.sleep(BASE_DELAY)

    # 6. Get followers for user one
    print_info_call("API", service, "Get followers", ITER_NUM)
    for _ in range(ITER_NUM):
        response = request_endpoint(f"/fllws/{main_user}", method="get", params={**api_message_amount, **api_latest_query},
                                    user_session=main_user_session)
        results.append(response)
    time.sleep(BASE_DELAY)

    # 7. All users unfollow user1
    print_info_call("API", service, "Unfollow users", ITER_NUM)
    for i in range(1, ITER_NUM):
        unfollow = api_unfollow_data_dummie(i)
        user_session = user_sessions[unfollow["unfollow"]]
        response = request_endpoint(f"/fllws/{unfollow['unfollow']}", method="post", data={"unfollow": "user0"},
                                    params=api_latest_query, user_session=user_session)
        response['endpoint'] = response['endpoint'].replace('fllws', 'unfllw')
        results.append(response)
    time.sleep(BASE_DELAY)

    # 8. Retrieve the latest status
    print_info_call("API", service, "Retrieve latest status", ITER_NUM)
    for _ in range(ITER_NUM):
        response = request_endpoint("/latest", method="get", params=api_latest_query)
        results.append(response)

    print(f"Finished API sequence for service {service} - iteration {iter}!", flush=True)

    return results


def run_api_seq_scenario(service, start):
    responses = sequential_interval_scenario(service, start, 0)
    df = pd.DataFrame(responses)

    return df
