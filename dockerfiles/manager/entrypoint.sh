#!/bin/sh
# Set strict error handling
set -e

# Source deployed contract addresses if available (written by deployer container)
if [ -f /deploy-data/deployed_addresses.env ]; then
    . /deploy-data/deployed_addresses.env
    export CHAIN_PROCESSOR_ADDRESS CHAIN_TEEAUTHENTICATOR_ADDRESS
fi

# Take ownership of the data directory.
# The '-R' flag makes it recursive.
chown -R appuser:appgroup ${MANAGER_DATA_FOLDER}

# Take ownership of the reports directory
chown -R appuser:appgroup "${MANAGER_REPORTS_FOLDER}"

# If LOG_SERVER_FOLDER is set, take ownership of it
if [ -n "${LOG_SERVER_FOLDER}" ]; then
    chown -R appuser:appgroup "${LOG_SERVER_FOLDER}"
fi

# Execute the main command passed to the script
exec "$@"
