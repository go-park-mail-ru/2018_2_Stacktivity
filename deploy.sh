#!/usr/bin/env bash

cd /home/ubuntu/blep/back
sudo docker-compose stop
sudo docker-compose up --build -d