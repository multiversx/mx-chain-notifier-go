FROM golang:1.17.6 as builder

MAINTAINER MultiversX

WORKDIR /elrond
COPY . .

WORKDIR /elrond/cmd/notifier

RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:20.04
COPY --from=builder /elrond/cmd/notifier /elrond

EXPOSE 8080

WORKDIR /elrond

ENTRYPOINT ["./notifier"]
CMD ["--api-type", "rabbit-api"]
