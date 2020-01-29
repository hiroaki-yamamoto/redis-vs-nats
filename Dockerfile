
FROM golang:alpine
RUN apk --no-cache --update upgrade && apk --no-cache add git gcc libc-dev ca-certificates
ENV GO111MODULE=on
RUN mkdir -p /opt/code
VOLUME [ "/opt/code" ]
WORKDIR /opt/code
ENTRYPOINT [ "./run.sh" ]
