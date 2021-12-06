FROM golang:1.16.0-alpine3.13

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

RUN go get -u github.com/cosmtrek/air

CMD ["air -c .air.toml"]
