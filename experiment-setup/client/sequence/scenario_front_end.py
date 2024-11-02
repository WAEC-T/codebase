import requests
import time

BASE_URL = "http://localhost:5000"

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
    response = session.request(method=method, url=url, data=data)
    print(f"Request to {url} | Status: {response.status_code}")
    return response

# 1. Access Public Timeline
request_endpoint("/public")
time.sleep(2)

# 2. Register 2 new super saiyajin users
request_endpoint("/register", method="post", data=register_data)
time.sleep(2)
request_endpoint("/register", method="post", data=register_data_kakaroto)
time.sleep(2)

# 3. Login with vegeta credentials
request_endpoint("/login", method="post", data=login_data)
time.sleep(2)

# 4. Access user timeline
request_endpoint("/vegeta")
time.sleep(2)

# 5. Follow kakaroto - YIKES
request_endpoint("/goku/follow")
time.sleep(2)

# 6. Post a new message
request_endpoint("/add_message", method="post", data=message_data)
time.sleep(2)

# 7. Unfollow kakaroto - YAY
request_endpoint("/goku/unfollow")
time.sleep(2)

# 8. Access Timeline again (logged-in user)
request_endpoint("/")
time.sleep(2)

# 9. Logout
request_endpoint("/logout")
time.sleep(2)

print("Finished frontend sequence!")