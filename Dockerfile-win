FROM alpine:3.3

RUN sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories
RUN apk add --no-cache ca-certificates curl

RUN curl -L https://github.com/jenkins-zh/jenkins-cli/releases/download/v0.0.23/jcli-windows-386.tar.gz|tar xzv && \
    mv jcli.exe /usr/local/bin/jcli.exe && \
    curl -L https://github.com/jenkins-zh/jenkins-cli/releases/download/v0.0.23/jcli-windows-386.txt -o jcli-windows-386.txt
