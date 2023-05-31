FROM golang:1.18-bullseye as base

WORKDIR /app

ENV GO111MODULE=on CGO_ENABLED=0

ADD . .

RUN go build -o /app/main /app/cmd/api/main.go

FROM alpine:3

WORKDIR /app

COPY --from=base /app/main /app/main

EXPOSE 8000

CMD ["/app/main"]