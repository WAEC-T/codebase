# run_experiment.py
import asyncio
import sys
from datetime import datetime
from pathlib import Path

from otii_config import create_otii_app, configure_multimeter
from data_collection import collect_data, generate_output, save_data, save_sequential_time
from orchestration import clean_database, manage_server_docker_service, trigger_clients
from host_sequence.scenario_api import run_api_seq_scenario
from host_sequence.scenario_page import run_page_seq_scenario

# Check the static addresses when preparing the setup!
SERVER_URL = "http://10.7.7.144:5000"
CLIENT_1_URL = "http://10.7.7.199:5001/trigger"
CLIENT_2_URL = "http://10.7.7.178:5001/trigger"
CLIENT_3_URL = "http://10.7.7.145:5001/trigger"

BASE_COMPOSE_FILES_LOCATION = '/media/mmcblk0p2/setup/compose_files/'
COOLDOWN = 30

SERVICES = {#"rust-actix": BASE_COMPOSE_FILES_LOCATION + 'rust-actix-compose-prod.yml'}
            "python-flask": BASE_COMPOSE_FILES_LOCATION + 'python-flask-compose-prod.yml'}

async def main(otii_project, device, out_path, service, run_mode):
    client_urls = [CLIENT_1_URL, CLIENT_2_URL, CLIENT_3_URL]
    time_seq_api_df = None
    time_seq_page_df = None
    print("Clearing DB on remote server...", flush=True)
    reset = clean_database(run_mode == "berries")
    if reset:
        print("DB cleared...", flush=True)
        print(f'Starting experiment - scenario: {run_mode}, service: {service}...', flush=True)
        start_time = datetime.now()
        otii_project.start_recording()
        if run_mode == "berries":
            result_clients_trigger = await trigger_clients(client_urls)
            print(result_clients_trigger, flush=True)
        elif run_mode == "sequential":
            time_seq_api_df = run_api_seq_scenario(service, start_time)
            clean_database(False)
            time_seq_page_df = run_page_seq_scenario(service, start_time)
        otii_project.stop_recording()
        t_delta = datetime.now() - start_time
        print(f"Scenario took {t_delta}", flush=True)
        print(f"Done with scenario {run_mode} for service {service}...", flush=True)
        df, recording_name = collect_data(otii_project, device)
        if time_seq_api_df is not None and time_seq_page_df is not None:
            save_data(time_seq_api_df, recording_name, out_path, run_mode, service)
            save_data(time_seq_page_df, recording_name, out_path, run_mode, service)
        save_data(df, recording_name, out_path, run_mode, service)
        generate_output(otii_project, device)
        

async def main_async(out_path, run_mode="berries"):
    otii_project, device = configure_multimeter(create_otii_app())
    for service, filepath in SERVICES.items():
        service_started = await manage_server_docker_service(
            SERVER_URL.removeprefix("http://").removesuffix(':5000'), 
            filepath
        )
        if service_started:
            await main(otii_project, device, out_path, service, run_mode)

if __name__ == "__main__":
    out_path = Path(sys.argv[1])
    run_mode = str(Path(sys.argv[2]))
    asyncio.run(main_async(out_path, run_mode))
