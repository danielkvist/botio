# Build stage
FROM golang:1.12.5-alpine3.9 AS build
RUN apk add --no-cache git

RUN adduser -D -u 10000 daniel
RUN mkdir /app/ && chown daniel /app/
USER daniel

WORKDIR /app/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tbot main.go

# Final stage
FROM alpine:3.9.4
LABEL maintainer="danielkvist@protonmail.com"

RUN apk add --no-cache ca-certificates 

RUN adduser -D -u 10000 daniel
RUN mkdir /data/ && chown daniel /data/
RUN mkdir /app/ && chown daniel /app/
USER daniel

COPY --from=build /app/tbot /app/tbot
ENTRYPOINT ["./app/tbot"]