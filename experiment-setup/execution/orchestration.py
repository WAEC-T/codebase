import asyncio
import httpx

CLIENT_1_URL = "http://10.0.0.3:5001/trigger"
CLIENT_2_URL = "http://10.0.0.2:5001/trigger"
CLIENT_3_URL = "http://10.0.0.5:5001/trigger"
SERVER_URL = "http://10.0.0.4:5000"

async def get_async(url):
    print(f"Starting scenario on {url}...", flush=True)
    timeout = httpx.Timeout(10.0, read=None)
    async with httpx.AsyncClient() as client:
        return await client.get(url, timeout=timeout)

async def trigger_clients():
    client_urls = [CLIENT_1_URL, CLIENT_2_URL, CLIENT_3_URL]
    results = await asyncio.gather(*map(get_async, client_urls))
    return results

def clear_server_db():
    import requests
    print("Clearing DB on remote server...", flush=True)
    r = requests.get(f"{SERVER_URL}/cleardb")
    if not r.ok:
        raise Exception("Failed to clear server DB")

# TODO: Implement this function and add to the workflow.
def server_running_verification():
    return True