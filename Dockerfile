FROM golang:1.16 AS builder

WORKDIR /work
COPY . .
RUN ls -hal && make linux

FROM alpine:3.10
ENV JOB_NAME "test"
COPY --from=builder /work/bin/linux/jcli /usr/bin/jcli
RUN jcli config generate -i=false > ~/.jenkins-cli.yaml
COPY bin/build.sh /usr/bin/jclih
RUN chmod +x /usr/bin/jclih

ENTRYPOINT ["jclih"]
