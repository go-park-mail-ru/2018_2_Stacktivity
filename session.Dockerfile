FROM ubuntu:18.04
MAINTAINER Silvman

WORKDIR /var/
ADD ./bin/session /var/

EXPOSE 8081
CMD ./session