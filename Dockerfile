FROM cgr.dev/chainguard/go:latest as src

WORKDIR /src
COPY . .

# migrate

FROM src as build-migrate

ENV GOCACHE=/root/.cache/go-build
RUN RUN --mount=type=cache,target="/tmp/.cache/go-build" go build -o ./build/migrate ./cmd/jobs/migrate

FROM cgr.dev/chainguard/glibc-dynamic:latest as migrate

COPY --from=build-migrate /src/build/migrate /bin/migrate
CMD ["/bin/migrate"]

# mothership

FROM src as build-mothership

ENV GOCACHE=/root/.cache/go-build
RUN RUN --mount=type=cache,target="/tmp/.cache/go-build" go build -o ./build/mothership ./cmd/workers/mothership

FROM cgr.dev/chainguard/glibc-dynamic:latest as mothership

COPY --from=build-mothership /src/build/mothership /bin/mothership
CMD ["/bin/mothership"]

# api

FROM src as build-api

ENV GOCACHE=/root/.cache/go-build
RUN RUN --mount=type=cache,target="/tmp/.cache/go-build" go build -o ./build/api ./cmd/api

FROM cgr.dev/chainguard/glibc-dynamic:latest as api

COPY --from=build-api /src/build/api /bin/api
CMD ["/bin/api"]
