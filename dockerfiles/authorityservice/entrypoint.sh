#!/bin/sh
set -e

# Source deployed contract addresses if available (written by deployer container)
if [ -f /deploy-data/deployed_addresses.env ]; then
    . /deploy-data/deployed_addresses.env
    export CHAIN_PROCESSOR_ADDRESS
fi

# Take ownership of the data directory so the non-root user can write.
if [ -n "${MANAGER_REPORTS_FOLDER}" ]; then
    chown -R appuser:appgroup "${MANAGER_REPORTS_FOLDER}"
fi

# Take ownership of the deploy artifacts directory so the non-root user can store uploads.
if [ -n "${DEPLOY_ARTIFACTS_PATH}" ]; then
    chown -R appuser:appgroup "${DEPLOY_ARTIFACTS_PATH}"
fi

exec "$@"
