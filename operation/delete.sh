#!/bin/bash

if [ "$EUID" -ne 0 ]; then 
    echo "Run as root."
else
    docker service rm turtle-mq
    docker swarm leave --force

    echo "Service deleted."
fi