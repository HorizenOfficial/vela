#!/bin/sh
set -e

CHAIN_RPC_URL="${CHAIN_RPC_URL:-http://chain:8545}"
DEPLOY_OUTPUT_DIR="${DEPLOY_OUTPUT_DIR:-/deploy-data}"
DEPLOY_FILE="${DEPLOY_OUTPUT_DIR}/deployed_addresses.env"

echo "Deployer: waiting for chain at ${CHAIN_RPC_URL}..."

# Wait for the chain to be ready (up to 60 seconds)
RETRIES=30
i=0
while [ "$i" -lt "$RETRIES" ]; do
    if curl -sf -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
        "${CHAIN_RPC_URL}" > /dev/null 2>&1; then
        echo "Deployer: chain is ready."
        break
    fi
    i=$((i + 1))
    if [ "$i" -eq "$RETRIES" ]; then
        echo "Deployer: ERROR - chain not ready after ${RETRIES} attempts."
        exit 1
    fi
    sleep 2
done

# Check idempotency: if addresses file exists and contracts are actually deployed, skip
if [ -f "${DEPLOY_FILE}" ]; then
    echo "Deployer: found existing ${DEPLOY_FILE}, verifying contracts..."
    . "${DEPLOY_FILE}"

    # Check if ProcessorEndpoint has code deployed
    CODE=$(curl -sf -X POST -H "Content-Type: application/json" \
        --data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getCode\",\"params\":[\"${CHAIN_PROCESSOR_ADDRESS}\",\"latest\"],\"id\":1}" \
        "${CHAIN_RPC_URL}" | sed 's/.*"result":"\([^"]*\)".*/\1/')

    if [ -n "${CODE}" ] && [ "${CODE}" != "0x" ]; then
        echo "Deployer: contracts already deployed at ProcessorEndpoint=${CHAIN_PROCESSOR_ADDRESS}. Skipping."
        exit 0
    fi
    echo "Deployer: contracts not found on chain, redeploying..."
fi

echo "Deployer: deploying contracts..."
npx hardhat run scripts/deploy/all.ts --network local

# Verify output file was created
if [ ! -f "${DEPLOY_FILE}" ]; then
    echo "Deployer: ERROR - deploy succeeded but ${DEPLOY_FILE} was not created."
    exit 1
fi

echo "Deployer: deployment complete."
cat "${DEPLOY_FILE}"
exit 0
