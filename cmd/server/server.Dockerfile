FROM golang:1.20-alpine

WORKDIR /app

COPY . .


EXPOSE 8080

CMD ["go", "run", "cmd/server/main.go"]
