#!/bin/bash

eval $(docker-machine env do-swarm-1)
docker-machine active
docker-compose build
docker-compose push
docker stack deploy --compose-file docker-compose.yml crypto-stack