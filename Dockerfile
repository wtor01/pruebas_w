FROM golang:1.18-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

ARG CMD_ROUTE

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build $CMD_ROUTE

FROM debian:buster-slim


RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN mkdir -p ./api/openapi

USER 1000


COPY --from=builder /app/build /app/build
COPY --from=builder /app/api/openapi /app/api/openapi

CMD ["/app/build"]