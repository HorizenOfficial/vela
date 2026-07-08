#!/usr/bin/env bash
# =============================================================================
# Reproducible Executor EIF build (Task 7 / R5). Base: Amazon Linux 2023.
#
# Produces an EIF with identical PCR0/1/2 from a git ref (the .eif file itself
# varies only in nitro-cli's non-measured BuildTime header), so a third party can
# rebuild from a tag and compare PCR0 against the on-chain `proposePcr0Swap`
# value during the timelock window.
#
# Usage:
#   dockerfiles/executor/build-eif.sh [GIT_REF] [OUTPUT_DIR]
#
#   GIT_REF     tag/commit to build (default: the exact tag on HEAD).
#   OUTPUT_DIR  where artifacts are written (default: ./eif-out).
#
# Emits into OUTPUT_DIR:
#   executor.eif        the enclave image
#   measurements.json   PCR0/1/2 (from `nitro-cli describe-eif`)
#   build-info.json     git ref, SOURCE_DATE_EPOCH, pins used, nitro-cli version
#   blobs.sha256        checksums of the nitro-cli bundled blobs
#
# Requirements: docker with buildx (BuildKit >= v0.13 for rewrite-timestamp),
# and permission to run a container that mounts the docker socket. Building an
# EIF does NOT require a Nitro-enabled host.
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# shellcheck disable=SC1091
source "${SCRIPT_DIR}/versions.env"

GIT_REF="${1:-$(git -C "${REPO_ROOT}" describe --tags --exact-match 2>/dev/null || true)}"
OUTPUT_DIR="${2:-${REPO_ROOT}/eif-out}"

if [ -z "${GIT_REF}" ]; then
  echo "ERROR: no GIT_REF given and HEAD is not a tagged release." >&2
  echo "       pass an explicit ref, e.g. $(basename "$0") v0.3.0" >&2
  exit 1
fi

# The build context (Dockerfile included) comes from `git archive ${GIT_REF}`,
# but versions.env above is sourced from THIS checkout - refuse to mix two
# revisions, or the pins may not match the archived Dockerfile.
if [ "$(git -C "${REPO_ROOT}" rev-parse HEAD)" != "$(git -C "${REPO_ROOT}" rev-parse "${GIT_REF}^{commit}")" ]; then
  echo "ERROR: HEAD is not ${GIT_REF} - versions.env would come from a different revision" >&2
  echo "       than the archived build context. Run: git checkout ${GIT_REF}" >&2
  exit 1
fi

# Fail closed on unresolved pins - a reproducible build with placeholders is a
# contradiction. See versions.env for how to resolve each value.
if grep -v '^[[:space:]]*#' "${SCRIPT_DIR}/versions.env" | grep -q '__TODO__'; then
  echo "ERROR: versions.env still contains __TODO__ placeholders." >&2
  echo "       Resolve every pinned input before a reproducible build." >&2
  exit 1
fi

GIT_VERSION="${GIT_REF}"
# SOURCE_DATE_EPOCH = commit time of the ref. git archive also stamps this as
# the mtime of every context file; BuildKit rewrites layer mtimes to it.
SOURCE_DATE_EPOCH="$(git -C "${REPO_ROOT}" log -1 --pretty=%ct "${GIT_REF}")"
export SOURCE_DATE_EPOCH

mkdir -p "${OUTPUT_DIR}"
IMAGE_TAR="${OUTPUT_DIR}/executor-image.tar"
IMAGE_REF="vela-executor:${GIT_VERSION}"
NITRO_CLI_IMAGE="vela-nitro-cli:${AL_RELEASEVER}"
BUILDER_NAME="vela-repro-$$"

# Dedicated buildx builder pinned to a known BuildKit (rewrite-timestamp needs
# the container driver + BuildKit >= v0.13). Removed on exit.
cleanup() { docker buildx rm "${BUILDER_NAME}" >/dev/null 2>&1 || true; }
trap cleanup EXIT
docker buildx create --name "${BUILDER_NAME}" --driver docker-container \
  --driver-opt "image=moby/buildkit:${BUILDKIT_VERSION}" >/dev/null

