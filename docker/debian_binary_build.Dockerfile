FROM golang:1.4.1
MAINTAINER Orlovsky Alexander <nordicdyno@gmail.com>

ENV GOPATH=/src
RUN mkdir -p /src/src/github.com/nordicdyno/resm-sketch
COPY . /src/src/github.com/nordicdyno/resm-sketch/
RUN go get -d -v github.com/nordicdyno/resm-sketch/resm && \
    go install github.com/nordicdyno/resm-sketch/resm
