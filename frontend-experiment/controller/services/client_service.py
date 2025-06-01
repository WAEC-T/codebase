import httpx


class ClientService:
    url: str

    def __init__(self, url: str):
        self.url = url

    def temp_print_error(self, exception):
        print("ERROR:")
        print(exception)

    def is_response_success(self, response) -> bool:
        if response.status_code > 299:
            print(response.text)
            return False
        else:
            return True

    async def start_scenario(self, minitwit_url) -> bool:
        async with httpx.AsyncClient() as client:
            try:
                response = await client.post(
                    f"{self.url}/start",
                    content=minitwit_url,
                    timeout=None,
                    headers={"Content-Type": "text/plain"},
                )

                return self.is_response_success(response)
            except Exception as e:
                self.temp_print_error(e)
                return False

    async def clear_db(self, database_string) -> bool: 
        async with httpx.AsyncClient() as client: 
            try:
                response = await client.post(
                    f"{self.url}/cleardb",
                    content=database_string,
                    timeout=None,
                    headers={"Content-Type": "text/plain"},
                )

                return self.is_response_success(response)
            except Exception as e:
                self.temp_print_error(e)
                return False