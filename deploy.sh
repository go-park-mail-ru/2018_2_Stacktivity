#!/usr/bin/env bash

sudo mv /home/ubuntu/backend.conf /etc/nginx/conf.d/backend.conf
sudo nginx -s reload
cd /home/ubuntu/blep/back
sudo docker-compose build
sudo docker-compose restart