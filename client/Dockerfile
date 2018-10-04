FROM golang:1.11.0

RUN mkdir /app 
ADD . /app/
WORKDIR /app 

RUN go get github.com/jasonlvhit/gocron
RUN go get github.com/sfreiberg/gotwilio
RUN go get github.com/billputer/go-namecheap

ENV KEY "default"
ENV USERNAME "username"
ENV PASSWORD "password"

RUN go build -o main
CMD ["./main"]
