# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the main application which includes web and judge logic
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install necessary runtime dependencies (like wget for healthcheck)
RUN apk add --no-cache tzdata ca-certificates wget

# Copy the single combined binary
COPY --from=builder /app/main .

# Copy necessary runtime assets (like UI files)
COPY --from=builder /app/ui ./ui
# Only copy pkg if it contains non-Go runtime assets (templates, etc.)
# COPY --from=builder /app/pkg ./pkg

# Create directories needed by EITHER webapp OR the internal judge logic
RUN mkdir -p /app/submissions /app/problems /app/temp && \
    chown nobody:nobody /app/submissions /app/problems /app/temp && \
    chmod +x /app/main

USER nobody

# Expose BOTH ports the application listens on
EXPOSE 8090
EXPOSE 8081

# Healthcheck the main web application endpoint
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -q --spider http://localhost:8090/health || exit 1

# Run the single combined binary
CMD ["./main"]