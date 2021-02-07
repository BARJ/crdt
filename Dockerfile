FROM golang:1.15-alpine
WORKDIR /src
RUN apk add --no-cache git bash
COPY . .
RUN go test ./... -cover

