# Client Code

You can test the client code locally by setting up a postgres database and MiniTwit implementation, e.g.,:

Postgres container: `docker run --rm --detach --env POSTGRES_HOST_AUTH_METHOD=trust --network host --name postgres-temp ghcr.io/3-1-research-project/postgres:17`

MiniTwit Python: 
```
$env:SECRET_KEY="waectsk"
$env:DATABASE_URL="postgresql://user:pass@localhost:5432/waect"

cd /python-flask/
gunicorn --bind 0.0.0.0:5000 wsgi:app
```

Then the client code can be run in headed mode by in slow motion by:

```
cd /client/

python run_scenario.py http://localhost:5000 False 200
```

## Container

A container for running the clients are pushed to the [GitHub Organization 3-1 Container Registry (called packages)](https://github.com/orgs/3-1-research-project/packages), and it is called `3-1-research-project/client:<insert version number>`

To run a container use the command below and update the port used. Make sure the `<port>` is the same. Also remember to change `<client-x>`, and `<version number>`

```bash
docker run --detach --name <client-x> --env PORT=<port> --publish <port>:<port> 3-1-research-project/client:<version number>
```

### How to limit the amount of CPUs used

If you want to limit the container to use 1 specific cpu, you can specify it as below

```bash
docker run --detach --name <client-x> --env PORT=<port> --publish <port>:<port> --cpuset-cpus="<core number>" 3-1-research-project/client:<version number>
```

Where `<core number>` is the specific core you want to limit to. Note, it is a 0 indexed list meaning if a CPU has 4 cores the possible values are 0, 1, 2, and 3.
