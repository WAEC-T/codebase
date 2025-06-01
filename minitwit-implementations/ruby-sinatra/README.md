# itu-minitwit

## Ruby setup guide

Requires Ruby 3.3.0, Bundler and Postgres

### Other prerequisites

```bashrc
sudo apt update && sudo apt install libpq-dev
```

### Setup and run

- Install required gems with `bundle install`
- Start the app with `bundle exec ruby myapp.rb`

# How to Run on Raspberry Pi

## Ubuntu Server 24.04.2 LTS

### Requirements
- Rust version >= 1.58 (See https://github.com/ruby/ruby/blob/master/doc/yjit/yjit.md)

### Installing and updating ruby to 3.4.0

See [this repository](https://github.com/rbenv/ruby-build) for a full guide

Run the [ruby-setup.sh](./ruby-setup.sh) script to install the correct ruby version. Note, the script does not install rust (see e.g., the [rest-setup.sh](../rust-actix/rust-setup.sh) script)

### How to Run Ruby with different configurations

#### Default
```bashrc
SECRET_KEY="waectsk" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" bundle exec ruby main.rb -o 0.0.0.0 -p 5000
```

#### YJIT
```bashrc
SECRET_KEY="waectsk" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" bundle exec ruby --yjit main.rb -o 0.0.0.0 -p 5000
```

#### jemalloc
```bashrc
LD_PRELOAD=/usr/lib/<device specific architecture>/libjemalloc.so.2 SECRET_KEY="waectsk" DATABASE_URL="postgresql://user:password@<ip-address>:5432/waect" bundle exec ruby main.rb -o 0.0.0.0 -p 5000
```
