
# syntax=docker/dockerfile:1

# Build stage
FROM golang:1-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=build /server /server

EXPOSE 8080

USER nobody:nogroup

CMD ["/server"]
