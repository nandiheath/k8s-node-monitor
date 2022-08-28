FROM golang:1.12.3 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go get -d -v ./...

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main cmd/monitor/main.go


FROM alpine:3.9.3

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /app/main main


EXPOSE $SERVER_PORT

ENTRYPOINT ["./main"]
