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

COPY ./lua_jwt_scopes_check /usr/local/custom/kong/plugins/lua_jwt_scopes_check
RUN cd /usr/local/custom/kong/plugins/lua_jwt_scopes_check && luarocks make kong-plugin-jwt-auth-0.1.0-1.rockspec --local
