FROM --platform=${BUILDPLATFORM} golang:1.18-alpine AS build
WORKDIR /src
# ENV CGO_ENABLED=0
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN apk add git
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/bot .

FROM scratch AS bin-unix
COPY --from=build /out/bot /

FROM bin-unix AS bin-linux
FROM bin-unix AS bin-darwin

FROM scratch AS bin-windows
COPY --from=build /out/bot /bot.exe

FROM bin-${TARGETOS} AS bin
