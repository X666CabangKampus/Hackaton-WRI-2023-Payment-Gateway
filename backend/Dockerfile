FROM golang:alpine

run apk update && apk add --no-cache git

WORKDIR /app

COPY . .

COPY ../conf/ /

RUN go mod tidy

RUN go build -o main .

ENTRYPOINT ["./main"]