FROM golang:1.11.0

RUN mkdir /app 
ADD . /app/
WORKDIR /app 

RUN go get github.com/billputer/go-namecheap
RUN go get github.com/gorilla/mux
RUN go build -o main

ENV KEY "default"

EXPOSE 8080
CMD ["./main"]