# =============================================================================
# Reproducible Executor enclave image (Task 7 / R5). Base: Amazon Linux 2023.
#
# PCR0 is a SHA-384 over the whole EIF; the runtime stage below becomes the
# application ramdisk, so EVERY byte and file timestamp in it lands in PCR0.
# This Dockerfile is therefore written to be a pure function of its build args.
#
# It is NOT meant to be built directly for releases: use ./build-eif.sh, which
# pins every input from versions.env and drives the deterministic buildx +
# nitro-cli pipeline (SOURCE_DATE_EPOCH + rewrite-timestamp + pinned releasever).
# A bare `docker build -f dockerfiles/executor/Dockerfile .` still works for dev
# (unpinned, .git auto-detected for the version string).
#
# AL2023 note: reproducibility relies on pinning the repo snapshot via
# --releasever (AL2023's versioned repositories keep those RPMs fetchable), NOT
# on per-package version strings. See versions.env.
# =============================================================================

# Global build args must be declared before the first FROM to be usable in it.
# Builder and runtime share ONE base image so the cgo (wasmtime-go) binary's
# glibc matches the runtime glibc exactly. build-eif.sh passes a @sha256 digest.
ARG BASE_IMAGE=amazonlinux:2023

# --- Build Stage ---
FROM ${BASE_IMAGE} AS builder

# AL2023 versioned-repo pin: freezes the entire package universe to one snapshot
# so gcc/binutils (which affect the cgo-linked binary) are deterministic. Empty
# in dev builds -> dnf uses the latest snapshot.
ARG AL_RELEASEVER=
ARG BUILDER_PACKAGES="gcc tar gzip git"
RUN dnf -y ${AL_RELEASEVER:+--releasever=${AL_RELEASEVER}} install ${BUILDER_PACKAGES} && \
    dnf clean all

# Install Go from the official tarball and verify its checksum. curl is
# preinstalled on AL2023 (curl-minimal), so no extra package is needed.
# Keep in sync with go.mod's `toolchain` directive (see GOTOOLCHAIN=local below).
ARG GO_VERSION=1.24.3
ARG GO_SHA256=
RUN curl -fsSL "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tar.gz && \
    if [ -n "${GO_SHA256}" ]; then \
      echo "${GO_SHA256}  /tmp/go.tar.gz" | sha256sum -c - ; \
    else \
      echo "WARNING: GO_SHA256 not set - skipping checksum verification (dev build)"; \
    fi && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
# Use ONLY the toolchain installed above; never auto-download the one named by
# go.mod's `toolchain` directive (that would fetch a different, un-pinned
# compiler at build time). GO_VERSION must satisfy go.mod, else the build fails.
ENV GOTOOLCHAIN=local

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Release builds use a pristine `git archive` context, which has NO .git/, so
# `git describe` cannot run - GIT_VERSION MUST be passed explicitly (build-eif.sh
# does this and sets RELEASE=1). A dev `docker build` (with .git copied by
# `COPY . .`) still auto-detects the version.
#
# -trimpath strips absolute build paths; -buildvcs=false drops the VCS stamp
# (path/dirty-flag nondeterminism) - the version already comes from ldflags.
ARG GIT_VERSION=
ARG RELEASE=
RUN set -eu; \
    if [ -n "${RELEASE}" ] && [ -z "${GIT_VERSION}" ]; then \
      echo "ERROR: GIT_VERSION must be set for release builds (git archive context has no .git to describe)" >&2; \
      exit 1; \
    fi; \
    if [ -z "${GIT_VERSION}" ]; then \
      GIT_VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
    fi; \
    CGO_ENABLED=1 GOOS=linux go build -trimpath -buildvcs=false \
      -ldflags "-X github.com/HorizenOfficial/vela/pkg/version.Version=${GIT_VERSION}" \
      -o /executor ./cmd/executor

# --- Runtime Stage ---
# Same BASE_IMAGE as the builder -> identical glibc for the cgo binary. This
# stage's flattened filesystem is measured into PCR0.
FROM ${BASE_IMAGE}

ARG AL_RELEASEVER=
# No NSM package: it does not exist on AL2023 and the executor's NSM path is
# pure-Go ioctl (github.com/hf/nsm). ca-certificates: outbound KMS TLS;
# shadow-utils: the service user.
ARG RUNTIME_PACKAGES="ca-certificates shadow-utils"
# Install runtime deps and strip every build-time-variable artifact in the same
# layer so it never reaches the measured rootfs:
#  - dnf cache/logs carry timestamps;
#  - the RPM database (sqlite on AL2023) and dnf history DB embed per-package
#    INSTALLTIME (content, not a file mtime, so BuildKit's rewrite-timestamp
#    cannot normalize it) - the runtime never runs dnf/rpm again, so both are
#    removed wholesale.
#  - ldconfig's aux-cache embeds the mtimes of the just-installed libraries as
#    seen at install time (again content, not the file's own mtime), so it too
#    varies per build and is removed. /etc/ld.so.cache is KEPT (needed by the
#    dynamically-linked cgo binary; its content is deterministic given the same
#    libraries - it stores no mtimes).
RUN dnf -y ${AL_RELEASEVER:+--releasever=${AL_RELEASEVER}} install ${RUNTIME_PACKAGES} && \
    dnf clean all && \
    rm -rf /var/cache/dnf /var/cache/ldconfig /var/log/dnf* /var/log/hawkey.log \
           /var/lib/rpm /var/lib/dnf /tmp/*

COPY --from=builder /executor /executor

# Create the service user, then normalize /etc/shadow: `useradd` stamps the
# last-password-change field (field 3, days since epoch) with the build DAY,
# which is content (not an mtime) and would make PCR0 differ across build days.
# Zero it, and drop the useradd backup files and never-logged-in log stubs.
RUN useradd -m appuser && \
    sed -i 's/^\(appuser:[^:]*\):[0-9]*:/\1:0:/' /etc/shadow && \
    rm -f /etc/passwd- /etc/group- /etc/shadow- /etc/gshadow- && \
    rm -rf /var/log/lastlog /var/log/faillog

# R2 key-continuity guard (Task 6): an enclave receives no environment at
# launch, so this MUST be baked into the image - nitro-cli folds ENV into the
# measured ramdisk, making the guard's state part of PCR0 and thus attestable.
# Default ON: a release/upgrade image refuses to generate a fresh keyset when
# the Manager reports no recovery data (a wiped/wrong data folder would orphan
# all encrypted state). ONLY a genesis (first-install) build may pass
# EXPECT_EXISTING_KEYSET=false; that image bootstraps the keyset once and must
# not remain in the accepted PCR0 set / KMS key policy afterwards.
ARG EXPECT_EXISTING_KEYSET=true
ENV EXECUTOR_EXPECT_EXISTING_KEYSET=${EXPECT_EXISTING_KEYSET}

USER appuser

ENTRYPOINT ["/executor"]
