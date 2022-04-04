FROM golang:1.17.6 as builder

MAINTAINER ElrondNetwork

WORKDIR /elrond
COPY . .

WORKDIR /elrond/cmd

RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:20.04
COPY --from=builder /elrond/cmd /elrond

EXPOSE 8080

WORKDIR /elrond

ENTRYPOINT ["./notifier"]
CMD ["--api-type", "notifier"]
