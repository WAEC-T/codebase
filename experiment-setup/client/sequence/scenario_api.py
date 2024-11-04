import requests
import time

BASE_URL = "http://localhost:5000/api"

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
    response = session.request(method=method, url=url, json=data, params=params)
    print(f"Request to {url} | Status: {response.status_code} | Response: {response.text}")
    return response

# 1. Register 2 spirit detectives
request_endpoint("/register", method="post", data=register_data, params=latest_query)
time.sleep(2)
request_endpoint("/register", method="post", data=register_data_2, params=latest_query)
time.sleep(2)

# 2. Post 2 messages
request_endpoint("/msgs/yusuke", method="post", data=message_data, params=latest_query)
time.sleep(2)
request_endpoint("/msgs/kuwabara", method="post", data=message_data_2, params=latest_query)
time.sleep(2)

# 3. Retrieve public messages
request_endpoint("/msgs", method="get", params={**message_amount, **latest_query})
time.sleep(2)

# 4. Retrieve yusuke messages
request_endpoint("/msgs/yusuke", method="get", params={**message_amount, **latest_query})
time.sleep(2)

# 7. Follow another user
request_endpoint("/fllws/yusuke", method="post", data=follow_data, params=latest_query)
time.sleep(2)

# 6. Get followers for the user
request_endpoint("/fllws/yusuke", method="get", params={**message_amount, **latest_query})
time.sleep(2)

# 7. Unfollow the user
request_endpoint("/fllws/yusuke", method="post", data=unfollow_data, params=latest_query)
time.sleep(2)

# 8. Retrieve the latest status
request_endpoint("/latest", method="get")
time.sleep(2)

print("Finished API sequence!")
