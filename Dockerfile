FROM golang:1.20.7 as builder

MAINTAINER MultiversX

WORKDIR /multiversx
COPY . .

WORKDIR /multiversx/cmd/notifier

RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:22.04
COPY --from=builder /multiversx/cmd/notifier /multiversx

EXPOSE 8080

WORKDIR /multiversx

ENTRYPOINT ["./notifier"]
CMD ["--api-type", "rabbit-api"]
