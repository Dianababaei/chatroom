FROM golang:1.20-alpine

WORKDIR /app

COPY . .

EXPOSE 8081

CMD ["go", "run", "cmd/client/main.go"]
