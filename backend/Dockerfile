############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

run apk update && apk add --no-cache git

WORKDIR /go-build/

COPY . .

RUN ls -la

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .

############################
# STEP 2 build a small image
############################
FROM alpine:latest

WORKDIR /app

COPY --from=builder /main /app/
COPY --from=builder /go-build/conf/ /app/

EXPOSE 8080

ENTRYPOINT ["/app/main"]