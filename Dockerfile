FROM golang:1.8

WORKDIR /go/src/
COPY ./public ./public
COPY ./templates ./templates
COPY ./src/go-teach-me ./go-teach-me

RUN go get github.com/google/uuid
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/sessions
RUN go get github.com/lib/pq
RUN go get golang.org/x/crypto/bcrypt

RUN go build go-teach-me/app

CMD ./app
