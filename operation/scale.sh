#!/bin/bash

if [ "$EUID" -ne 0 ]; then 
    echo "Run as root."
else
    if [ ! "$#" -eq 1 ]; then
        echo "1 argument required, $# provided."
        return 1
    fi
    if ! [[ $1 =~ ^[0-9]+$ ]]; then
        echo "Numeric argument required, $1 provided."
        return 1
    fi
    if [ $1 -gt 10 -o $1 -lt 1 ]; then
        echo "Value $1 is too big or too small."; return 1
    fi

    python3 operation/diff-update-compose.py $1
    if [ ! $? -eq 0 ]; then
        return 1
    fi
    docker compose --file generated.docker-compose.yaml up --remove-orphans -d
    echo "Service scaled to $1."
fi