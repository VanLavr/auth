FROM golang:latest

WORKDIR ./

COPY . .

RUN go build -o ./bin/main ./cmd/main.go

EXPOSE 8080

CMD ["./bin/main"]