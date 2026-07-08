# Reproducible Executor EIF Build (R5)

This document describes how to reproducibly build the Executor enclave image so
that the `PCR0` (and `PCR1`/`PCR2`) values registered on-chain via
`TeeAuthenticator.proposePcr0Swap` can be independently recomputed and verified
by third parties during the timelock window.

Design context: `EXECUTOR_TEE_UPGRADE_DESIGN.md` (R5) and Task 7 in
`EXECUTOR_TEE_UPGRADE_TASKS.md`.

## Why this is hard

`PCR0` is a SHA-384 over the **whole EIF**, which is assembled from three
layers, each of which must be deterministic:

1. **The Go binary** — with `CGO_ENABLED=1` (required by `wasmtime-go` and
   `go-ethereum/crypto/secp256k1`) the binary links the builder image's glibc,
   so the base image and the gcc/binutils versions affect its bytes.
2. **The runtime image filesystem** — it becomes the application ramdisk, so
   every byte, file timestamp, **and** build-time-variable *content* (dnf logs,
   the RPM/dnf databases' `INSTALLTIME`, `/etc/shadow`'s last-change day) lands
   in `PCR0`.
3. **The EIF packaging** — `nitro-cli build-enclave` bundles its own
   kernel/init/bootstrap blobs (`/usr/share/nitro_enclaves/blobs/`) into the
   EIF, so the exact `nitro-cli` version determines part of `PCR0`/`PCR1`.

Note that BuildKit's `rewrite-timestamp` normalizes file **mtimes** only; it
does **not** fix content-embedded dates, which is why the RPM database is
removed wholesale and `/etc/shadow` is normalized explicitly (see the
Dockerfile comments).

## Files

| File | Role |
|------|------|
| `dockerfiles/executor/Dockerfile` | Deterministic builder + runtime stages (Amazon Linux 2023); a pure function of its build args. |
| `dockerfiles/executor/versions.env` | Single source of truth for every pinned input (base digest, AL2023 releasever, Go checksum, package lists, BuildKit version). |
| `dockerfiles/executor/nitro-cli.Dockerfile` | Pinned `nitro-cli` toolbox used for EIF packaging. |
| `dockerfiles/executor/build-eif.sh` | Drives the deterministic buildx + nitro-cli pipeline and emits the EIF and measurements. |
| `.github/workflows/reproducible-eif.yml` | Builds twice on independent runners and asserts identical PCRs. |

## Pinning the inputs

Every value in `versions.env` is pinned. When bumping the toolchain, set the
values being replaced to `__TODO__` until re-resolved: `build-eif.sh` refuses to
run (and the CI check stays informational) while any placeholder remains. Each
value carries the exact command to resolve it in a comment; in summary:

- **Base image digest** — `docker buildx imagetools inspect amazonlinux:2023 --format '{{.Manifest.Digest}}'`. Builder and runtime share one base so the cgo binary's glibc matches the runtime glibc.
- **AL2023 releasever** — the versioned repo snapshot (e.g. `2023.6.20250428`) that freezes every `dnf`-installed package at once. Listed at https://cdn.amazonlinux.com/al2023/core/mirrors/ ; also derivable from `/etc/os-release` inside the pinned base.
- **Go checksum** — from https://go.dev/dl/ for `go<version>.linux-amd64.tar.gz`.
- **Package names** — `BUILDER_PACKAGES` / `RUNTIME_PACKAGES` / `NITRO_CLI_PACKAGES`; versions are fixed by the releasever, so no per-package EVRs are pinned. There is no NSM package on AL2023 — the executor's NSM path is pure-Go (see `versions.env`). `NITRO_CLI_PACKAGES` must include `aws-nitro-enclaves-cli-devel`, which ships the enclave blobs that `build-enclave` bakes into the EIF.
- **BuildKit version** — must be `>= v0.13` for the `rewrite-timestamp` exporter option.

## Building

```bash
# from the repo root, on a checkout that has the tag
./dockerfiles/executor/build-eif.sh v0.3.0 ./eif-out
```

Building an EIF does **not** require a Nitro-enabled host (only *running* an
enclave does). The script:

1. streams a pristine `git archive` of the tag as the build context (no `.git`,
   no untracked files) and passes `GIT_VERSION=<tag>` and `RELEASE=1`;
2. builds with `SOURCE_DATE_EPOCH` = the tag's commit time, `--provenance=false`,
   and `--output type=docker,...,rewrite-timestamp=true` on a buildx builder
   pinned to `moby/buildkit:<BUILDKIT_VERSION>`;
3. loads the image and runs `nitro-cli build-enclave` inside the pinned toolbox
   (host Docker socket mounted);
4. emits `executor.eif`, `measurements.json` (PCR0/1/2), `build-info.json`, and
   `blobs.sha256`.

## Verifying a release

Given a published release (git tag + expected PCRs), a third party runs:

```bash
git checkout <tag>
./dockerfiles/executor/build-eif.sh <tag> ./eif-out
jq .Measurements.PCR0 eif-out/measurements.json     # compare to the on-chain proposePcr0Swap value
```

CI performs the same build twice on independent runners and fails if any of
`PCR0`/`PCR1`/`PCR2` differ.

Per release, publish: the git tag, `PCR0`/`PCR1`/`PCR2`, the `nitro-cli`
version, and the base-image digest.

> **The verifiable output is the PCR set, not the `.eif` file hash.** `nitro-cli`
> stamps a wall-clock `BuildTime` (and docker `LastTagTime`) into the EIF's
> metadata **header**, which is *not* part of the measured sections. Two builds
> from the same source therefore produce EIFs that differ in the header (first
> differing byte is in the metadata region) but carry **identical PCR0/1/2** —
> and PCR0 is the on-chain trust anchor. The underlying Docker image *is*
> byte-reproducible (same layer digests), which is what makes the ramdisk (PCR2)
> and hence PCR0 deterministic. Do not gate verification on the EIF file hash.

## Known limits

- **RPM snapshot longevity.** AL2023's versioned repositories keep a pinned
  `releasever` snapshot fetchable long-term, so pinned packages remain available
  well beyond the timelock window (a strict improvement over AL2, whose mirrors
  did not serve old RPMs). If AWS ever retires a snapshot, vendor/mirror the RPMs
  with the release artifacts; the R5 guarantee only needs the timelock window.
- **Minimal image (future hardening).** A `scratch` + copied-glibc runtime image
  would sidestep the dnf/RPM-DB and `useradd` nondeterminism entirely, but must
  ship the exact glibc runtime **and** the NSS libraries (`libnss_*`) that glibc
  `dlopen`s at runtime for name resolution — these do not show up in `ldd` and
  are easy to miss. Kept as a follow-up rather than the default because it
  cannot be validated without a full enclave run.
