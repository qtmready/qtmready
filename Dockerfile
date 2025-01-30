FROM cgr.dev/chainguard/wolfi-base AS builder

WORKDIR /src

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
  libgit2-dev=1.7.2-r0

# Copy source code
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  go mod download

COPY . .

# Check libgit2 version
RUN apk info libgit2-dev

RUN pkg-config --cflags --libs --static libgit2

# Explicitly remove rpath
RUN sed -i 's/-R\/usr\/lib//g' /usr/lib/pkgconfig/libgit2.pc

RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
  go build -x -v -tags static,system_libgit2 \
  -o /build/quantm \
  ./cmd/quantm

# Runtime Stage (Static base)
FROM cgr.dev/chainguard/static AS quantm

COPY --from=builder /build/quantm /bin/quantm

ENTRYPOINT ["/bin/quantm"]
CMD []
