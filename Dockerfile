FROM golang:alpine as builder

RUN apk add --no-cache git gcc libc-dev
RUN go get github.com/Kong/go-pluginserver
RUN go get github.com/dgrijalva/jwt-go

RUN mkdir /go-plugins
COPY /plugins/example.go /go-plugins/
RUN go build -buildmode plugin -o /go-plugins/example.so /go-plugins/example.go

FROM kong:2.0.1-alpine

COPY --from=builder /go/bin/go-pluginserver /usr/local/bin/
RUN mkdir /tmp/go-plugins
COPY --from=builder /go-plugins/example.so /tmp/go-plugins
