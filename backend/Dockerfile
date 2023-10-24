FROM golang:alpine

run apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main .

ENTRYPOINT ["/app/main"]