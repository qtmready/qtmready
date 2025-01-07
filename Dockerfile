FROM cgr.dev/chainguard/wolfi-base AS builder

ARG version=1.21

WORKDIR /src

RUN apk update && apk search --no-cache -v -d libgit2-dev

# Install build dependencies - go build toolchain, and git2go build dependencies (libgit2)
RUN apk update && \
  apk add --no-cache \
  go \
  build-base \
  cmake \
  pkgconf \
  openssl-dev \
  pcre2-dev \
  zlib-dev \
  libssh2-dev \
  libgit2-dev=1.7.0-r2

# Copy source code
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  go mod download

COPY . .

# Build the quantm binary with static linking
RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
  go build -tags static,system_libgit2 \
  -o /build/quantm \
  ./cmd/quantm

# Runtime Stage (Static base)
FROM cgr.dev/chainguard/static AS quantm

COPY --from=builder /build/quantm /bin/quantm

ENTRYPOINT ["/bin/quantm"]
CMD []
