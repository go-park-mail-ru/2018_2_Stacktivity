#!/usr/bin/env bash

# chmod 0777 ./deploy_local.sh
# if permission denied

# make sure that your ports 8082 (public-api), 8083 (game), 9090 (prometheus), 3000 (grafana) are open

# you will not be able to use prometheus or grafana localy with project configs; remove reverse proxy adaptation settings in
# prometheus and grafana configs or set up reverse proxy with the help of nginx.conf

# change manually your frontend origin in pkg/apps/game-serever/config.go and pkg/apps/public-api-server/config.go

# please do not rename directory of project, or you will have to change 2018_2_stacktivity_postgres_1
# to whatever your docker-compose will choose as the name of postgres container

sudo chmod 0777 ./make_bin.sh
./make_bin.sh
sudo docker-compose up -d
sudo docker cp storage/migrations/1_create_user_table.up.sql 2018_2_stacktivity_postgres_1:/var/
sudo docker exec 2018_2_stacktivity_postgres_1 useradd docker
sudo docker exec 2018_2_stacktivity_postgres_1 su docker -c "psql -d docker -a -f /var/1_create_user_table.up.sql"