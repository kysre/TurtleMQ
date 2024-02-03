#!/bin/bash

if [ "$EUID" -ne 0 ]; then 
    echo "Run as root."
else
    docker image build --tag leader:latest --target leader .
    docker image build --tag data-node:latest --target data-node .

    docker swarm leave --force
    docker swarm init --advertise-addr 192.168.99.100

    docker service create --replicas 1 --name turtle-mq leader

    echo "Service initialized."
fi