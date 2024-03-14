# build stage
FROM golang:alpine as stage-builder

RUN apk update
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY cmd cmd

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o binary ./cmd/

# release stage
FROM alpine:latest as stage-release

COPY --from=stage-builder /app/binary /

ENTRYPOINT ["/binary"]
