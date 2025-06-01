#!/bin/sh

apt update && apt install -y \
    libpq-dev \
    build-essential \
    ca-certificates \
    curl \
    libz-dev \
    libffi-dev \
    libedit-dev \
    libyaml-dev \
    libjemalloc2 

curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y \
    && echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> /root/.bashrc \
    && /root/.cargo/bin/rustc --version

export PATH="/root/.cargo/bin:${PATH}"
