FROM alpine

COPY ./netdisk /tmp/netdisk

WORKDIR /tmp/

RUN chmod +x netdisk

EXPOSE 8080
ENTRYPOINT ["./netdisk"]