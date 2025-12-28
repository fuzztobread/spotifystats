FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o spotistats cmd/spotistats/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/spotistats .
COPY --from=builder /app/web ./web
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/scripts/start.sh .

RUN chmod +x start.sh

EXPOSE 8080

CMD ["./start.sh"]
