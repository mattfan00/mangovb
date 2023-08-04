# syntax=docker/dockerfile:1
FROM golang:1.19-alpine

# install gcc needed for cgo
RUN apk add gcc musl-dev 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o scraper ./cmd/scraper/main.go

ENTRYPOINT [ "./scraper" ]
