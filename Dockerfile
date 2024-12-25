FROM golang:1.23 AS builder

WORKDIR /work
COPY . .
RUN CGO_ENABLED=0 go build -v -a -o jcli .

FROM alpine:3.10
COPY --from=builder /work/jcli /usr/bin/jcli
RUN jcli config generate -i=false > ~/.jenkins-cli.yaml

ENTRYPOINT ["jcli"]
