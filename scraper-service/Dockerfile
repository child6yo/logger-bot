FROM golang:1.24.3-alpine 

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
COPY ./ ./

RUN go mod tidy

RUN go build -o /scraper-service ./cmd/main.go

CMD ["/scraper-service"]
