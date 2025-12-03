# --- Build Stage ---
FROM debian:bullseye-slim AS builder

# Install base dependencies
RUN apt-get update && apt-get install -y wget

# Manually install a specific version of Go
ENV GO_VERSION=1.24.0
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

# Verify Go installation
RUN go version

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the authority service binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /authorityservice ./cmd/authorityservice

# --- Runtime Stage ---
FROM alpine:latest

# Install su-exec for user switching
RUN apk add --no-cache su-exec

# Default location for data folder (can be overridden by docker compose or docker run)
ENV AUTHORITY_SERVICE_DATA_PATH="/data"

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the entrypoint script
COPY dockerfiles/authorityservice/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Copy the built binary from the builder stage
COPY --from=builder /authorityservice /authorityservice

VOLUME ${AUTHORITY_SERVICE_DATA_PATH}

# Set the entrypoint to our script.
# The user is still root here, so the script can chown.
ENTRYPOINT ["/entrypoint.sh"]

# The default command is to run the authority service as appuser.
# The entrypoint script will execute this using `su-exec`.
CMD ["su-exec", "appuser", "/authorityservice"]
