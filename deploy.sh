#!/usr/bin/env bash

cd /home/ubuntu/blep/back
sudo docker-compose down
sudo docker-compose build
sudo docker-compose up -d