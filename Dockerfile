FROM golang:1.24-alpine

WORKDIR /app

# Installing build dependencies
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

# runtime dependencies
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/pkg ./pkg

RUN mkdir -p /app/submissions /app/problems

RUN chmod +x /app/main
RUN chown -R nobody:nobody /app

USER nobody

EXPOSE 8090

CMD ["./main"]