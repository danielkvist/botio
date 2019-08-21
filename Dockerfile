# Build stage
FROM golang:1.12.9-alpine3.10 AS build
RUN apk add --no-cache git
WORKDIR /app/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o botio main.go

# Final stage
FROM alpine:3.10.2
LABEL maintainer="danielkvist@protonmail.com"
RUN apk add --no-cache ca-certificates
COPY --from=build /app/botio /app/botio
ENTRYPOINT ["./app/botio"]