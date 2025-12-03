This folder contains an example docker compose to be used for debug/demo.<br>
It starts a dev chain using [Foundry Anvill](https://getfoundry.sh/anvil/overview) and two separate processes with the processor and the manager

## Instructions:

1)  (Skip this step if you want to use the Docker hub official images)
    To build locally the docker images needed by the docker compose run *from the project root folder* the following commands:

    ```
    docker build -t horizen/cce-executor -f dockerfiles/executor/Dockerfile . 
    docker build -t horizen/cce-manager -f dockerfiles/manager/Dockerfile . 
    docker build -t horizen/cce-authorityservice -f dockerfiles/authorityservice/Dockerfile .
    docker build -t horizen/cce-chain -f dockerfiles/chain/Dockerfile . 
    ```

2) Switch to "dockerfiles" folder

    ```
    cd dockerfiles
    ```

3) Create an .env file using .env.template as draft


4) Start the environment with:

    ```
    docker compose up 
    ```

## Additional info:

- exposed addresses/ports (see `.env.template` for defaults):
  - authorityservice: listens on `${AUTHORITY_SERVICE_LISTEN_ADDRESS}` on `${AUTHORITY_SERVICE_IP_ADDRESS}` inside the internal network.

- the manager database and the chain data is persisted in two docker volumes (horizen-cce-manager-data and horizen-cce-chain-data).<br>
  To start from scratch, delete the volumes.
- the authority service shares the manager database volume (horizen-cce-manager-data) so it can read the same deanonymization reports.
- to connect to the chain from Metamask, use the following parameters:
   - rpc url: http://localhost:8545
   - chainid: 31337

## Where to go next: 

- The Anvil chain node is created empty: to have a running dev environment you must deploy the contracts using the hardhat scripts in the contracts/ folder. After having deployed them, be sure to update the  CHAIN_PROCESSOR_ADDRESS in the .env file with the address of the ProcessorEndpoint smart contract, and restart the docker compose.
