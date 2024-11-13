# run_experiment.py
import asyncio
import sys
from datetime import datetime
from pathlib import Path

from otii_config import create_otii_app, configure_multimeter
from data_collection import collect_data, generate_output, save_data, save_sequential_time
from orchestration import clean_database, manage_server_docker_service, trigger_clients
#from client.host_sequence.scenario_api import run_api_seq_scenario
#from client.host_sequence.scenario_page import run_page_seq_scenario


# Check the static addresses when preparing the setup!
SERVER_URL = "http://10.7.7.144:5000"
CLIENT_1_URL = "http://10.7.7.177:5001/trigger"
# CLIENT_2_URL = "http://10.0.0.2:5001/trigger"
# CLIENT_3_URL = "http://10.0.0.5:5001/trigger"

SERVICES = ["python-flask"]

async def main(otii_project, device, out_path, service, run_mode="standard"):
    print("Clearing DB on remote server...", flush=True)
    
    reset = clean_database()

    if reset:
        print("Starting scenario on three clients...", flush=True)
        start_time = datetime.now()

        otii_project.start_recording()
        if run_mode == "standard":
            client_urls = [CLIENT_1_URL]
            result_clients_trigger = await trigger_clients(client_urls)
            print(result_clients_trigger, flush=True)
        # elif run_mode == "sequential":
        #     time_seq_api_df = await asyncio.create_task(run_api_seq_scenario(service, start_time))
        #     time_seq_page_df = await asyncio.create_task(run_page_seq_scenario(service, start_time))

        otii_project.stop_recording()
        t_delta = datetime.now() - start_time
        print(f"Scenario took {t_delta}", flush=True)
        print(f"Done with scenario {run_mode}...", flush=True)

        # Collect and save data

        df, recording_name = collect_data(otii_project, device)
        return
        # save_data(df, recording_name, out_path)
        # if time_seq_api_df and time_seq_page_df:
        #     save_sequential_time(time_seq_api_df, time_seq_page_df, recording_name, out_path)
        # generate_output(otii_project, device)

if __name__ == "__main__":
    out_path = Path(sys.argv[1])
    otii_project, device = configure_multimeter(create_otii_app())
    print("otii_project: ", otii_project)
    print("device: ", device)
    for service in SERVICES:
        service_started = manage_server_docker_service(SERVER_URL.removeprefix("http://"), service);
        if service_started:
                asyncio.run(main(otii_project, device, out_path, service))
