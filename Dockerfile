FROM alpine:3.8

WORKDIR "/opt"

ADD release/faas-swagger /opt/

RUN chmod 777 /opt/*

CMD ["/opt/faas-swagger"]
