# --- Build Stage ---
FROM debian:bullseye-slim AS builder

# Install base dependencies
RUN apt-get update && apt-get install -y wget gcc

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

# Build the manager binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /executor ./cmd/executor

# Create a non-root user
RUN addgroup appgroup && adduser appuser && adduser appuser appgroup

# Set the user
USER appuser

# Copy the built binary from the builder stage
#COPY --from=builder /executor /executor

# Set the entrypoint
ENTRYPOINT ["/executor"]
