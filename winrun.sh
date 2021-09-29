#!/bin/bash

docker-compose down
docker-compose -d -f docker-compose-windows.yaml up

./backend/backend