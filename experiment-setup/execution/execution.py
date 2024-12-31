# run_experiment.py
import asyncio
import sys
from datetime import datetime
from pathlib import Path

import pandas

from otii_config import create_otii_app, configure_multimeter
from data_collection import collect_data, generate_output, save_data, save_sequential_time
from orchestration import clean_database, manage_server_docker_service, trigger_clients
from host_sequence.scenario_api import run_api_seq_scenario
from host_sequence.scenario_page import run_page_seq_scenario

# Check the static addresses when preparing the setup!
SERVER_URL = "YOUR_IP_HERE:PORT/trigger"
CLIENT_1_URL = "YOUR_IP_HERE:PORT/trigger"
CLIENT_2_URL = "YOUR_IP_HERE:PORT/trigger"
CLIENT_3_URL = "YOUR_IP_HERE:PORT/trigger"

BASE_COMPOSE_FILES_LOCATION = '/media/mmcblk0p2/setup/compose_files/'
COOLDOWN = 30

SERVICES = {"rust-actix": BASE_COMPOSE_FILES_LOCATION + 'rust-actix-compose-prod.yml',
            "python-flask": BASE_COMPOSE_FILES_LOCATION + 'python-flask-compose-prod.yml',
            "go-gorilla": BASE_COMPOSE_FILES_LOCATION + 'go-gorilla-compose-prod.yml'
            }


async def execute_experiment(otii_project, device, out_path, service, run_mode, iterations):
    client_urls = [CLIENT_1_URL, CLIENT_2_URL, CLIENT_3_URL]
    experiment_power_data = pandas.DataFrame()

    print(f'Starting experiment - scenario: {run_mode}, service: {service}...', flush=True)

    recording_name = f"batch_{iterations}_start_{datetime.now().timestamp()}"

    if(run_mode == "berries" and iterations > 0):
        for iteration in range(iterations):
            experiment_data_iteration = await execute_berries(client_urls, otii_project, device, service, iteration)
            experiment_power_data = pandas.concat([experiment_power_data, experiment_data_iteration], ignore_index=True)

    elif(run_mode == "sequential"):
        experiment_power_data, name = execute_sequential(otii_project, device, service)
        recording_name = name
    else:
        raise ValueError(f"Invalid run mode: {run_mode}")
    save_data(experiment_power_data, out_path, service, run_mode, recording_name)
    generate_output(otii_project, device)


def execute_sequential(otii_project, device, service):
    time_seq_api_df = None
    time_seq_page_df = None
    if prepare_database(False):
        print(f'Executing ~ sequential ~ scenario for service: {service}...', flush=True)
        start_time = datetime.now()
        otii_project.start_recording()
        time_seq_api_df = run_api_seq_scenario(service, start_time)
        clean_database(False)
        time_seq_page_df = run_page_seq_scenario(service, start_time)
        otii_project.stop_recording()
        t_delta = datetime.now() - start_time
        print(f"Sequential scenario for service {service} took ~ {t_delta} seconds ~", flush=True)
        dataframe, recording_name = collect_data(otii_project, device)
        if time_seq_api_df is not None and time_seq_page_df is not None:
            save_sequential_time(time_seq_api_df, time_seq_page_df, out_path, service, run_mode, recording_name)
        return dataframe, recording_name


async def execute_berries(client_urls, otii_project, device, service, iteration):
    if prepare_database(True):
        print(f'Executing ~ berries ~ scenario for serviceservice: {service}...', flush=True)
        start_time = datetime.now()
        otii_project.start_recording()
        if run_mode == "berries":
            result_clients_trigger = await trigger_clients(client_urls)
            print(result_clients_trigger, flush=True)
        otii_project.stop_recording()
        t_delta = datetime.now() - start_time
        print(f"Berries scenario for service {service} iteration {iteration} took ~ {t_delta} seconds ~", flush=True)
        dataframe, _ = collect_data(otii_project, device, iteration)
        return dataframe


def prepare_database(dump_data):
    print("Clearing database on remote server... with dump data? {dump_data}", flush=True)
    reset = clean_database(dump_data)
    if reset: print("Database procedure was done successfully ...", flush=True)
    else: print("Database procedure error !")
    return reset


async def main(run_mode, out_path, iterations):
    ssh_target = SERVER_URL.removeprefix("http://").removesuffix(':5000')
    otii_project, device = configure_multimeter(create_otii_app())
    for service, filepath in SERVICES.items():
        service_started = await manage_server_docker_service(
            ssh_target, 
            filepath
        )
        if service_started:
            await execute_experiment(otii_project, device, out_path, service, run_mode, iterations)
    await manage_server_docker_service(ssh_target, "", True)


if __name__ == "__main__":
    run_mode = str(Path(sys.argv[1])) if len(sys.argv) > 1 else "berries"
    out_path = Path(sys.argv[2]) if len(sys.argv) > 2 else Path("data/out")
    asyncio.run(main(run_mode, out_path, 10))
