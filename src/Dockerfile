FROM golang:1.25-bookworm AS base

WORKDIR /src

RUN apt-get update \
    && apt-get install -y --no-install-recommends gcc libc6-dev \
    && rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src .

FROM base AS development

RUN go install github.com/go-delve/delve/cmd/dlv@v1.25.2 \
    && mkdir -p /data

ENV GIN_MODE=debug
EXPOSE 8080 2345

# Rebuild mounted source on each container start; keep SQLite in /data.
CMD ["sh", "-c", "go build -gcflags='all=-N -l' -o /tmp/api . && cd /data && exec dlv exec /tmp/api --headless --listen=:2345 --api-version=2 --accept-multiclient --continue"]

FROM base AS builder

RUN go build -trimpath -ldflags="-s -w" -o /out/api .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --gid 10001 app \
    && useradd --uid 10001 --gid app --no-create-home app \
    && mkdir /data \
    && chown app:app /data

COPY --from=builder /out/api /usr/local/bin/api

WORKDIR /data
USER app
ENV GIN_MODE=release

EXPOSE 8080

CMD ["/usr/local/bin/api"]
