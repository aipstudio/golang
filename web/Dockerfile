FROM golang

WORKDIR /opt/test
ADD . /opt/test

EXPOSE 8080

RUN go build

ENTRYPOINT [ "/opt/test/test" ]
