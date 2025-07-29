# ---------- Stage 1: Go backend ----------
FROM golang:1.23-alpine AS go-build
WORKDIR /app

# Copy go mod files
COPY app/backend/go.mod app/backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY app/backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# ---------- Stage 2: Vue 3 frontend ----------
FROM node:20-alpine AS vue-build
WORKDIR /app

# Copy package files
COPY app/frontend/package*.json ./
RUN npm ci

# Copy frontend source code
COPY app/frontend/ ./
RUN npm run build

# ---------- Stage 3: Final runtime ----------
FROM nginx:alpine

# Install supervisor and curl for health checks
RUN apk add --no-cache supervisor curl

# Copy Go backend binary
COPY --from=go-build /app/main /usr/local/bin/main
RUN chmod +x /usr/local/bin/main

# Copy Vue.js frontend build
COPY --from=vue-build /app/dist /usr/share/nginx/html

# Copy configuration files
COPY nginx.conf /etc/nginx/nginx.conf
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Create supervisor configuration
RUN mkdir -p /etc/supervisor/conf.d
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Create necessary directories
RUN mkdir -p /var/log/supervisor

# Expose ports
EXPOSE 80 8000

# Start supervisor to manage all processes
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]