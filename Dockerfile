FROM golang:1.14

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -v

CMD ["/app/rmq-go"]