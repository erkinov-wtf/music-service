# Build stage
FROM golang:latest AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    TZ=${TIMEZONE}

WORKDIR /app

# Copy the entire project
COPY . .

# Download dependencies
RUN go mod download

# Build the binary
WORKDIR /app/cmd/music-service
RUN go build -o main .

# Atlas installation
RUN curl -sSf https://atlasgo.sh | sh

# Final stage for Go app
FROM debian:bullseye-slim AS final

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/cmd/music-service/main /app/main
COPY --from=builder /usr/local/bin/atlas /usr/local/bin/atlas

# Copy the migrations directory
COPY /migrations /app/migrations
COPY configs/config.yml /app/configs/config.yml
COPY configs/local.yml /app/configs/local.yml
COPY configs/release.yml /app/configs/release.yml

# Creating a shell script to run migrations and start the app
RUN echo '#!/bin/sh\n\
    /usr/local/bin/atlas migrate hash \n\
    /usr/local/bin/atlas migrate apply --dir "file:///app/migrations"  \
    --url "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"\n\
    /app/main' > /app/start.sh
RUN chmod +x /app/start.sh

CMD ["/app/start.sh"]