FROM alpine:3.14.0
ARG UNAME=nobody
ARG UID=1000
ARG GID=1000

COPY ./ekspose /usr/local/bin/ekspose

RUN addgroup -S $GID && adduser -S $UID -G $GID
USER $UID

ENTRYPOINT [ "/usr/local/bin/ekspose" ]