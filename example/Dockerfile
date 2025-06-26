FROM golang:1.24.3-alpine 

RUN apk add --no-cache git

WORKDIR /app  

COPY . .  
RUN go build -o /logs-task task.go 

CMD ["/logs-task"]