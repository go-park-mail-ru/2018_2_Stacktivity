FROM ubuntu:18.04
MAINTAINER Silvman

WORKDIR /var/
ADD ./bin/public-api /var/

EXPOSE 8082
CMD ./public-api