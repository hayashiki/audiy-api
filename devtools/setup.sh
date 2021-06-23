go get -u github.com/99designs/gqlgen
gqlgen init

mkdir schema
put on schema files

rm -rf graph

gqlgen

mkdir -p cmd/gqlserver
mv server.go cmd/gqlserver/main.go

mkdir -p infrastructure/auth
mkdir -p infrastructure/ds
mkdir -p interfaces/api

mkdir -p application/usecase
mkdir -p domain/model
mkdir -p domain/repository
mkdir -p etc/config
mkdir -p etc/secrets
mkdir -p etc/utils

mkdir -p interfaces/middleware

