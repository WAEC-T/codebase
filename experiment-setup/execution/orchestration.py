import asyncio
import asyncssh
import httpx
import psycopg2
from dotenv import load_dotenv
import os

load_dotenv("../../../.env.prod")
DATABASE_URL = os.environ.get("DATABASE_URL")
SSH_USER=os.environ.get("SSH_USER")
SSH_PASS=os.environ.get("SSH_PASS")

async def get_async(url):
    print(f"Starting scenario on {url}...", flush=True)
    timeout = httpx.Timeout(10.0, read=None)
    async with httpx.AsyncClient() as client:
        return await client.get(url, timeout=timeout)

async def trigger_clients(clients: list[str]):
    results = await asyncio.gather(*map(get_async, clients))
    return results

def clean_database():
    try:
        with psycopg2.connect(DATABASE_URL) as conn:
            with conn.cursor() as cur:
                cur.execute("TRUNCATE TABLE users CASCADE;")
                cur.execute("TRUNCATE TABLE messages CASCADE;")
                cur.execute("TRUNCATE TABLE followers CASCADE;")
                conn.commit()
        return True
    except Exception as e:
        print(f"Database cleaning failed: {e}")
        return False 

async def manage_server_docker_service(host, docker_compose_dir):
    try:
        async with await asyncssh.connect(host, username=SSH_USER, password=SSH_PASS) as conn:

            _ = await conn.run("docker-compose stop", check=True)
            print("Container stopped.", flush=True)

            _ = await conn.run(f"cd {docker_compose_dir} && docker-compose up -d", check=True)
            print(f"Docker service {docker_compose_dir} started successfully...", flush=True)

            return True 
    except Exception as e:
        print(f"Failed to manage Docker on remote server: {e}", flush=True)
        return False

