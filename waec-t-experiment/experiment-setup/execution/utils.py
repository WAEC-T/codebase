import random

TERMINAL_COLORS = ['\033[31m','\033[32m','\033[33m','\033[34m','\033[35m','\033[36m','\033[91m','\033[92m','\033[93m','\033[94m','\033[95m','\033[96m','\033[97m']
RESET = '\033[0m'

def print_random_color(string):
    random_color = random.choice(TERMINAL_COLORS)
    print(f"{random_color}{string}{RESET}", flush=True)

def print_info_call(scenario, service, endpoint, iter_num):
    print(f"Starting {scenario} sequence for service {service} - Call: {endpoint} with {iter_num} requests")