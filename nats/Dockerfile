
FROM golang:alpine AS build
RUN apk --no-cache --update upgrade && apk --no-cache add git gcc libc-dev ca-certificates
ENV GO111MODULE=on
RUN mkdir -p /opt/code
COPY . /opt/code
WORKDIR /opt/code
RUN go mod download
RUN go build -o /tmp/build ./nats

FROM alpine
COPY --from=build /tmp/build /bin/app
RUN mkdir /opt/code
VOLUME [ "/opt/code" ]
ENTRYPOINT [ "app" ]
