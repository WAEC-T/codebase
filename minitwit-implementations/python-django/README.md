# itu-minitwit

## Python setup guide

### Setup and run

If not using the dev container, it is recommended to use a virtual environment:

Make venv
`python3 -m venv venv`

Go into Venv
`source venv/bin/activate`

Regardless of using dev container or virtual environment, the requirements for the implementation still needs to be installed: 

Install requirements (if not done)
`pip3 install -r requirements.txt`

Then navigate to the src folder and run the following: 

```bash
python3 ./manage_prod.py migrate --run-syncdb
python3 ./manage_prod.py collectstatic
python3 ./manage_prod.py runserver 0.0.0.0:5000 --insecure
```
Note: currently the static files (css) is only served correctly in production using the --insecure flag. While this might not be recommended in a real production setting, we found it better than to serve the static files via a proxy of sort, as Django does not support serving it otherwise. Another alternative to serve the static files within Django is to set DEBUG=True in the prod_settings, but that is not recommended either. 

### Access the App Locally

Once the server is running, visit:
[http://localhost:5000](http://localhost:5000)

If localhost does not work try 127.0.0.1:5000

# How to Run on Raspberry Pi

## Ubuntu Server 24.04.2 LTS

### Default
```bashrc
SECRET_KEY="waect" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" python3 ./manage_prod.py runserver 0.0.0.0:5000 --workers 4 
```

### Jemalloc
```bashrc
LD_PRELOAD=/usr/lib/<device specific architecture>/libjemalloc.so.2 SECRET_KEY="waect" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" python3 ./manage_prod.py runserver 0.0.0.0:5000 --workers 4 
```
