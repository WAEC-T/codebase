from otii_tcp_client import otii_client

def create_otii_app():
    client = otii_client.OtiiClient().connect()
    return client

def configure_multimeter(otii_app):
    devices = otii_app.get_devices()
    if len(devices) == 0:
        raise Exception("No Arc or Ace connected!")
    device = devices[0]

    device.enable_channel('mp', True)
    device.enable_channel('mc', True)
    device.enable_channel('mv', True)

    device.set_main_voltage(5.0)
    device.set_exp_voltage(4.9)
    device.set_max_current(2.5)

    project = otii_app.get_active_project()
    return project, device
