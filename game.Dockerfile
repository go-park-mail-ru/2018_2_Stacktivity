FROM ubuntu:18.04
MAINTAINER Silvman

WORKDIR /var/
ADD game /var/

EXPOSE 8083
CMD ./game