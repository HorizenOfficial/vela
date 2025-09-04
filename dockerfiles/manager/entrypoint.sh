#!/bin/sh
# Set strict error handling
set -e

# Take ownership of the data directory.
# The '-R' flag makes it recursive.
chown -R appuser:appgroup ${MANAGER_DATA_FOLDER}

# Execute the main command passed to the script
exec "$@"
