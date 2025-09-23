#!/bin/sh
# Set strict error handling
set -e

# Take ownership of the data directory for the 'foundry' user.
# The '-R' flag makes it recursive.
chown -R foundry:foundry ${DATA_FOLDER}

# Execute the main command passed to the script,
# but run it as the 'foundry' user using gosu.
exec gosu foundry anvil --host 0.0.0.0 --port ${RPCPORT} --state ${DATA_FOLDER}/anvil_state.json --state-interval 5