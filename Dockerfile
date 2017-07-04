FROM golang:1.8

WORKDIR /go/src/go-teach-me
COPY ./public ./public
COPY ./templates ./templates
COPY ./src/go-teach-me .
COPY ./app.go .

RUN go-wrapper download
RUN go-wrapper install


CMD ["go-wrapper", "run"] # ["app"]
