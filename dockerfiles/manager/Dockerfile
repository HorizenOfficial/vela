# --- Build Stage ---
FROM debian:bullseye-slim AS builder

# Install base dependencies
RUN apt-get update && apt-get install -y nodejs npm wget

# Manually install a specific version of Go
ENV GO_VERSION=1.24.0
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

# Verify Go installation
RUN go version

# Install solc
ENV SOLC_VERSION=0.8.28
RUN wget https://github.com/ethereum/solidity/releases/download/v$SOLC_VERSION/solc-static-linux \
    -O /usr/local/bin/solc && \
    chmod +x /usr/local/bin/solc

# Verify solc installation
RUN solc --version

# Set the working directory inside the container
WORKDIR /app

# Copy contract dependency files and install dependencies
# This is done before copying the rest of the source code to leverage Docker caching
COPY contracts/package.json contracts/package-lock.json ./contracts/
WORKDIR /app/contracts
RUN npm ci
WORKDIR /app

# Copy go.mod and go.sum to download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the manager binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /manager ./cmd/manager

# --- Runtime Stage ---
FROM alpine:latest

# Install su-exec for user switching
RUN apk add --no-cache su-exec

#default location for data folder (can be overridden by docker compose or docker run)
ENV MANAGER_DATA_FOLDER="/data"
#default location for reports folder (can be overridden by docker compose or docker run)
ENV MANAGER_REPORTS_FOLDER="/reports"

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the entrypoint script
COPY dockerfiles/manager/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Copy the built binary from the builder stage
COPY --from=builder /manager /manager

VOLUME ${MANAGER_DATA_FOLDER}
VOLUME ${MANAGER_REPORTS_FOLDER}

# Set the entrypoint to our script.
# The user is still root here, so the script can chown.
ENTRYPOINT ["/entrypoint.sh"]

# The default command is to run the manager as appuser.
# The entrypoint script will execute this using `su-exec`.
CMD ["su-exec", "appuser", "/manager"]
