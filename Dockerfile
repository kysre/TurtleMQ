FROM python:3.10-slim AS data-node

WORKDIR /usr/src/app

COPY datanode/requirements.txt ./
RUN pip install -r requirements.txt

COPY ./datanode/src ./src

EXPOSE 1234

FROM golang:1.21-bookworm AS leader-build
WORKDIR /srv/build

ARG no_proxy
ARG NO_PROXY

RUN apt update --fix-missing

ADD Makefile leader/go.mod leader/go.sum ./
RUN go mod download

COPY ./leader .
RUN go build -o $@ leader ./cmd/$@


FROM debian:bookworm as leader

RUN apt update --fix-missing && \
    apt-get upgrade -y && \
    apt install -y ca-certificates && \
    apt install -y tzdata && \
    ln -sf /usr/share/zoneinfo/UTC /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean

# Debug Tools
RUN apt install -y nano vim

WORKDIR /srv/build

COPY --from=leader-build /srv/build/. /srv/build

COPY --from=leader-build /srv/build/leader /bin/

RUN echo $(ls /bin)

EXPOSE 8080
EXPOSE 9000

ENTRYPOINT ["/bin/leader", "serve"]
