FROM centos:7
MAINTAINER Hugo Hiden <hhiden@redhat.com>
ARG BINARY=./influx

COPY ${BINARY} /opt/influx
ENTRYPOINT ["/opt/influx"]
