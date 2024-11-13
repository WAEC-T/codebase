from pathlib import Path

def collect_data(otii_project, device):
    """Collect data from the most recent recording."""
    project = otii_project

    # Get statistics for the recording
    recording = project.get_last_recording()
    info = recording.get_channel_info(device.id, 'mc')
    statistics = recording.get_channel_statistics(device.id, 'mc', info['from'], info['to'])

    # Print the statistics
    print(f'From:        {info["from"]} s')
    print(f'To:          {info["to"]} s')
    print(f'Offset:      {info["offset"]} s')
    print(f'Sample rate: {info["sample_rate"]}')

    print(f'Min:         {statistics["min"]:.5} A')
    print(f'Max:         {statistics["max"]:.5} A')
    print(f'Average:     {statistics["average"]:.5} A')
    print(f'Energy:      {statistics["energy"] / 3600:.5} Wh')



def generate_output(otii_project, device):
    """Generate and print statistics for the main channels in the last recording."""
    recording = otii_project.get_last_recording()
    if not recording:
        raise Exception("No recording found.")

    # Print statistics for main current
    stats = recording.get_channel_statistics(device, "mc")
    print(f"Main Current (mc): Min={stats['min']}, Max={stats['max']}, Avg={stats['avg']}, Energy={stats['energy']}", flush=True)

    # Print statistics for main voltage and main power
    for channel in ["mv", "mp"]:
        stats = recording.get_channel_statistics(device, channel)
        print(f"{channel}: Min={stats['min']}, Max={stats['max']}, Avg={stats['avg']}", flush=True)

def save_sequential_time(dataframe_api, dataframe_page, recording_name, out_path):
    """Save sequential time data for API and Page as JSON files."""
    api_path = Path(out_path, f"sequential_time_api_{recording_name}.json")
    page_path = Path(out_path, f"sequential_time_page_{recording_name}.json")

    dataframe_api.to_json(api_path, orient="records", lines=True)
    dataframe_page.to_json(page_path, orient="records", lines=True)

def save_data(dataframe, recording_name, out_path):
    """Save the collected data as a CSV file."""
    csv_path = Path(out_path, f"{recording_name}.csv")
    dataframe.to_csv(csv_path)