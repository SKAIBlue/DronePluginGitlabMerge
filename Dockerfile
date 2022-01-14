FROM alpine
ADD main /bin/
ENTRYPOINT /bin/main