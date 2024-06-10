# Dockerfile for sc-bos
# This dockerfile performs the entire build itself.
# To enable fetching of private dependencies, an npmrc file must be be injected using a secret.
#
# To build, assuming your .npmrc is set up on your machine, run:
#     docker/podman build --secret=id=npmrc,src=$HOME/.npmrc .

FROM --platform=$BUILDPLATFORM node:20 AS build_ui

WORKDIR /src

ENV YARN_CACHE_FOLDER=/yarn-cache

COPY ui/conductor/package.json ui/conductor/yarn.lock ./
RUN --mount=type=cache,target=/yarn-cache \
    --mount=type=secret,id=npmrc,target=/root/.npmrc \
    yarn install --frozen-lockfile --check-files

COPY ui/conductor ./
ARG GIT_VERSION="(unknown)"
ENV GIT_VERSION=$GIT_VERSION
RUN yarn build

FROM --platform=$BUILDPLATFORM golang:1.22-alpine3.20 AS build_go

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

COPY cmd ./cmd/
COPY internal ./internal/
COPY pkg ./pkg/

# set by the build engine
ARG TARGETARCH
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=$TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o sc-bos ./cmd/bos

FROM alpine:3.20
LABEL vendor="Vanti Ltd"

COPY --from=build_go /src/sc-bos /app/
COPY --from=build_ui /src/dist /app/ops-ui/
COPY default/cfg/ /cfg/
COPY default/ui-config/ /app/ui-config/

EXPOSE 443
EXPOSE 23557
EXPOSE 7000-7999

VOLUME ["/cfg", "/data"]

ENTRYPOINT ["/app/sc-bos", "--appconf=/cfg/app.conf.json", "--sysconf=/cfg/system.conf.json", "--data=/data"]
