# syntax=docker/dockerfile:1
# ----- Stage 1: Build the Go binary, generate API docs, and gather assets -----
FROM golang:1.24 AS builder
WORKDIR /src

# First copy go.mod/go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire repository so that we get all source files and assets.
COPY . .

# Install swaggo and generate API docs.
# This runs from the repository root and uses your provided command.
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g server.go -d "cmd/server,edupage,icanteen" --parseInternal

# Build the Go binary from the ./cmd/server folder.
WORKDIR /src/cmd/server
RUN CGO_ENABLED=0 go build -o /server .

# ----- Stage 2: Final runtime image -----
FROM alpine:latest

# Install CA certificates if needed.
RUN apk add --no-cache ca-certificates

# Copy the built Go server binary (it may expect assets relative to a fixed location).
COPY --from=builder /server /server

# Copy the assets from the cmd/server/web folder into the image.
# They will be available at /cmd/server/web.
COPY --from=builder /src/cmd/server/web /cmd/server/web

# Copy the docs folder (from repository root) into the image.
COPY --from=builder /src/docs /docs

# Expose port 8080.
EXPOSE 8080

# Run the server.
ENTRYPOINT ["/server"]
