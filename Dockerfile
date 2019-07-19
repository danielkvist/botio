# Build stage
FROM golang:1.12.7-alpine3.10 AS build
RUN apk add --no-cache git
RUN adduser -D -u 10000 daniel && \
    mkdir /app/ && chown daniel /app/
USER daniel
WORKDIR /app/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o tbot main.go

# Final stage
FROM alpine:3.10
LABEL maintainer="danielkvist@protonmail.com"
RUN apk add --no-cache ca-certificates
COPY --from=build /app/tbot /app/tbot
ENTRYPOINT ["./app/tbot"]