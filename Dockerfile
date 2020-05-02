FROM golang:1.12 AS builder

WORKDIR /work
COPY . .
RUN ls -hal && make linux

FROM alpine:3.10
COPY --from=builder /work/bin/linux/jcli /usr/bin/jcli

ENTRYPOINT ["jcli"]
