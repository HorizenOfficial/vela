# --- Build Stage ---
FROM debian:bullseye-slim AS builder

# Install base dependencies
RUN apt-get update && apt-get install -y wget ca-certificates gnupg git

# Manuallt install Node.js 20 from NodeSource (the one in debian:bullseye-slim is too old)
RUN mkdir -p /etc/apt/keyrings && \
    wget -qO- https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" > /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && apt-get install -y nodejs

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
ARG GIT_VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X github.com/horizen-pes/pkg/version.Version=${GIT_VERSION}" \
    -o /manager ./cmd/manager

# --- Runtime Stage ---
FROM alpine:latest

# Install su-exec for user switching
RUN apk add --no-cache su-exec

#default location for data folder (can be overridden by docker compose or docker run)
ENV MANAGER_DATA_FOLDER="/data"
# default location for shared reports/artifacts data (entrypoint derives subfolders)
ENV SHARED_DATA_FOLDER="/shared-data"
#default location for logs folder (can be overridden by docker compose or docker run)
ENV LOG_SERVER_FOLDER="/logs"

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the entrypoint script
COPY dockerfiles/manager/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Copy the built binary from the builder stage
COPY --from=builder /manager /manager

VOLUME ${MANAGER_DATA_FOLDER}
VOLUME ${SHARED_DATA_FOLDER}
VOLUME ${LOG_SERVER_FOLDER}

# Set the entrypoint to our script.
# The user is still root here, so the script can chown.
ENTRYPOINT ["/entrypoint.sh"]

# The default command is to run the manager as appuser.
# The entrypoint script will execute this using `su-exec`.
CMD ["su-exec", "appuser", "/manager"]
