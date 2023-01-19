FROM golang:1.19.3-alpine3.16 AS builder

WORKDIR /app

COPY . .

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

CMD CompileDaemon \
    --build="go build -o ./main ./app/main.go" \
    --command="./main -i ./site -o ./build" \
    --exclude-dir="./build/" \
    --include "*.htm" \
    --include "*.css" \
    --include "*.js" \
    --include "*.png"
