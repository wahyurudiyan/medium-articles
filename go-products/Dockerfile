FROM golang:1.24 AS builder

COPY . /src
WORKDIR /src

RUN go mod tidy
RUN mkdir -p bin/
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/server ./cmd/*

FROM debian:trixie-backports

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
COPY ./configs/config.yaml .

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./server", "-conf", "/data/conf"]
