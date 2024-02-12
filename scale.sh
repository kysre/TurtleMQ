#!/bin/bash

if [ "$EUID" -ne 0 ]; then 
    echo "Run as root."
else
    if [ ! "$#" -eq 1 ]; then
        echo "1 argument required, $# provided."
        return 1
    fi

    python3 operation/diff-update-compose.py $1
    if [ ! $? -eq 0 ]; then
        return 1
    fi
    docker compose --file docker-compose.yaml up --remove-orphans -d
    docker compose --file docker-compose.yaml up -d --no-deps --force-recreate prometheus
    echo "Service scaled to $1."
fi