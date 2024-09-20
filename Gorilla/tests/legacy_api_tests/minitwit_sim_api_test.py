import json
import base64
import os
import requests

BASE_URL = os.getenv("MT_BASEURL", "http://localhost:15001")
USERNAME = 'simulator'
PWD = 'super_safe!'
CREDENTIALS = ':'.join([USERNAME, PWD]).encode('ascii')
ENCODED_CREDENTIALS = base64.b64encode(CREDENTIALS).decode()
HEADERS = {'Connection': 'close',
           'Content-Type': 'application/json',
           f'Authorization': f'Basic {ENCODED_CREDENTIALS}'}
    
def test_cleandb():
        url = f"{BASE_URL}/cleandb"
        response = requests.post(url, headers=HEADERS)
        assert response.ok, response

def test_latest():
        # post something to updaet LATEST
        url = f"{BASE_URL}/register"
        data = {'username': 'test', 'email': 'test@test', 'pwd': 'foo'}
        params = {'latest': 1337}
        response = requests.post(url, data=json.dumps(data),
                                params=params, headers=HEADERS)
        assert response.ok, response

        # verify that latest was updated
        url = f'{BASE_URL}/latest'
        response = requests.get(url, headers=HEADERS)
        assert response.ok, response
        assert response.json()['latest'] == 1337, response


def test_register():
        username = 'a'
        email = 'a@a.a'
        pwd = 'a'
        data = {'username': username, 'email': email, 'pwd': pwd}
        params = {'latest': 1}
        response = requests.post(f'{BASE_URL}/register',
                                data=json.dumps(data), headers=HEADERS, params=params)
        assert response.ok, response
        # TODO: add another assertion that it is really there

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 1, response


def test_create_msg():
        username = 'a'
        data = {'content': 'Blub!'}
        url = f'{BASE_URL}/msgs/{username}'
        params = {'latest': 2}
        response = requests.post(url, data=json.dumps(data),
                                headers=HEADERS, params=params)
        assert response.ok, response

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 2, response


def test_get_latest_user_msgs():
        username = 'a'

        query = {'no': 20, 'latest': 3}
        url = f'{BASE_URL}/msgs/{username}'
        response = requests.get(url, headers=HEADERS, params=query)
        assert response.status_code == 200, response

        got_it_earlier = False
        for msg in response.json():
            if msg['content'] == 'Blub!' and msg['user'] == username:
                got_it_earlier = True

        assert got_it_earlier, response

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 3, response


def test_get_latest_msgs():
        username = 'a'
        query = {'no': 20, 'latest': 4}
        url = f'{BASE_URL}/msgs'
        response = requests.get(url, headers=HEADERS, params=query)
        assert response.status_code == 200, response

        got_it_earlier = False
        for msg in response.json():
            if msg['content'] == 'Blub!' and msg['user'] == username:
                got_it_earlier = True

        assert got_it_earlier, response

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 4, response


def test_register_b():
        username = 'b'
        email = 'b@b.b'
        pwd = 'b'
        data = {'username': username, 'email': email, 'pwd': pwd}
        params = {'latest': 5}
        response = requests.post(f'{BASE_URL}/register', data=json.dumps(data),
                                headers=HEADERS, params=params)
        assert response.ok, response
        # TODO: add another assertion that it is really there

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 5, response


def test_register_c():
        username = 'c'
        email = 'c@c.c'
        pwd = 'c'
        data = {'username': username, 'email': email, 'pwd': pwd}
        params = {'latest': 6}
        response = requests.post(f'{BASE_URL}/register', data=json.dumps(data),
                                headers=HEADERS, params=params)
        assert response.ok, response

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 6, response


def test_follow_user():
        username = 'a'
        url = f'{BASE_URL}/fllws/{username}'
        data = {'follow': 'b'}
        params = {'latest': 7}
        response = requests.post(url, data=json.dumps(data),
                                headers=HEADERS, params=params)
        assert response.ok, response

        data = {'follow': 'c'}
        params = {'latest': 8}
        response = requests.post(url, data=json.dumps(data),
                                headers=HEADERS, params=params)
        assert response.ok, response

        query = {'no': 20, 'latest': 9}
        response = requests.get(url, headers=HEADERS, params=query)
        assert response.ok, response

        json_data = response.json()
        assert "b" in json_data["follows"], response
        assert "c" in json_data["follows"], response

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 9, response


def test_a_unfollows_b():
        username = 'a'
        url = f'{BASE_URL}/fllws/{username}'

        #  first send unfollow command
        data = {'unfollow': 'b'}
        params = {'latest': 10}
        response = requests.post(url, data=json.dumps(data),
                                headers=HEADERS, params=params)
        assert response.ok, response

        # then verify that b is no longer in follows list
        query = {'no': 20, 'latest': 11}
        response = requests.get(url, params=query, headers=HEADERS)
        assert response.ok, response
        assert 'b' not in response.json()['follows'], response

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 11, response 
