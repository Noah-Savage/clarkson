# Multi-stage build for Clarkson

# Stage 1: Build backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /build
RUN apk add --no-cache sqlite-dev gcc musl-dev

COPY backend/ .
RUN CGO_ENABLED=1 GOOS=linux go build -o clarkson-server main.go models.go handlers.go routes.go notifications.go uploads.go imports.go reports.go

# Stage 2: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /build
COPY frontend/ .

RUN npm install && npm run build

# Stage 3: Runtime
FROM alpine:3.18

RUN apk add --no-cache ca-certificates sqlite-libs curl

RUN addgroup -g 100 users && \
    adduser -D -u 99 -G users nobody

RUN mkdir -p /config /assets /app && \
    chown -R nobody:users /config /assets /app

WORKDIR /app

COPY --from=backend-builder /build/clarkson-server .
COPY --from=frontend-builder /build/dist ./public

RUN chown -R nobody:users /app

USER nobody

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

ENV PORT=3000 \
    CONFIG_PATH=/config \
    ASSETS_PATH=/assets \
    JWT_SECRET=change-me-in-production

CMD ["./clarkson-server"]
