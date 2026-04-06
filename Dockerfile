# ==========================================
# Stage 1: Build the application
# ==========================================
# Use $BUILDPLATFORM to always compile natively on the machine running docker.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /out/talon .

# ==========================================
# Stage 2: Final minimal image
# ==========================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -S appuser
RUN mkdir -p /data && chown appuser:appgroup /data

WORKDIR /app

COPY --from=builder /out/talon .
RUN chown appuser:appgroup talon

USER appuser

ENV DB_PATH=/data/talon.db

VOLUME ["/data"]

CMD ["./talon"]
