FROM golang

WORKDIR /opt/mm
ADD . /opt/mm

EXPOSE 3333

CMD go get github.com/OsoianMarcel/claymore-go
RUN go build

ENTRYPOINT [ "/opt/mm/mm" ]
