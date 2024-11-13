from pathlib import Path
from otii_tcp_client.arc import Channel

def collect_data(otii_project, device):
    recording = otii_project.get_last_recording()
    df = recording.get_dataframe(device, (Channel.MAIN_CURRENT, Channel.MAIN_VOLTAGE, Channel.MAIN_POWER))
    return df, recording.name

def generate_output(otii_project, device):
    recording = otii_project.get_last_recording()
    minimum, maximum, avg, energy = recording.get_complete_channel_statistics(device, Channel.MAIN_CURRENT)
    print(f"{Channel.MAIN_CURRENT.name}: {minimum}, {maximum}, {avg}, {energy}", flush=True)

    for channel in (Channel.MAIN_VOLTAGE, Channel.MAIN_POWER):
        minimum, maximum, avg = recording.get_complete_channel_statistics(device, channel)
        print(f"{channel.name}: {minimum}, {maximum}, {avg}", flush=True)

def save_sequential_time(dataframe_api , dataframe_page, recording_name, out_path):
    dataframe_api.to_json(Path(out_path,f"sequential_time_api_{recording_name}.json"), orient="records", lines=True)
    dataframe_page.to_json(Path(out_path,f"sequential_time_page_{recording_name}.json"), orient="records", lines=True)

def save_data(dataframe, recording_name, out_path):
    dataframe.to_csv(Path(out_path, f"{recording_name}.csv"))
