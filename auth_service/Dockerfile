FROM golang:1.25-alpine AS base
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY . .
EXPOSE 5002
CMD [ "air" ]

FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main ./cmd/api/main.go

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/bin/main .
EXPOSE 5002
CMD [ "./main" ]