#!/bin/sh
# Set strict error handling
set -e

# Take ownership of the data directory.
# The '-R' flag makes it recursive.
chown -R appuser:appgroup /tmp/horizen-pes-data

# Execute the main command passed to the script
exec "$@"
