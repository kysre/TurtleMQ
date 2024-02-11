#!/bin/bash

if [ "$EUID" -ne 0 ]; then 
    echo "Run as root."
else
    docker compose --file docker-compose.yaml down
    rm -f docker-compose.yaml
    echo "Service deleted."
fi