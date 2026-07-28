# Pinned nitro-cli toolbox for reproducible EIF packaging (Task 7 / R5).
# Base: Amazon Linux 2023.
#
# nitro-cli bundles the kernel/init/bootstrap blobs (/usr/share/nitro_enclaves/
# blobs/) that land in PCR0/PCR1, so its exact version is part of the
# reproducibility contract. Pinning AL_RELEASEVER (below) freezes the
# nitro-cli RPM version, hence the blobs. This toolbox image is NOT itself
# measured, so its base need not be digest-pinned. build-eif.sh runs
# `nitro-cli build-enclave` inside it against the host Docker daemon (socket
# mounted); no docker CLI is installed because nitro-cli uses the Docker API
# directly over the socket.
ARG BASE_IMAGE=amazonlinux:2023
FROM ${BASE_IMAGE}

ARG AL_RELEASEVER=
# aws-nitro-enclaves-cli-devel ships the enclave blobs (kernel/init/cmdline/nsm
# under /usr/share/nitro_enclaves/blobs/) that build-enclave bakes into the EIF;
# the base cli package alone does not include them.
ARG NITRO_CLI_PACKAGES="aws-nitro-enclaves-cli aws-nitro-enclaves-cli-devel"
# On AL2023 these are in the default repos (amazon-linux-extras no longer
# exists), so a plain dnf install suffices.
RUN dnf -y ${AL_RELEASEVER:+--releasever=${AL_RELEASEVER}} install ${NITRO_CLI_PACKAGES} && \
    dnf clean all

ENTRYPOINT []
