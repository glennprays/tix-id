# build stage
FROM golang:1.20.4-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# redis stage
FROM redis:6.2-alpine3.17 AS redis
RUN apk add --no-cache bash

# run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=redis /usr/local/bin/redis-server /usr/local/bin/
COPY .env .

# Install timezone data
RUN apk add --no-cache tzdata \
    && cp -r /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

CMD ["sh", "-c", "redis-server --daemonize yes & /app/main"]
