# Use the official Foundry image which has all dependencies pre-installed
FROM ghcr.io/foundry-rs/foundry

#default location for data folder (can be overridden by docker compose or docker run)
ENV DATA_FOLDER="/data"
#default location for rpc port (can be overridden by docker compose or docker run)
ENV RPCPORT=8545

# Switch to root to install packages and set up permissions
USER root

# Install gosu, a lightweight sudo replacement
RUN apt-get update && apt-get install -y --no-install-recommends curl && \
    curl -sSL -o /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/1.17/gosu-$(dpkg --print-architecture)" && \
    chmod +x /usr/local/bin/gosu && \
    rm -rf /var/lib/apt/lists/*

# Copy the entrypoint script and make it executable
COPY dockerfiles/chain/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Expose the default Anvil RPC port
EXPOSE 8545

# Create a directory for the chain state
VOLUME ${DATA_FOLDER}

# Set the entrypoint to our script.
# The container will start as root to run this.
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
