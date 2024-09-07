FROM golang:1.21-alpine

# Install bash
RUN apk add --no-cache bash

WORKDIR /var/controller

COPY . .
RUN go mod download
CMD ["go", "run", "cmd/main.go"]
