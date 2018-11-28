FROM ubuntu:18.04

MAINTAINER Falcon22

RUN apt-get -y update
RUN apt-get install -y systemd
RUN apt-get install -y redis-server

ENV PGVER 10
ENV GOVER 1.10

RUN apt-get install -y postgresql-$PGVER

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root

RUN apt install -y golang-$GOVER git

ENV GOROOT /usr/lib/go-$GOVER
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH/src/2018_2_Stacktivity/
ADD . $GOPATH/src/2018_2_Stacktivity/

RUN go install ./cmd/public-api/ && go install ./cmd/game/ && go install ./cmd/session/

EXPOSE 3000

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql -d docker -a -f ./storage/migrations/1_create_user_table.up.sql &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
EXPOSE 6379

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root
CMD service postgresql start && redis-server --daemonize yes && session && game && public-api
