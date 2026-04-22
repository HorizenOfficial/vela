# --- Build Stage ---
FROM amazonlinux:2 AS builder

RUN yum install -y gcc tar gzip wget git

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

# GIT_VERSION is auto-detected from .git/ when not provided via --build-arg.
# Do not exclude .git/ in .dockerignore, or version falls back to "dev".
ARG GIT_VERSION=
RUN if [ -z "$GIT_VERSION" ]; then \
      GIT_VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
    fi && \
    CGO_ENABLED=1 GOOS=linux go build \
      -ldflags "-X github.com/HorizenOfficial/vela/pkg/version.Version=${GIT_VERSION}" \
      -o /executor ./cmd/executor

# --- Runtime Stage ---
FROM amazonlinux:2

RUN yum install -y aws-nitro-enclaves-nsm ca-certificates shadow-utils  && yum clean all

COPY --from=builder /executor /executor

RUN useradd -m appuser
USER appuser

ENTRYPOINT ["/executor"]
