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
    && go install github.com/air-verse/air@v1.64.5 \
    && mkdir -p /data

ENV GIN_MODE=debug

EXPOSE 8080 2345

CMD ["air", "-c", ".air.toml"]

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
