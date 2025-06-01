from datetime import datetime
from pathlib import Path
from utils import print_random_color
import pandas as pd
import pytz

RESULTS_SECTION = '\033[32m-------------------------------\n      EXPERIMENT INFO:\n-------------------------------\033[0m\n'

def collect_data(otii_project, device, iteration=None):
    """Collect data from the most recent recording."""
    project = otii_project
    recording = project.get_last_recording()

    samples_amount = recording.get_channel_data_count(device.id, 'mp')
    power_records = recording.get_channel_data(device.id, 'mp', 0, samples_amount)
    current_records = recording.get_channel_data(device.id, 'mc', 0, samples_amount)
    voltage_records = recording.get_channel_data(device.id, 'mv', 0, samples_amount)

    dataframe = create_experiment_dataframe(power_records, current_records, voltage_records, timestamp_from_start_test(recording.name), iteration)
    return dataframe, recording.name

def timestamp_from_start_test(experiment_start_datetime):
    dt = datetime.strptime(experiment_start_datetime, "%Y-%m-%d %H:%M:%S")
    local_zone = pytz.timezone('Europe/Copenhagen')
    dt_local = local_zone.localize(dt)
    dt_utc = dt_local.astimezone(pytz.UTC)
    return dt_utc.timestamp()

def create_experiment_dataframe(power_records, current_records, voltage_records, start_timestamp, iteration):
    power_df = handle_timestamp(pd.DataFrame(power_records).rename(columns={'values': 'power'}), start_timestamp)
    current_df = handle_timestamp(pd.DataFrame(current_records).rename(columns={'values': 'current'}), start_timestamp)
    voltage_df = handle_timestamp(pd.DataFrame(voltage_records).rename(columns={'values': 'voltage'}), start_timestamp)

    merged_df = pd.merge(power_df[['timestamp', 'power']], current_df[['timestamp', 'current']], on='timestamp', how='outer')
    merged_df = pd.merge(merged_df, voltage_df[['timestamp', 'voltage']], on='timestamp', how='outer')

    if iteration is not None:
        merged_df['iteration'] = iteration

    return merged_df

def handle_timestamp(dataframe, start_timestamp):
    dataframe['timestamp'] = dataframe['interval'].cumsum() + start_timestamp
    return dataframe

def generate_output(otii_project, device):
    """Generate and print statistics for the main channels in the last recording."""
    recording = otii_project.get_last_recording()
    if not recording:
        raise Exception("No recording found.")

    info = recording.get_channel_info(device.id, 'mp')
    print(f'{RESULTS_SECTION}Start-Stop: {info["from"]} s - {info["to"]} s, Offset: {info["offset"]} s, Sample rate: {info["sample_rate"]}')

    for channel, variable in {"mc": "Current (A)", "mv": "Voltage (V)", "mp": "Power (W):"}.items():
        stats = recording.get_channel_statistics(device.id, channel, info['from'], info['to'])
        print_random_color(f"{variable}: Min={stats['min']}, Max={stats['max']}, Avg={stats['average']}")

def save_sequential_time(dataframe_api, dataframe_page, out_path, service_name, run_mode, recording_name):
    """Save sequential time data for API and Page as JSON files."""
    service_path = Path(out_path, service_name)
    target_folder = service_path / run_mode
    target_folder.mkdir(parents=True, exist_ok=True)

    api_path = target_folder / f"{service_name}-sequential-api-{recording_name}.csv"
    page_path = target_folder / f"{service_name}-sequential-page-{recording_name}.csv"
    
    dataframe_api.to_csv(api_path, index=False)
    dataframe_page.to_csv(page_path, index=False)

def save_data(dataframe, out_path, service_name, run_mode, recording_name):
    """Save the collected data as a CSV file."""
    service_path = Path(out_path, service_name)
    target_folder = service_path / run_mode
    target_folder.mkdir(parents=True, exist_ok=True)
    csv_path = target_folder / f"{service_name}-power-{recording_name}.csv"
    dataframe.to_csv(csv_path, index=False)