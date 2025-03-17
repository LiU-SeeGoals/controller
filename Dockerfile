FROM golang:1.21-alpine

# Install bash
RUN apk add --no-cache bash

COPY .bash_history /root/.bash_history
WORKDIR /var/controller

COPY . .
RUN go mod download
CMD ["go", "run", "cmd/main.go"]
