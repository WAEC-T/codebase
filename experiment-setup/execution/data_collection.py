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

def save_data(df, recording_name, out_path):
    df.to_csv(Path(out_path, f"{recording_name}.csv"))
