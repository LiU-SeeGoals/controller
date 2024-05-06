FROM golang:1.21-alpine

WORKDIR /var/controller

COPY go.mod ./
COPY go.sum ./
RUN go mod download
CMD ["go", "run", "cmd/main.go"]
