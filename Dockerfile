FROM golang:1.15.7 as builder

MAINTAINER ElrondNetwork

WORKDIR /elrond
COPY . .

WORKDIR /elrond/cmd
RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:18.04
COPY --from=builder /elrond/cmd /elrond

EXPOSE 8080

WORKDIR /elrond

CMD ["./notifier"]
