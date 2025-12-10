#!/bin/sh
set -e

# Take ownership of the data directory so the non-root user can write.
if [ -n "${MANAGER_REPORTS_FOLDER}" ]; then
    chown -R appuser:appgroup "${MANAGER_REPORTS_FOLDER}"
fi

exec "$@"
