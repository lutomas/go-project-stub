FROM golang:1.14.9

# Build service
COPY . /go/src/github.com/lutomas/go-project-stub
WORKDIR /go/src/github.com/lutomas/go-project-stub

RUN GO111MODULE=off make install-install-linux

FROM alpine:latest
ENTRYPOINT ["/bin/main-app", "serve"]

EXPOSE 9701
