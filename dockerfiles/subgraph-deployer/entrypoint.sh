#!/bin/sh
set -e

DEPLOY_DATA_DIR="${DEPLOY_DATA_DIR:-/deploy-data}"
DEPLOY_FILE="${DEPLOY_DATA_DIR}/deployed_addresses.env"
GRAPH_NODE_URL="${GRAPH_NODE_URL:-http://subgraph-node:8020}"
IPFS_URL="${IPFS_URL:-http://subgraph-ipfs:5001}"
SUBGRAPH_NAME="${SUBGRAPH_NAME:-hcce}"

echo "Subgraph deployer: starting..."

# --- 1. Read deployed contract addresses ---
if [ ! -f "${DEPLOY_FILE}" ]; then
    echo "Subgraph deployer: ERROR - ${DEPLOY_FILE} not found."
    exit 1
fi
echo "Subgraph deployer: reading deployed addresses from ${DEPLOY_FILE}"
. "${DEPLOY_FILE}"

if [ -z "${CHAIN_PROCESSOR_ADDRESS}" ]; then
    echo "Subgraph deployer: ERROR - CHAIN_PROCESSOR_ADDRESS not set in ${DEPLOY_FILE}."
    exit 1
fi
echo "Subgraph deployer: ProcessorEndpoint address = ${CHAIN_PROCESSOR_ADDRESS}"

# --- 2. Wait for Graph Node to be ready ---
echo "Subgraph deployer: waiting for Graph Node at ${GRAPH_NODE_URL}..."
RETRIES=60
i=0
while [ "$i" -lt "$RETRIES" ]; do
    if curl -sf -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"subgraph_create","id":1,"params":{"name":"__healthcheck__"}}' \
        "${GRAPH_NODE_URL}" > /dev/null 2>&1; then
        echo "Subgraph deployer: Graph Node is ready."
        break
    fi
    i=$((i + 1))
    if [ "$i" -eq "$RETRIES" ]; then
        echo "Subgraph deployer: ERROR - Graph Node not ready after ${RETRIES} attempts."
        exit 1
    fi
    sleep 2
done

# --- 3. Generate subgraph-local.yaml from the template ---
echo "Subgraph deployer: generating subgraph-local.yaml..."
sed -e "s/network: horizen-testnet/network: local/" \
    -e "s/address: \"0x<contract_address>\"/address: \"${CHAIN_PROCESSOR_ADDRESS}\"/" \
    -e "s/startBlock: [0-9]*/startBlock: 0/" \
    subgraph.yaml > subgraph-local.yaml

echo "Subgraph deployer: generated subgraph-local.yaml:"
cat subgraph-local.yaml

# --- 4. Generate code from schema ---
echo "Subgraph deployer: running codegen..."
npx graph codegen subgraph-local.yaml

# --- 5. Create subgraph on Graph Node (ignore error if already exists) ---
echo "Subgraph deployer: creating subgraph '${SUBGRAPH_NAME}'..."
npx graph create --node "${GRAPH_NODE_URL}" "${SUBGRAPH_NAME}" || true

# --- 6. Deploy subgraph ---
echo "Subgraph deployer: deploying subgraph '${SUBGRAPH_NAME}'..."
npx graph deploy \
    --node "${GRAPH_NODE_URL}" \
    --ipfs "${IPFS_URL}" \
    --version-label v0.0.1 \
    "${SUBGRAPH_NAME}" \
    subgraph-local.yaml

# --- 7. Wait for the subgraph to sync at least one block ---
GRAPHQL_URL="http://subgraph-node:8000/subgraphs/name/${SUBGRAPH_NAME}"
echo "Subgraph deployer: waiting for subgraph to start syncing at ${GRAPHQL_URL}..."
RETRIES=60
i=0
while [ "$i" -lt "$RETRIES" ]; do
    RESULT=$(curl -sf -X POST -H "Content-Type: application/json" \
        --data '{"query":"{ _meta { block { number } } }"}' \
        "${GRAPHQL_URL}" 2>/dev/null || echo "")
    if echo "${RESULT}" | grep -q '"number"'; then
        echo "Subgraph deployer: subgraph is syncing. ${RESULT}"
        break
    fi
    i=$((i + 1))
    if [ "$i" -eq "$RETRIES" ]; then
        echo "Subgraph deployer: WARNING - subgraph not syncing after ${RETRIES} attempts, proceeding anyway."
        break
    fi
    sleep 2
done

echo "Subgraph deployer: done."
exit 0
