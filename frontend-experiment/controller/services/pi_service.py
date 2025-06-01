import paramiko


class PiService:
    ip: str
    username: str
    password: str

    def __init__(self, ip: str, username: str, password: str):
        self.ip = ip
        self.username = username
        self.password = password

    def get_raspberry_pi_temp(self):
        command = "cat /sys/class/thermal/thermal_zone0/temp"  

        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(self.ip, username=self.username, password=self.password)

        _, stdout, stderr = ssh.exec_command(command)
        temp_output = stdout.read().decode().strip()

        ssh.close()

        if temp_output.isdigit():
            cpu_temp = int(temp_output) / 1000.0
            return cpu_temp
        else:
            raise ValueError(f"Failed to retrieve valid temperature data. stderror: {stderr.read().decode().strip()}")
