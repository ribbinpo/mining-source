FROM golang:1.23.4-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 4000

CMD ["./main"]