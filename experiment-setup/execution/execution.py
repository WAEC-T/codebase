# run_experiment.py
import asyncio
import requests
import sys
from datetime import datetime
from pathlib import Path

from otti_config import create_otii_app, configure_multimeter
from data_collection import collect_data, generate_output, save_data
from orchestration import trigger_clients

# Check the static addresses when preparing the setup!
SERVER_URL = "http://10.0.0.4:5000"
CLIENT_1_URL = "http://10.0.0.3:5001/trigger"
CLIENT_2_URL = "http://10.0.0.2:5001/trigger"
CLIENT_3_URL = "http://10.0.0.5:5001/trigger"

async def main(otii_project, device, out_path):
    print("Clearing DB on remote server...", flush=True)
    r = requests.get(f"{SERVER_URL}/cleardb")

    if r.ok:
        print("Starting scenario on three clients...", flush=True)
        start_time = datetime.now()
        client_urls = [CLIENT_1_URL, CLIENT_2_URL, CLIENT_3_URL]

        otii_project.start_recording()

        results = trigger_clients()
        # TODO: add here functionality to not only print but to check if it is ok! Otherwise try to trigger it again.
        print(results, flush=True)

        otii_project.stop_recording()

        t_delta = datetime.now() - start_time
        print(f"Scenario took {t_delta}", flush=True)
        print("Done with scenario...", flush=True)

        # Collect and save data
        df, recording_name = collect_data(otii_project, device)
        save_data(df, recording_name, out_path)
        generate_output(otii_project, device)

if __name__ == "__main__":
    out_path = Path(sys.argv[1])
    otii_project, device = configure_multimeter(create_otii_app())
    for _ in range(10):
        asyncio.run(main(otii_project, device, out_path))
