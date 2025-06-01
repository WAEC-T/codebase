import argparse
import time
from services.client_service import ClientService
from services.otii_service import OtiiService
from services.pi_service import PiService
import asyncio
import json
import datetime

USERNAME_SERVER_PI = "admin"
PASSWORD_SERVER_PI = "admin"


async def run(
    client_urls: str,
    minitwit_url: str,
    database_connection_string: str,
    num_cores: str,
    data_name: str = "default",
    use_otii: bool = False,
    username_ssh_server: str = None,
    password_ssh_server: str = None,
):
    if use_otii:
        otii_service = OtiiService()
        otii_project, device = otii_service.configure_multimeter()

    print("urls", client_urls)

    new_list = []
    for url in client_urls:
        url = url.strip()
        for i in range(int(num_cores)):

            list_chars = list(url)
            list_chars[-1] = str(int(list_chars[-1]) + i)
            new_url = "".join(list_chars)
            new_list.append(new_url)

    print(new_list)

    client_services = [ClientService(url=url) for url in new_list]

    server_ip = minitwit_url.split("//")[1].split(":")[0]

    pi_service = PiService(
        ip=server_ip,
        username=(
            username_ssh_server
            if username_ssh_server is not None
            else USERNAME_SERVER_PI
        ),
        password=(
            password_ssh_server
            if password_ssh_server is not None
            else PASSWORD_SERVER_PI
        ),
    )

    print("Starting Experiment")
    print(f"For data_name: {data_name}")

    for i in range(5):
        print(f"Collecting temperature data pre experiment {i}")
        temperature_data_pre_experiment = pi_service.get_raspberry_pi_temp()
        print("Pre experiment temp: " + str(temperature_data_pre_experiment))

        if use_otii:
            print(f"Starting recording {i}")
            otii_project.start_recording()

        print(f"Starting scenario {i}")

        results = await asyncio.gather(
            *map(
                lambda client_service: start_scenario_for_client(
                    client_service, minitwit_url
                ),
                client_services,
            )
        )

        print(results)
        if False in results:
            print("One or more scenario runs failed, stopping measurement")
            if use_otii:
                otii_project.stop_recording()
            return

        if use_otii:
            otii_project.stop_recording()
        print(f"Done recording {i}")

        print(f"Collecting temperature data post experiment {i}")
        timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
        with open(f"{data_name}_temperature_{i}_{timestamp}.csv", "w") as file:
            file.write(
                f"CPU Temperature Pre experiment: {temperature_data_pre_experiment}°C\n"
            )
            file.write(
                f"CPU Temperature Post experiment: {pi_service.get_raspberry_pi_temp()}°C\n"
            )

        print("Clearing Database")
        db_cleared = await client_services[0].clear_db(database_connection_string)
        if db_cleared:
            print("DB Cleared successfully")
        else:
            print("Error Clearing the DB")

        otii_service.collect_data(otii_project, device, data_name, i)


async def start_scenario_for_client(client_service: ClientService, minitwit_url: str):
    return await client_service.start_scenario(minitwit_url=minitwit_url)


def get_json_data(data_name: str):
    with open(data_name, "r") as file:
        return json.load(file)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Run an energy measurement test")
    parser.add_argument(
        "client_urls",
        help="Space separated list of urls of the client applications",
        type=lambda arg: arg.split(","),
    )
    parser.add_argument("minitwit_url", help="Url of the MiniTwit application")
    parser.add_argument(
        "database_connection_string", help="Connection string of the MiniTwit database"
    )
    parser.add_argument(
        "num_cores", help="Number of cores to use for the test in each client"
    )
    parser.add_argument(
        "use_otii",
        action="store_true",
        help="Flag to indicate whether to use Otii for measurements",
    )
    args = parser.parse_args()
    asyncio.run(
        run(
            args.client_urls,
            args.minitwit_url,
            args.database_connection_string,
            args.num_cores,
            args.use_otii,
        )
    )
