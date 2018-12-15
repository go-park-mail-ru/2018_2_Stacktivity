FROM ubuntu:18.04
MAINTAINER Falcon22

RUN apt-get -y update
ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

VOLUME ["/var/lib/postgresql/data"]

ADD ./storage/migrations/1_create_user_table.up.sql /var/

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql -d docker -a -f /var/1_create_user_table.up.sql &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
CMD ["/usr/lib/postgresql/10/bin/postgres", "-D", "/var/lib/postgresql/10/main", "-c", "config_file=/etc/postgresql/10/main/postgresql.conf"]