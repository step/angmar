FROM golang:alpine
WORKDIR /go/src/github.com/step/angmar/

ADD . .
RUN apk update && apk add --no-cache git ca-certificates make && update-ca-certificates && go get ./... && make angmar

FROM alpine
WORKDIR /app
COPY --from=0 /go/src/github.com/step/angmar/bin/angmar ./angmar
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
ENTRYPOINT ["sh","-c", "/app/angmar -queue $QUEUE -redis-address $REDIS_ADDRESS -redis-db $REDIS_DB -log-filename $LOG_FILE_NAME -log-path $LOG_FILE_PATH -source-path $SOURCE_PATH"]