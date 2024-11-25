FROM alpine:edge AS builder
WORKDIR /usr/src/waect-rust
RUN apk update && apk add --no-cache rust cargo libpq-dev

COPY . .
RUN cargo build --release 

FROM alpine:edge
WORKDIR /usr/src/waect-rust
COPY --from=builder /usr/src/waect-rust/target/release/waect-rust ./
COPY src/frontend/static/ ./src/frontend/static/
RUN apk update && apk add libc6-compat libgcc libpq

EXPOSE 5000

CMD ["./waect-rust"]