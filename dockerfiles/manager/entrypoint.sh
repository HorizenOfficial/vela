#!/bin/sh
# Set strict error handling
set -e

# Take ownership of the data directory.
# The '-R' flag makes it recursive.
chown -R appuser:appgroup ${MANAGER_DATA_FOLDER}

# If MANAGER_REPORTS_FOLDER is set, take ownership of it.
if [ -n "${MANAGER_REPORTS_FOLDER}" ]; then
    chown -R appuser:appgroup "${MANAGER_REPORTS_FOLDER}"
fi


# If LOG_SERVER_FOLDER is set, take ownership of it
if [ -n "${LOG_SERVER_FOLDER}" ]; then
    chown -R appuser:appgroup "${LOG_SERVER_FOLDER}"
fi

# Execute the main command passed to the script
exec "$@"
