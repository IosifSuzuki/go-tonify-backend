#!/usr/bin/env bash

docker stop $(docker ps -q)
docker system prune -a --volumes -f