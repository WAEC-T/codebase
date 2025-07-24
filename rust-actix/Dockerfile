FROM debian:bullseye-slim AS builder
WORKDIR /usr/src/waect-rust
RUN apt-get update && apt-get install -y curl build-essential libpq-dev libssl-dev pkg-config

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="/root/.cargo/bin:$PATH"

COPY . .
RUN cargo build --release

FROM debian:bullseye-slim
WORKDIR /usr/src/waect-rust
COPY --from=builder /usr/src/waect-rust/target/release/waect-rust ./
COPY src/frontend/static/ ./src/frontend/static/
RUN apt-get update && apt-get install -y libpq5 && rm -rf /var/lib/apt/lists/*

EXPOSE 5000

CMD ["./waect-rust"]