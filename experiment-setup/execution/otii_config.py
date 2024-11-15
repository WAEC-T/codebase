# equipment_config.py

from otii_tcp_client import otii_client

def create_otii_app():
    # Connect to the Otii 3 application
    client = otii_client.OtiiClient().connect()
    return client

# TODO: change this to bootstrap a new project if the project doesn't exist or create direct in the equipment?
def configure_multimeter(otii_app):
    # Based on the example from
    # https://github.com/qoitech/otii-tcp-client-python/blob/master/examples/basic_measurement.py
    devices = otii_app.get_devices()
    if len(devices) == 0:
        raise Exception("No Arc or Ace connected!")
    device = devices[0]

    device.enable_channel('mp', True)

    # Set device parameters
    device.set_main_voltage(5.0)  # Set main voltage to 5V
    device.set_exp_voltage(4.9)   # Set expansion voltage to 4.9V
    device.set_max_current(2.0)   # Set max current to 2Ax

    # Get the active project
    project = otii_app.get_active_project()
    return project, device
