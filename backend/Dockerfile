# Multi-stage build: compile the API, ship a tiny static image.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /everyreview-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /everyreview-api /usr/local/bin/everyreview-api
EXPOSE 8080
ENTRYPOINT ["everyreview-api"]
