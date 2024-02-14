FROM ubuntu:focal

RUN apt-get update && \
    apt-get install -y curl && \
    apt-get install -y git

RUN curl -OL https://dl.google.com/go/go1.21.5.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.21.5.linux-arm64.tar.gz && \
    rm go1.21.5.linux-arm64.tar.gz

ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

WORKDIR /home/app

COPY . /home/app/

RUN go mod tidy
RUN go mod vendor

RUN chmod +x main.sh
RUN chmod +x main.go

ENTRYPOINT [ "/home/app/main.sh" ]
