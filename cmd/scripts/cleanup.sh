#!/bin/bash

docker-compose down --remove-orphans
docker container prune -f
docker image prune -a -f
docker volume prune -f
docker builder prune --all -f
