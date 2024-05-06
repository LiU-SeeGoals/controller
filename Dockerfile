FROM golang:1.21-alpine

WORKDIR /var/controller

COPY . .
RUN go mod download
CMD ["go", "run", "cmd/main.go"]
