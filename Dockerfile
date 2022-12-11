FROM golang:1.18-bullseye as base

WORKDIR /app

ENV GO111MODULE=on CGO_ENABLED=0

COPY go.mod go.sum /app/

COPY ./cmd/api/ /app/cmd/api/
COPY ./internal/ /app/internal/

RUN go build -o /app/main /app/cmd/api/main.go

FROM alpine:3.16

COPY --from=base /app/main main
RUN mkdir migrations

COPY ./internal/db/migrations/ ./migrations

EXPOSE 8000

CMD ["/main"]