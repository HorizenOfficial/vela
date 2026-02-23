FROM node:20-slim

RUN apt-get update && apt-get install -y --no-install-recommends curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app/contracts

# Copy dependency files first for Docker layer caching
COPY contracts/package.json contracts/package-lock.json ./
RUN npm ci

# Copy contract sources and compile
COPY contracts/ ./
RUN npx hardhat compile

# Copy nitro_prover artifacts needed by deploy script at runtime.
# Using a wildcard so the build succeeds even if the directory doesn't exist yet;
# the deploy will fail at runtime if the artifacts are actually missing.
COPY nitro_prove[r]/ /app/nitro_prover/

COPY dockerfiles/deployer/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
