# Testing scripts using containers

In the bootstrap folder run the following commands to build and execute.

## Client

```sh
docker build -f "./bootstrap/test/client/Dockerfile" -t "test-client-provision" .

docker run -it test-client-provision /bin/sh
```

Verify if the python is installed as well as the scenario script and the venv folder are in the work directory /home/waect!

## Server

This particular test due to the complexity of installing docker inside a docker container (funny) uses a dind alpine base image (docker in docker). Therefore, it is possible to execute the whole script under it but in fact doesnt test the docker start in normal conditions. This should be done during the set up. However, all the other things like alpine semantics and packages can be explicity checked.

```sh
docker build -f "./bootstrap/test/server/Dockerfile" -t "test-server-provision" .

#Start the container in detached mode:
docker run --privileged -d --rm --name test-server-provision test-server-provision

#Enter in interactive mode using sh
docker exec -it test-server-provision /bin/sh

#Inside the container execute the provisioning bash script
./server-provisioning.sh
```

## Future usefull commands

Stop docker command:

```sh
service docker stop
```

Remove docker to automatically start in the os:

```sh
rc-update del docker default
```

List all services running:

```sh
rc-status
```

Show all startup services

```sh
rc-update -v show
```
