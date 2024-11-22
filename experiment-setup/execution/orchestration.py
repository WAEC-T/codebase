import asyncio
import asyncssh
import httpx
import psycopg2
from dotenv import load_dotenv
import os

load_dotenv("../../.env.prod")
DATABASE_URL = os.environ.get("DATABASE_URL")
SSH_USER=os.environ.get("SSH_USER")
SSH_PASS=os.environ.get("SSH_PASS")

async def get_async(url):
    print(f"Starting scenario on {url}...", flush=True)
    timeout = httpx.Timeout(10.0, read=None)
    async with httpx.AsyncClient() as client:
        try:
            response = await client.get(url, timeout=timeout)
            return response
        except Exception as e:
            print(f"Error with {url}: {e}", flush=True)
        finally:
            print(f"Finished scenario on {url}.", flush=True)

async def trigger_clients(clients: list[str]):
    results = await asyncio.gather(*map(get_async, clients))
    print("All clients finished execution...", flush=True)
    return results

def clean_database(with_dump: bool):
    try:
        with psycopg2.connect(DATABASE_URL) as conn:
            with conn.cursor() as cur:
                cur.execute("TRUNCATE TABLE users CASCADE;")
                cur.execute("TRUNCATE TABLE messages CASCADE;")
                cur.execute("TRUNCATE TABLE followers CASCADE;")
                if(with_dump):
                    with open('dump.sql', 'r') as f:
                        sql_commands = f.read()
                        cur.execute(sql_commands)
                    conn.commit()
        return True
    except Exception as e:
        print(f"Database cleaning failed: {e}")
        return False 

async def manage_server_docker_service(host, docker_compose_file, last_service=False):
    try:
        async with await asyncssh.connect(host, username=SSH_USER, password=SSH_PASS, known_hosts=None) as conn:
            stop_cmd = "docker ps -a -q | xargs -r docker stop"
            up_cmd = f"docker-compose -f {docker_compose_file} up -d"

            result = await conn.run(stop_cmd, check=False)

            if result.returncode == 0: 
                print("Containers stopped successfully.", flush=True)
            else:
                print(f"Failed to stop containers. Error: {result.stderr}", flush=True)

            if (last_service): 
                return True
            
            result = await conn.run(up_cmd, check=False)

            if result.returncode == 0:
                print(f"Docker service {docker_compose_file} started successfully...", flush=True)
            else:
                print(f"Failed to start Docker Compose. Error: {result.stderr}", flush=True)

            return True 
    except Exception as e:
        print(f"Failed to manage Docker on remote server: {e}", flush=True)
        return False