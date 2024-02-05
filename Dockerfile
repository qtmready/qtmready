FROM cgr.dev/chainguard/go:latest as src

WORKDIR /src
COPY . .

# migrate

FROM src as build-migrate
RUN go build -o ./build/migrate ./cmd/jobs/migrate

FROM cgr.dev/chainguard/glibc-dynamic:latest as migrate

COPY --from=build-migrate /src/build/migrate /bin/migrate
CMD ["/bin/migrate"]

# mothership

FROM src as build-mothership
RUN go build -o ./build/mothership ./cmd/workers/mothership

FROM cgr.dev/chainguard/glibc-dynamic:latest as mothership

COPY --from=build-mothership /src/build/mothership /bin/mothership
CMD ["/bin/mothership"]

# API

FROM src as build-api
RUN go build -o ./build/api ./cmd/api/main.go

FROM cgr.dev/chainguard/glibc-dynamic:latest as api

COPY --from=build-api /src/build/api /bin/api
CMD ["/bin/api"]
