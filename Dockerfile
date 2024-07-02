FROM golang:1.22-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod tidy

COPY . .

RUN go build -o main cmd/app/main.go

EXPOSE 8080

CMD ["./main"]