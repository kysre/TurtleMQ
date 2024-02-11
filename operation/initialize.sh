#!/bin/bash

if [ "$EUID" -ne 0 ]; then 
    echo "Run as root."
else
    cp operation/initial.docker-compose.yaml docker-compose.yaml
    docker compose --file docker-compose.yaml up --remove-orphans -d
    echo "Service initialized."
fi