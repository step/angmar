FROM golang:alpine
WORKDIR /go/src/github.com/step/angmar/

ADD . .
RUN apk update && apk add --no-cache git ca-certificates make && update-ca-certificates && go get ./... && make angmar

FROM golang:alpine
WORKDIR /app
COPY --from=0 /go/src/github.com/step/angmar/bin/angmar ./angmar
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
EXPOSE 8009
ENTRYPOINT [ "/app/angmar" ]