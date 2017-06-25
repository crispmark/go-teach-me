FROM golang:1.8

WORKDIR /go/src/
COPY ./public ./public
COPY ./src/go-teach-me ./go-teach-me

RUN go get github.com/lib/pq
RUN go get golang.org/x/crypto/bcrypt

RUN go build go-teach-me/app

CMD ./app
