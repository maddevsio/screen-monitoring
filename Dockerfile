FROM golang:1.7
MAINTAINER Gennady Karev <pendolf666@gmail.com>

ADD . /go/src/github.com/maddevsio/screen-monitoring

WORKDIR /go/src/github.com/maddevsio/screen-monitoring

RUN go get -v && go build -v

EXPOSE 8080

CMD ["./screen-monitoring"]
