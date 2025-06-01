import base64
import requests

USERNAME = 'simulator'
PWD = 'super_safe!'
CREDENTIALS = ':'.join([USERNAME, PWD]).encode('ascii')
ENCODED_CREDENTIALS = base64.b64encode(CREDENTIALS).decode()
HEADERS = {'Connection': 'close',
           'Content-Type': 'application/json',
           f'Authorization': f'Basic {ENCODED_CREDENTIALS}'}

def create_new_session():
    session = requests.Session()
    session.headers.update({
        'Connection': 'close',
        'Content-Type': 'application/json',
        'Authorization': f'Basic {ENCODED_CREDENTIALS}'
    })
    return session
print(ENCODED_CREDENTIALS)