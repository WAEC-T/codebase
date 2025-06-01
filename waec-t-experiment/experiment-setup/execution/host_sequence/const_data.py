
def register_data_dummie(i):
    return {
        'username': f'user{i}',
        'email': f'user{i}@waect.com',
        'password': 'waect123',
        'password2': 'waect123'
    }

def login_data_dummie(i):
    return {
        'username': f'user{i}',
        'password': 'waect123'
    }

def api_register_data_dummie(i):
    return {
        'username': f'user{i}',
        'email': f'user{i}@waect.com',
        'pwd': 'waect123',
        'pwd2': 'waect123'
    }

def api_follow_data_dummie(i):
    return {
        'follow': f'user{i}'
    }

def api_unfollow_data_dummie(i):
    return {
        'unfollow': f'user{i}'
    }

message_data = {
    "text": "O melhor guerreiro não é aquele que sempre ganha, mas o que mantém o seu orgulho mesmo na derrota."
}

api_message_data = {
    "content": "O melhor guerreiro não é aquele que sempre ganha, mas o que mantém o seu orgulho mesmo na derrota."
}
