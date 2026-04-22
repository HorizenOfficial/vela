#!/bin/sh
set -e

# Source deployed contract addresses if available (written by deployer container)
if [ -f /deploy-data/deployed_addresses.env ]; then
    . /deploy-data/deployed_addresses.env
    export CHAIN_PROCESSOR_ADDRESS
fi

: "${SHARED_DATA_FOLDER:=/shared-data}"
: "${MANAGER_REPORTS_FOLDER:=${SHARED_DATA_FOLDER}/reports}"
: "${DEPLOY_ARTIFACTS_PATH:=${SHARED_DATA_FOLDER}/artifacts}"

export MANAGER_REPORTS_FOLDER
export DEPLOY_ARTIFACTS_PATH

# Ensure shared subfolders exist before dropping privileges.
mkdir -p "${MANAGER_REPORTS_FOLDER}" "${DEPLOY_ARTIFACTS_PATH}"

# Take ownership of the shared data directories so the non-root user can write.
chown -R appuser:appgroup "${MANAGER_REPORTS_FOLDER}"
chown -R appuser:appgroup "${DEPLOY_ARTIFACTS_PATH}"

exec "$@"
