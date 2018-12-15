FROM ubuntu:18.04
MAINTAINER Silvman

WORKDIR /var/
ADD ./bin/game /var/

EXPOSE 8083
CMD ./game