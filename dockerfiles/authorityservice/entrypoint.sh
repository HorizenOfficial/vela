#!/bin/sh
set -e

# Take ownership of the data directory so the non-root user can write.
if [ -n "${AUTHORITY_SERVICE_DATA_PATH}" ]; then
    chown -R appuser:appgroup "${AUTHORITY_SERVICE_DATA_PATH}"
fi

exec "$@"
