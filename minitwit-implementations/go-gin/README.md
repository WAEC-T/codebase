# itu-minitwit
## Go setup guide
Uses Go 1.22.2
### Other prerequisites
```bashrc
ENV PATH="/usr/local/go/bin:$PATH"
```

For installing 1.22.2 specifically:
```bashrc
RUN wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz \
    && rm go1.22.2.linux-amd64.tar.gz
```
### Setup and run
- Start the app with `go build src/main.go `
# How to Run on Raspberry Pi
## Ubuntu Server 24.04.2 LTS
### Requirements
- Go version = 1.22.2

### How to Run Go with different configurations
#### Default
```
go build src/main.go 
```
```
./main 
```
#### PGO
PGO is a bit complex as it needs to be trained first, but the training should be enables on port 6060 of this go project

According to the docs, having a valid PGO file in the project with automatically enable PGO when building
See (https://go.dev/doc/pgo#building) for more

#### jemalloc
After building go
```
LD_PRELOAD=/usr/lib/<device specific architecture>/libjemalloc.so.2 ./main
```