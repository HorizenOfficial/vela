# --- Build Stage ---
FROM amazonlinux:2023 AS builder

RUN dnf install -y gcc tar gzip wget git

# Install Go
ENV GO_VERSION=1.24.0
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /kms-proxy ./cmd/kms_proxy

# --- Runtime Stage ---
FROM amazonlinux:2023

RUN dnf install -y ca-certificates && dnf clean all

COPY --from=builder /kms-proxy /kms-proxy

ENTRYPOINT ["/kms-proxy"]
