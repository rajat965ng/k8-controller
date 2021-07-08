FROM alpine:3.14.0
MAINTAINER
ARG UNAME=nobody
ARG UID=1000
ARG GID=1000

RUN apk add openjdk11
# Create a group and user
RUN addgroup -S $GID && adduser -S $UID -G $GID

# Tell docker that all future commands should run as the appuser user
USER $UID

#RUN adduser --disabled-password --gecos '' --home \$pwd --ingroup $UNAME --no-create-home --uid "$UID" $UNAME
#RUN groupadd -g \$GID \$UNAME

#RUN useradd -m -u \$UID -g \$GID -s /bin/java \$UNAME
#USER \$UNAME
RUN java -version && javac -version