# Copyright © 2022, 2024, Breu, Inc. <info@breu.io>
#
# We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
# is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
# the Software under the Apache License, Version 2.0, in which case the following will apply:
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
# the License.
#
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
# an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.

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
