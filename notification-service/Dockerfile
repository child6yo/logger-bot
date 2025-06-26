FROM golang:1.24.3-alpine 

RUN apk add --no-cache git

WORKDIR /notification-service

COPY go.mod go.sum ./
COPY ./ ./

RUN go mod tidy

RUN go build -o notification-service ./cmd/main.go

CMD ["./notification-service"]