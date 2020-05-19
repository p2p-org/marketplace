#!/bin/bash

docker build -t marketplace .
docker network create dwh-network
docker-compose -f demo.yml up -d
