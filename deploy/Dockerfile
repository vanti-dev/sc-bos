FROM docker.io/node:18-alpine3.18 AS ui-build

RUN apk add --update --no-cache git

COPY ui/conductor ui/conductor
WORKDIR ui/conductor
RUN --mount=type=secret,id=npmrc,target=/root/.npmrc \
    yarn install && yarn build

FROM docker.io/golang:1.21-alpine3.18 AS ac-build

WORKDIR /work
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
RUN go build -o areacontroller ./cmd/bos

FROM docker.io/alpine:3.18

COPY --from=ui-build /ui/conductor/dist /static
COPY --from=ac-build /work/areacontroller /areacontroller

ENTRYPOINT ["/areacontroller"]
