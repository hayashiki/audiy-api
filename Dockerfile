FROM golang:1.17-alpine3.13 as go-builder
RUN apk --update --no-cache add make git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
#COPY credentials ./credentials
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -ldflags="-w -s ${PACKAGE}/version.Version=$BUILD_VERSION -X ${PACKAGE}/version.BuildTime=$(date --utc  +%FT%TZ)" -v -o bin/gqlserver cmd/gqlserver/main.go
RUN make build-api

FROM gcr.io/distroless/base

WORKDIR /app

# Copy the binary to the production image from the builder stage.
COPY --from=go-builder /app/bin/gqlserver /app/gqlserver
#COPY --from=go-builder /app/credentials /app/credentials

# Run the web service on container startup.
CMD ["./gqlserver"]
