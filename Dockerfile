FROM golang:alpine as builder

ENV GOPATH /go

WORKDIR /go/src

COPY . /go/src/go-discord-bot

RUN cd /go/src/go-discord-bot && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM alpine

WORKDIR /app

COPY --from=builder /go/src/go-discord-bot /app
COPY .env /app

# EXPOSE

CMD ["./go-discord-bot"]