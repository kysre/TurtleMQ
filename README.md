# TurtleMQ

This is a MVP of a message queueing system.

## How to test locally?

Test with docker compose:

```
docker compose -f test.docker-compose.yaml up -d
docker compose -f test.docker-compose.yaml logs -f
docker compose -f test.docker-compose.yaml down -d
```
