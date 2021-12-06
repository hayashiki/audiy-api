FROM golang:1.17-alpine3.13 as go-builder
WORKDIR /app
COPY go.mod go.sum ./

RUN apk add --no-cache upx make git && \
    go version && \
    go mod download

# ffmpeg
COPY --from=mwader/static-ffmpeg:4.4.0 /ffmpeg /tmp/ffmpeg
COPY --from=mwader/static-ffmpeg:4.4.0 /ffprobe /tmp/ffprobe
RUN upx /tmp/ffmpeg

COPY . .
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -trimpath -ldflags="-w -s ${PACKAGE}/version.Version=$BUILD_VERSION -X ${PACKAGE}/version.BuildTime=$(date --utc  +%FT%TZ)" -v -o bin/gqlserver cmd/gqlserver/main.go
RUN make build-api

FROM gcr.io/distroless/base

WORKDIR /app
COPY --from=go-builder /app/bin/gqlserver /app/gqlserver
COPY --from=go-builder /tmp/ffmpeg /usr/bin/ffmpeg
COPY --from=go-builder /tmp/ffprobe /usr/bin/ffprobe

CMD ["./gqlserver"]
