FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# RUN apk add --no-cache tzdata ca-certificates


COPY --from=builder /app/main .

COPY --from=builder /app/ui ./ui
# COPY --from=builder /app/pkg ./pkg

RUN mkdir -p /app/submissions /app/problems /app/temp && \
    chown nobody:nobody /app/submissions /app/problems /app/temp && \
    chmod +x /app/main

USER nobody

EXPOSE 8090
EXPOSE 8081


CMD ["./main"]