FROM golang:1.26.2 AS builder

WORKDIR /multiversx
COPY . .

WORKDIR /multiversx/cmd/notifier

RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:24.04
COPY --from=builder /multiversx/cmd/notifier /multiversx

EXPOSE 8080

WORKDIR /multiversx

ENTRYPOINT ["./notifier"]
CMD ["--publisher-type", "rabbitmq"]
