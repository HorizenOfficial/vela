This folder contains an example docker compose to be used for debug/demo.<br>
It starts a dev chain using [Foundry Anvill](https://getfoundry.sh/anvil/overview) and two separate processes with the processor and the manager

## Instructions:

1)  build locally the docker images needed by the docker compose run *from the project root folder* the following commands:

    ```
    docker build -t horizen-pes-executor -f dockerfiles/executor/Dockerfile . 
    docker build -t horizen-pes-manager -f dockerfiles/manager/Dockerfile . 
    docker build -t horizen-pes-chain -f dockerfiles/chain/Dockerfile . 
    ```

2) Create an .env file using .env.template as draft

3) Start the environment with:

    ```
    cd dockerfiles
    docker compose up docker-compose.yml
    ```

## Additional info:

- the manager database and the chain data is persisted in two volumes (horizen-pes-manager-data and horizen-pes-chain-data).<br>
  To start from scratch, delete the volumes.
- to connect to the chain from Metamask, use the following parameters:
   - rpc url: http://localhost:8543
   - chainid: 31337

