FROM alpine:3.18

COPY frontier /usr/bin/frontier

EXPOSE 8080
EXPOSE 5556
ENTRYPOINT ["frontier"]