FROM golang:1.16-buster AS build

WORKDIR /

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y git curl openssh-client gcc g++ musl-dev

RUN mkdir -p /src

COPY ./ /src/

RUN cd /src && go get ./... && go build

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*
        
COPY --from=build               /src/go-wspubsub / 
COPY config/route.json         /config/route.json 
COPY config/config.json        /config/config.json 

ENV PORT=80

ENTRYPOINT ["/go-wspubsub"]