FROM golang:1.21-bookworm AS build
WORKDIR /srv/build

ARG no_proxy
ARG NO_PROXY
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN apt update --fix-missing

ADD Makefile go.mod go.sum ./

RUN --mount=type=secret,id=ci-job-token \
    git config --global credential.helper store && \
    make dependencies && \
    PROTOC=/srv/protoc3/bin/protoc PROTOC_OPTIONS="-I/srv/protoc3/include -I." make generate && \
    rm ~/.git-credentials

COPY . .
RUN make leader


FROM debian:bookworm as final

RUN apt update && apt install -y ca-certificates

RUN apt update --fix-missing && \
    apt-get upgrade -y && \
    apt install -y ca-certificates && \
    apt install -y tzdata && \
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip && \
    unzip protoc-3.7.1-linux-x86_64.zip -d /srv/protoc3 && \
    apt-get clean

# Debug Tools
RUN apt install -y nano vim \
    curl traceroute iputils-ping net-tools \
    redis-tools screen bc && \
    apt clean

WORKDIR /srv/build

COPY --from=build /srv/build/. /srv/build

COPY --from=build /srv/build/leader /bin/

RUN echo $(ls /bin)

ENTRYPOINT ["/bin/leader", "serve"]