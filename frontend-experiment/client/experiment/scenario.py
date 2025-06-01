import asyncio
import random
import string
from playwright.async_api import async_playwright, Page
import time
import json
import requests
from experiment.actions import navigate_to, press_link, fill_input, submit_input
import logging
import os
from datetime import datetime
import sys

BASE_URL = "http://host.docker.internal:5000"


def create_new_session():
    session = requests.Session()
    session.headers.update(
        {
            "Connection": "close",
            "Content-Type": "application/json",
            "Authorization": "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh",
        }
    )
    return session


def test_create_msg(number, username):
    session = create_new_session()
    # username = "hej"
    data = {"content": str(number) + "-Blub!"}
    url = f"{BASE_URL}/api/msgs/{username}"
    header = {
        "Authorization": "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh",
        "Content-Type": "application/json",
    }
    params = {"latest": 2}
    response = session.post(url, data=json.dumps(data), params=params, headers=header)
    assert response.ok


def create_user(username):
    session = create_new_session()
    email = "hej@hej.com"
    pwd = "hej"
    data = {"username": username, "email": email, "pwd": pwd}
    url = f"{BASE_URL}/api/register"
    header = {
        "Authorization": "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh",
        "Content-Type": "application/json",
    }
    params = {"latest": 2}
    response = session.post(url, data=json.dumps(data), params=params, headers=header)
    assert response.ok


async def request(command, *args):
    start = time.time()
    result = await command(*args)
    end = time.time()
    logging.info(
        f"Request: {command} | start: {start:<20.6f} | end: {end:<20.6f} | delta: {end - start}"
    )

    # return result, end - start

async def create_many_users(nameList):
    for name in nameList:
        await create_user(name)


async def playwright_visit_user(page, user):
    await press_link(page, "public timeline")
    await navigate_to(page, "/user/" + user)


async def un_follow_users(page: Page, user, follow_string):
    fllw_string = str(follow_string) + " user"
    await press_link(page, fllw_string)
    await press_link(page, "public timeline")


async def login(page, username, password):
    await navigate_to(page, "/login")
    await fill_input(page, "username", username)
    await fill_input(page, "password", password)
    await submit_input(page)


async def logout(page):
    await navigate_to(page, "/logout")


async def create_msg(page, count=30):
    await press_link(page, "my timeline")  # press_link(page, "public timeline")
    for i in range(count):
        await fill_input(page, "text", "tweet")
        await submit_input(page)


async def create_user(page, username):
    await press_link(page, "sign up")
    await fill_input(page, "username", username)
    await fill_input(page, "email", "hej@hej.dk")
    await fill_input(page, "password", "hej")
    await fill_input(page, "password2", "hej")
    await submit_input(page)


async def execute_scenario(nameList, url=BASE_URL, headless=True, slow_mo=None):
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=headless, slow_mo=slow_mo)
        page = await browser.new_page(base_url=url)
        await request(page.goto, "/login")

        # CREATE (register)
        for user in nameList:
            await request(create_user, page, user)

        # CREATE (login), CREATE-READ (30 msg pr user to fill timeline)
        for user in nameList:
            await request(login, page, user, "hej")
            await request(create_msg, page, 30)
            await request(logout, page)

        # # CREATE (login), READ (many different users in loop)
        # for user in nameList:
        #     await request(login, page, user, "hej")
        #     for name in nameList:
        #         await request(playwright_visit_user, page, name)
        #     await request(logout, page)

        # CREATE (login), READ-CREATE-READ (Follow Users)
        for user in nameList:
            await request(login, page, user, "hej")
            for name in nameList:
                if name == user:
                    continue
                await request(playwright_visit_user, page, name)
                await request(un_follow_users, page, name, "Follow")
            await request(logout, page)

        # CREATE (login), READ-DELETE-READ (Unfollow Users)
        for user in nameList:
            await request(login, page, user, "hej")
            for name in nameList:
                if name == user:
                    continue
                await request(playwright_visit_user, page, name)
                await request(un_follow_users, page, name, "Unfollow")
            await request(logout, page)


def setup_logging():
    log_dir = "logs"
    os.makedirs(log_dir, exist_ok=True)

    log_filename = os.path.join(
        log_dir, f"log_{datetime.now().strftime('%Y-%m-%d_%H-%M-%S')}.log"
    )

    logging.basicConfig(
        filename=log_filename,
        level=logging.INFO,
        format="%(asctime)s - %(levelname)s - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )


def makeNameList():
    randomString = "".join(random.choices(string.ascii_uppercase, k=8))
    return [f"User{i+1}{randomString}" for i in range(10)]


async def run(url):
    # Generate a list of 20 usernames
    nameList = makeNameList()

    # Make logfile
    setup_logging()

    # Run scenario
    await execute_scenario(nameList, url)


async def main(url, headless=True, slow_mo=None):
    # Generate a list of 20 usernames
    nameList = makeNameList()


    setup_logging()

    await execute_scenario(nameList, url=url, headless=headless, slow_mo=slow_mo)
