# MiniTwit Implementations Devcontainer

This container spins up a data on the local network, such that a MiniTwit implementation can be tested in a real-life like scenario

## Pull Container Permissions 

If you get some docker pull permissions issues, you can pull the containers before opening the devcontainer. This is done by pulling the images used in the [compose.yaml](./compose.yaml) file.

```
docker pull ghcr.io/3-1-research-project/postgres:<add current version>
docker pull ghcr.io/3-1-research-project/client:<add current version>
```
