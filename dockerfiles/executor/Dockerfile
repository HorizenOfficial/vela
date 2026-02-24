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

ENV GOPRIVATE=github.com/HorizenOfficial/*

COPY go.mod go.sum ./
RUN --mount=type=secret,id=gh_token \
    git config --global url."https://$(cat /run/secrets/gh_token):@github.com/".insteadOf "https://github.com/" && \
    go mod download && \
    rm -f /root/.gitconfig

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /executor ./cmd/executor

# --- Runtime Stage ---
FROM amazonlinux:2

RUN yum install -y aws-nitro-enclaves-nsm ca-certificates shadow-utils  && yum clean all

RUN useradd -m appuser
USER appuser

COPY --from=builder /executor /executor

ENTRYPOINT ["/executor"]
