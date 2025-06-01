# itu-minitwit

## C-sharp setup guide

### Other prerequisites

```bashrc
sudo apt update && sudo apt install libpq-dev
```

### Setup and run

- Install specific dotnet with `sudo apt-get install -y dotnet-sdk-8.0`
- Build/Publish `dotnet publish`
- Go into folder `cd src/Minitwit.Web/bin/Release/net8.0/publish`

# How to Run on Raspberry Pi

## Ubuntu Server 24.04.2 LTS

### Default
```bashrc
ASPNETCORE_URLS=http://0.0.0.0:5000 DATABASE_URL=postgresql://user:password@<ip-address>:5432/waect ./Minitwit.Web
```

### Jemalloc
```bashrc
LD_PRELOAD=/usr/lib/<device specific architecture>/libjemalloc.so.2 ASPNETCORE_URLS=http://0.0.0.0:5000 DATABASE_URL=postgresql://user:password@<ip-address>:5432/waect ./Minitwit.Web
```

### QuickJit False
```bashrc
DOTNET_TC_QuickJit=0 ASPNETCORE_URLS=http://0.0.0.0:5000 DATABASE_URL=postgresql://user:password@<ip-address>:5432/waect ./Minitwit.Web
```