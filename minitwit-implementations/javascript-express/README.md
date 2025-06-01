# itu-minitwit

## Javascript Express setup guide

### Other prerequisites

```bashrc
sudo apt update && sudo apt install libpq-dev
```

Make sure to have node package manager

```bashrc
sudo apt install npm
```

### Setup and run

- Install required packages with `npm install`
- Start the app with `npm start`

# How to Run on Raspberry Pi

## Ubuntu Server 24.04.2 LTS

### Default

```bashrc
SECRET_KEY="waectsk" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" npm run start-node 
```

### Jemalloc

```bashrc
LD_PRELOAD=/usr/lib/<device specific architecture>/libjemalloc.so.2 SECRET_KEY="waectsk" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" npm run start-node 
```

### Jitless

```bashrc
SECRET_KEY="waectsk" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" npm run start-node --jitless
```