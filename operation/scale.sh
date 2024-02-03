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

    docker service scale turtle-mq=$1

    echo "Service scaled to $1."
fi