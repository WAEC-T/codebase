import csv
import datetime
from otii_tcp_client import otii_connection, otii as otii_application
from otii_tcp_client import arc as otii_arc


class OtiiService:
    otii_app: any

    def __init__(self, host: str = "127.0.0.1", port: int = 1905):
        connection = otii_connection.OtiiConnection(host, port)
        connect_response = connection.connect_to_server(try_for_seconds=10)
        if connect_response["type"] == "error":
            raise Exception(
                f'Exit! Error code: {connect_response["errorcode"]}, '
                f'Description: {connect_response["payload"]["message"]}'
            )
        self.otii_app = otii_application.Otii(connection)

    def configure_multimeter(self) -> tuple[any, any]:
        # Based on the example from
        # https://github.com/qoitech/otii-tcp-client-python/blob/master/examples/basic_measurement.py
        devices = self.otii_app.get_devices()
        if len(devices) == 0:
            raise Exception("No Arc or Ace connected!")
        device = devices[0]

        device.enable_channel("mc", True)
        device.enable_channel("mv", True)
        device.enable_channel("mp", True)

        device.set_channel_samplerate("mc", 10000)
        device.set_channel_samplerate("mv", 10000)
        device.set_channel_samplerate("mp", 10000)

        project = self.otii_app.get_active_project()

        return project, device

    def collect_data(
        self,
        otii_project: otii_application.project.Project,
        device: otii_arc.Arc,
        schema_path: str,
        iteration: int,
    ):

        recording = otii_project.get_last_recording()

        current_count = recording.get_channel_data_count(device.id, "mc")

        timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
        schema = schema_path.split("/")[-1].split(".")[0]
        filename = f"{schema}_{iteration}_{timestamp}.csv"

        with open(filename, mode="w") as data_file:

            data_writer = csv.writer(
                data_file, delimiter=",", quotechar='"', quoting=csv.QUOTE_MINIMAL
            )
            data_writer.writerow(["seconds", "current", "voltage", "power"])

            print("Writing data to CSV...")

            index = 0
            while index < current_count:
                current_data = recording.get_channel_data(device.id, "mc", index, 1000)[
                    "values"
                ]
                voltage_data = recording.get_channel_data(device.id, "mv", index, 1000)[
                    "values"
                ]
                power_data = recording.get_channel_data(device.id, "mp", index, 1000)[
                    "values"
                ]

                for i in range(len(current_data)):
                    data_writer.writerow(
                        [
                            (index + i) / 10000,
                            current_data[i],
                            voltage_data[i],
                            power_data[i],
                        ]
                    )

                index += len(current_data)

        print(f"Data saved to {filename}!")