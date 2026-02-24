FROM node:20-slim

RUN apt-get update && apt-get install -y --no-install-recommends curl jq && rm -rf /var/lib/apt/lists/*

WORKDIR /app/subgraph

# Copy dependency files first for Docker layer caching
COPY subgraphs/hcce/package.json subgraphs/hcce/package-lock.json ./
RUN npm ci

# Copy subgraph sources
COPY subgraphs/hcce/ ./

COPY dockerfiles/subgraph-deployer/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
