FROM cgr.dev/chainguard/go:latest AS base
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  go mod download

FROM base AS src

WORKDIR /src
COPY . .
RUN git status


# migrate
FROM src AS build-migrate
LABEL io.quantm.artifacts.app="quantm"
LABEL io.quantm.artifacts.component="migrate"

RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
  go build -o ./build/migrate ./cmd/jobs/migrate

FROM cgr.dev/chainguard/git:latest-glibc AS migrate

COPY --from=build-migrate /src/build/migrate /bin/migrate

ENTRYPOINT [ ]
CMD ["/bin/migrate"]


# mothership
FROM src AS build-mothership
LABEL io.quantm.artifacts.app="quantm"
LABEL io.quantm.artifacts.component="mothership"

RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
  go build -o ./build/mothership ./cmd/workers/mothership

FROM cgr.dev/chainguard/git:latest-glibc AS mothership

COPY --from=build-mothership /src/build/mothership /bin/

ENTRYPOINT [ ]
CMD ["/bin/mothership"]


# api
FROM src AS build-api
LABEL io.quantm.artifacts.app="quantm"
LABEL io.quantm.artifacts.component="api"

RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
  go build -o ./build/api ./cmd/api

FROM cgr.dev/chainguard/git:latest-glibc AS api

COPY --from=build-api /src/build/api /bin/api

ENTRYPOINT [ ]
CMD ["/bin/api"]
