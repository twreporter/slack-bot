# syntax = docker/dockerfile:1-experimental

FROM --platform=${BUILDPLATFORM} golang:1.18-alpine AS base
WORKDIR /src
# ENV CGO_ENABLED=0

RUN apk add git
COPY go.* .
RUN go mod download
COPY . .

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/bot .

FROM base AS unit-test
RUN --mount=type=cache,target=/root/.cache/go-build go test -v .

FROM scratch AS bin-unix
COPY --from=build /out/bot /

FROM bin-unix AS bin-linux
FROM bin-unix AS bin-darwin

FROM scratch AS bin-windows
COPY --from=build /out/bot /bot.exe

FROM bin-${TARGETOS} AS bin
