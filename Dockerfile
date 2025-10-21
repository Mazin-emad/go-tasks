FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-backend ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

COPY --from=builder /app/todo-backend .

COPY --from=builder /app/cmd/migrate/migrations ./migrations

RUN chown -R appuser:appuser /root/

USER appuser

EXPOSE 4001

CMD ["./todo-backend"]
