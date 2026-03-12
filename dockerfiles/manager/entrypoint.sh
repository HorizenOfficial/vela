#!/bin/sh
# Set strict error handling
set -e

# Source deployed contract addresses if available (written by deployer container)
if [ -f /deploy-data/deployed_addresses.env ]; then
    . /deploy-data/deployed_addresses.env
    export CHAIN_PROCESSOR_ADDRESS CHAIN_TEEAUTHENTICATOR_ADDRESS
fi

: "${MANAGER_DATA_FOLDER:=/data}"
: "${SHARED_DATA_FOLDER:=/shared-data}"
: "${MANAGER_REPORTS_FOLDER:=${SHARED_DATA_FOLDER}/reports}"
: "${MANAGER_ARTIFACTS_PATH:=${SHARED_DATA_FOLDER}/artifacts}"

export MANAGER_REPORTS_FOLDER
export MANAGER_ARTIFACTS_PATH

# Ensure data folders exist before dropping privileges.
mkdir -p "${MANAGER_DATA_FOLDER}" "${MANAGER_REPORTS_FOLDER}" "${MANAGER_ARTIFACTS_PATH}"

# Take ownership of the data directory.
# The '-R' flag makes it recursive.
chown -R appuser:appgroup "${MANAGER_DATA_FOLDER}"

# Take ownership of the reports directory
chown -R appuser:appgroup "${MANAGER_REPORTS_FOLDER}"

# Take ownership of the artifacts directory
chown -R appuser:appgroup "${MANAGER_ARTIFACTS_PATH}"

# If LOG_SERVER_FOLDER is set, take ownership of it
if [ -n "${LOG_SERVER_FOLDER}" ]; then
    mkdir -p "${LOG_SERVER_FOLDER}"
    chown -R appuser:appgroup "${LOG_SERVER_FOLDER}"
fi

# Execute the main command passed to the script
exec "$@"