echo ">> Building reproducible executor image from ${GIT_REF} (SOURCE_DATE_EPOCH=${SOURCE_DATE_EPOCH}, releasever=${AL_RELEASEVER})"
# Pristine context straight from git: no .git, no untracked files. BuildKit
# rewrites every layer file timestamp to SOURCE_DATE_EPOCH; --provenance=false
# drops the SLSA attestation (it embeds build time).
git -C "${REPO_ROOT}" archive --format=tar "${GIT_REF}" | \
  docker buildx --builder "${BUILDER_NAME}" build \
    --file dockerfiles/executor/Dockerfile \
    --build-arg BASE_IMAGE="${BASE_IMAGE}" \
    --build-arg AL_RELEASEVER="${AL_RELEASEVER}" \
    --build-arg BUILDER_PACKAGES="${BUILDER_PACKAGES}" \
    --build-arg RUNTIME_PACKAGES="${RUNTIME_PACKAGES}" \
    --build-arg GO_VERSION="${GO_VERSION}" \
    --build-arg GO_SHA256="${GO_SHA256}" \
    --build-arg GIT_VERSION="${GIT_VERSION}" \
    --build-arg RELEASE=1 \
    --provenance=false \
    --output "type=docker,name=${IMAGE_REF},dest=${IMAGE_TAR},rewrite-timestamp=true" \
    -

echo ">> Loading image into the local Docker daemon"
docker load -i "${IMAGE_TAR}" >/dev/null

echo ">> Building pinned nitro-cli toolbox (${NITRO_CLI_IMAGE})"
docker buildx --builder "${BUILDER_NAME}" build \
  --file "${SCRIPT_DIR}/nitro-cli.Dockerfile" \
  --build-arg BASE_IMAGE="${NITRO_CLI_BASE}" \
  --build-arg AL_RELEASEVER="${AL_RELEASEVER}" \
  --build-arg NITRO_CLI_PACKAGES="${NITRO_CLI_PACKAGES}" \
  --provenance=false \
  --load \
  -t "${NITRO_CLI_IMAGE}" \
  "${SCRIPT_DIR}"

echo ">> Packaging EIF with nitro-cli build-enclave"
# nitro-cli reads the target image from the host Docker daemon via the mounted
# socket and assembles the EIF (no Nitro hardware required to build).
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "${OUTPUT_DIR}:/out" \
  "${NITRO_CLI_IMAGE}" \
  nitro-cli build-enclave \
    --docker-uri "${IMAGE_REF}" \
    --output-file /out/executor.eif \
  > "${OUTPUT_DIR}/build-enclave.log" 2>&1 || { cat "${OUTPUT_DIR}/build-enclave.log"; exit 1; }

echo ">> Reading canonical measurements (nitro-cli describe-eif)"
docker run --rm \
  -v "${OUTPUT_DIR}:/out" \
  "${NITRO_CLI_IMAGE}" \
  nitro-cli describe-eif --eif-path /out/executor.eif \
  > "${OUTPUT_DIR}/measurements.json"

echo ">> Recording blob checksums and build provenance"
# No fallback: missing blobs means the toolbox lacks aws-nitro-enclaves-cli-devel
# (the package that ships them) - that is a broken build, not an optional extra.
docker run --rm "${NITRO_CLI_IMAGE}" \
  sh -c 'sha256sum /usr/share/nitro_enclaves/blobs/*' \
  > "${OUTPUT_DIR}/blobs.sha256"
NITRO_CLI_VERSION="$(docker run --rm "${NITRO_CLI_IMAGE}" rpm -q aws-nitro-enclaves-cli 2>/dev/null || echo unknown)"

cat > "${OUTPUT_DIR}/build-info.json" <<EOF
{
  "gitRef": "${GIT_VERSION}",
  "sourceDateEpoch": ${SOURCE_DATE_EPOCH},
  "baseImage": "${BASE_IMAGE}",
  "alReleasever": "${AL_RELEASEVER}",
  "goVersion": "${GO_VERSION}",
  "nitroCliVersion": "${NITRO_CLI_VERSION}",
  "buildkitVersion": "${BUILDKIT_VERSION}"
}
EOF

echo ">> Done. PCR measurements:"
cat "${OUTPUT_DIR}/measurements.json"
