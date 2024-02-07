#!/bin/bash

if [ "$EUID" -ne 0 ]; then 
    echo "Run as root."
else
    docker compose --file generated.docker-compose.yaml down
    rm -f generated.docker-compose.yaml
    echo "Service deleted."
fi