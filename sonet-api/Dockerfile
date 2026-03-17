FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /sonote-api ./cmd/sonet/

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /sonote-api /usr/local/bin/sonote-api
COPY --from=builder /app/internal/database/migrations /migrations

EXPOSE 8080
CMD ["sonote-api"]
