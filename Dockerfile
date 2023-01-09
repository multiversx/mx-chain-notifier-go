FROM golang:1.17.6 as builder

MAINTAINER MultiversX

WORKDIR /mx
COPY . .

WORKDIR /mx/cmd/notifier

RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:20.04
COPY --from=builder /mx/cmd/notifier /mx

EXPOSE 8080

WORKDIR /mx

ENTRYPOINT ["./notifier"]
CMD ["--api-type", "rabbit-api"]
