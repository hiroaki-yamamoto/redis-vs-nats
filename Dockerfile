
FROM golang:alpine
ARG PKGNAME
ENV PKGNAME=${PKGNAME}
RUN apk --no-cache --update upgrade && apk --no-cache add git gcc libc-dev ca-certificates
ENV GO111MODULE=on
RUN mkdir -p /opt/code
VOLUME [ "/opt/code" ]
WORKDIR /opt/code
ENTRYPOINT [ "./run.sh" ]
