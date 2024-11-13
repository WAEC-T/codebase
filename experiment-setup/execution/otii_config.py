# equipment_config.py

from otii_tcp_client import otii_connection, otii as otii_application
from otii_tcp_client.arc import Channel

def create_otii_app(host="127.0.0.1", port=1905):
    # Connect to the Otii 3 application
    connection = otii_connection.OtiiConnection(host, port)
    connect_response = connection.connect_to_server(try_for_seconds=10)
    if connect_response["type"] == "error":
        raise Exception(
            f'Exit! Error code: {connect_response["errorcode"]}, '
            f'Description: {connect_response["payload"]["message"]}'
        )
    otii_app = otii_application.Otii(connection)

    return otii_app

# TODO: change this to bootstrap a new project if the project doesn't exist or create direct in the equipment?
def configure_multimeter(otii_app):
    # Based on the example from
    # https://github.com/qoitech/otii-tcp-client-python/blob/master/examples/basic_measurement.py
    devices = otii_app.get_devices()
    if len(devices) == 0:
        raise Exception("No Arc or Ace connected!")
    device = devices[0]

    device.enable_channel(Channel.MAIN_CURRENT)
    device.enable_channel(Channel.MAIN_VOLTAGE)
    device.enable_channel(Channel.MAIN_POWER)

    device.set_channel_samplerate(Channel.MAIN_CURRENT, 10000)
    device.set_channel_samplerate(Channel.MAIN_VOLTAGE, 10000)
    device.set_channel_samplerate(Channel.MAIN_POWER, 10000)

    # Get the active project
    project = otii_app.get_active_project()

    return project, device
