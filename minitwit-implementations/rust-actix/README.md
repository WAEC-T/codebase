# Getting started

## Installing rust

Follow the guide to installing the rust language

[Rust installation guide](https://www.rust-lang.org/learn/get-started)

## Running the program

First you need to build the project. This is done from the root foulder in the project running the build command.

```bashrc
cargo build
```

After it has successfully build and no errors occoured run the project.

```bashrc
cargo run --bin waect-rust
```

In the case of of dependencies not being installed run - do not run unless errors are encountered

```bashrc
cargo install --path ./
```

Then in your bowser of choice connect to [localhost:5000](http://localhost:5000)

## Running the api

```bashrc
cargo run --bin api
```

## Resources

## Frameworks

- Actix web framework

## Dependencies

- actix-files - version "0.6.5"
- actix-web - version "4"
- actix-session - version "0.9.0"
- askama - version "0.12.1"
- askama_actix - version "0.14.0"
- chrono - version "0.4.34"
- filters - version "0.4.0"
- rusqlite - version "0.30.0"
- actix-web-flash-messages - version "0.4"
- serde - version "1.0.196"
- pwhash - version "1"
- actix-identity - version "0.7.0"
- md-5 - version "0.10.6"
- uuid - version "1.7.0"

## Important libraries

- Askama for rendering templates
- rusqlite for database handling
- pwhash to verify and create user secrets

# How to Run on Raspberry Pi

## Ubuntu Server 24.04.2 LTS

### Requirements
None

### Installing Rust

Run [rust-setup.sh](./rust-setup.sh)

### How to Run Ruby with different configurations

To build a Rust executable follow the approach below

```bashrc
cargo build --release
```

#### Default

Run the executable generated from the steps in the previous section

```bashrc
# cd codebase/rust-actix/
chmod +x target/release/waect-rust

./target/release/waect-rust
```

#### jemalloc

Using jemalloc does not require a rebuild, but setting an environment variable when running MiniTwit

Build the executable like the baseline (see above), and the start the MiniTwit application using the following command

```bashrc
LD_PRELOAD=/usr/lib/<device specific architecture>/libjemalloc.so.2 ./waect-rust
```

#### PGO

`PGO` takes a bit more work, as some workflow data has to be generated. We've already generated the data, but for transparency and replication reasons we show below how the data was generated

```bashrc
# cd /codebase/rust-actix/

rustup component add llvm-tools-preview

rm -rf /tmp/pgo-data

RUSTFLAGS="-Cprofile-generate=/tmp/pgo-data -Cllvm-args=-pgo-warn-missing-function" cargo build --release

chmod +x /target/release/waect-rust

./target/release/waect-rust

# Run 1-2 scenarios

~/.rustup/toolchains/<toolchain>/lib/rustlib/<target-triple>/bin/llvm-profdata merge -o /tmp/pgo-data/rust-merged.profdata /tmp/pgo-data

RUSTFLAGS="-Cprofile-use=/<exact path>/rust-merged.profdata" cargo build --release

```

Now the generated `waect-rust` has been optimized and can be used. The data used for optimizing the Rust implementation is stored in this repository (see [rust-merged.profdata](./rust-merged.profdata))
