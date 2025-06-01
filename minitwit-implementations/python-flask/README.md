# itu-minitwit

## Python setup guide

### Setup and run

Make venv
`python3 -m venv venv`

Go into Venv
`source venv/bin/activate`

Install requirements (if not done)
`pip3 install -r requirements.txt`

# How to Run on Raspberry Pi

## Ubuntu Server 24.04.2 LTS

### Default
```bashrc
SECRET_KEY="waectsk" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" python -m gunicorn --bind 0.0.0.0:5000 --workers 4 wsgi:app
```

### Jemalloc
```bashrc
LD_PRELOAD=/usr/lib/<device specific architecture>/libjemalloc.so.2 SECRET_KEY="waectsk" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" python -m gunicorn --bind 0.0.0.0:5000 --workers 4 wsgi:app
```