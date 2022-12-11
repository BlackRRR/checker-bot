FROM golang:latest

WORKDIR /checker-bot

COPY . .

RUN go build ./cmd/checker-bot

CMD ["./checker-bot"]