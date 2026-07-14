This folder contains an example docker compose to be used for debug/demo.<br>
It starts a dev chain using [Foundry Anvil](https://getfoundry.sh/anvil/overview), automatically deploys the smart contracts, and runs the processor manager and authority service.

## Instructions:

1)  (Skip this step if you want to use the Docker hub official images)
    To build locally the docker images needed by the docker compose run *from the project root folder* the following commands.

    Build the images:

    ```
    docker build -t horizen/cce-executor -f dockerfiles/executor/Dockerfile .
    docker build -t horizen/cce-manager -f dockerfiles/manager/Dockerfile .
    docker build -t horizen/kms-proxy -f dockerfiles/kms_proxy/Dockerfile .
    docker build -t horizen/cce-authorityservice -f dockerfiles/authorityservice/Dockerfile .
    docker build -t horizen/cce-chain -f dockerfiles/chain/Dockerfile .
    docker build -t horizen/cce-deployer -f dockerfiles/deployer/Dockerfile .
    docker build -t horizen/cce-subgraph-deployer -f dockerfiles/subgraph-deployer/Dockerfile .
    ```

    > The `cce-executor` command above is a convenient **dev** build (unpinned,
    > version auto-detected from `.git`). For a **reproducible** enclave image
    > whose `PCR0`/`PCR1`/`PCR2` can be verified against the on-chain value, use
    > `dockerfiles/executor/build-eif.sh` instead — see
    > [`docs/design/REPRODUCIBLE_EIF_BUILD.md`](../docs/design/REPRODUCIBLE_EIF_BUILD.md).

    > The executor image bakes the key-continuity guard
    > `EXECUTOR_EXPECT_EXISTING_KEYSET=true` (the production upgrade variant): with
    > the guard on, the Executor refuses to bootstrap a new keyset when the Manager
    > has no recovery data — which is exactly the situation on a first start. In this
    > dev compose the flag is forwarded from `.env` at runtime, and `.env.dev` sets
    > it to `false`, so a fresh environment bootstraps normally without rebuilding
    > the image. Make sure your `.env` contains `EXECUTOR_EXPECT_EXISTING_KEYSET=false`
    > if you created it before this variable was added.

2) Switch to "dockerfiles" folder

    ```
    cd dockerfiles
    ```

3) Create an .env file. For local development you can copy the ready-made dev defaults:

    ```
    cp .env.dev .env
    ```

    Or create one from scratch using `.env.template` as a reference.

4) Start the environment with:

    ```
    docker compose up
    ```

    The startup sequence is:
    1. **chain** (Anvil) starts and becomes available
    2. **subgraph-postgres**, **subgraph-ipfs** start (Graph Node infrastructure)
    3. **deployer** connects to the chain, deploys all smart contracts, writes the deployed addresses to a shared volume, and exits
    4. **subgraph-node** (Graph Node) starts, connects to the chain, and becomes healthy
    5. **subgraph-deployer** reads the deployed contract addresses, generates a local subgraph manifest, and deploys the subgraph to Graph Node, then exits
    6. **manager** and **authorityservice** start, reading the deployed contract addresses and querying the subgraph automatically

## Additional info:

- exposed addresses/ports (see `.env.template` for defaults):
- authorityservice: listens on `${AUTHORITY_SERVICE_LISTEN_ADDRESS}` on `${AUTHORITY_SERVICE_IP_ADDRESS}` inside the internal network.

- the manager database and chain data are persisted in docker volumes (`horizen-cce-manager-data` for the DB, `horizen-cce-chain-data` for chain data).<br>
  To start from scratch, delete the volumes.
- deployed contract addresses are stored in the `horizen-cce-deploy-data` volume. The deployer checks this on startup and skips deployment if contracts are already present on the chain.
- shared runtime files for manager and authorityservice are stored in `horizen-cce-shared-data`, with `reports/` for deanonymization outputs and `artifacts/` for uploaded deploy WASM blobs.
- to connect to the chain from Metamask, use the following parameters:
   - rpc url: http://localhost:8545
   - chainid: 31337

Authority service requires chain connectivity env vars (forwarded via docker-compose): `CHAIN_RPC_PROTOCOL`, `CHAIN_RPC_ADDRESS`, `CHAIN_RPC_PORT`, `CHAIN_PROCESSOR_ADDRESS`.
Authority service now reads events from the subgraph: set `AUTHORITY_SERVICE_SUBGRAPH_URL` (and keep chain RPC settings for chain ID checks).
For WASM deploy v1, ensure these are configured consistently in `.env`:
- `DEPLOYER_ADMIN`: this address is bootstrapped on-chain as the initial allowed deployer for `DEPLOYAPP`.
- `SHARED_DATA_FOLDER`: docker entrypoints derive `${SHARED_DATA_FOLDER}/reports` and `${SHARED_DATA_FOLDER}/artifacts` automatically for manager and authorityservice.
- `DEPLOY_ARTIFACTS_MAX_SIZE_MB`: optional upload limit (`0` means unlimited).

## Restarting and volume management

- **Restart without deleting volumes**: the deployer detects existing contracts and skips deployment. Fast restart.
- **Chain data deleted** (`docker volume rm dockerfiles_horizen-cce-chain-data`): the deployer detects contracts are missing from the chain and redeploys.
- **Deploy data deleted** (`docker volume rm dockerfiles_horizen-cce-deploy-data`): the deployer redeploys (same addresses since Anvil is deterministic with the same nonce).
- **Contracts modified**: rebuild the deployer image, delete both volumes, and restart.
- **Manager data deleted** (`docker volume rm dockerfiles_horizen-cce-manager-data`): the keyset recovery data is lost, so the Executor must bootstrap a fresh keyset on the next start. This works only with `EXECUTOR_EXPECT_EXISTING_KEYSET=false` in `.env` (the `.env.dev` default); with `true` the Executor aborts the handshake instead of generating a keyset — useful to reproduce the production upgrade guard, fatal on a first start. Also delete the chain volume: the on-chain `teeSigner` registered by the old keyset no longer matches the new one.

## Where to go next
The system is up and running, and you can deploy an app on-chain via the descriptor flow. Each deploy derives its own `applicationId` from the on-chain `requestId`, so there is no need to rename the WASM file or reserve a fixed id.

Practical how-to for the `vela-nova` test app (Private transfer):
- go to https://github.com/HorizenOfficial/vela-nova/releases/tag/v0.1.0
- download `payment_app.wasm`
- use the nova-linux wallet executable to submit the deploy and interact with the app:

    ```
    novaw deployapp --wasm /absolute/path/to/payment_app.wasm --max-value-fee "100 wei"
    ```

The wallet uploads the WASM to the authority service (`POST /deploy/upload`) and submits the on-chain deploy request; the Manager picks it up, forwards the artifact to the Executor, and the TEE verifies the WASM fingerprint against the on-chain descriptor before loading the module.

If you submit deploys from a different wallet, grant it first the `DEPLOYAPP` role with the ProcessorEndpoint management script (`contracts/scripts/management/addAllowedDeployer.ts`) using the admin account.

Use `wallet.conf.template` as wallet config file, with the following properties set to connect to this dev environment:

```
rpcUrl=http://localhost:8545
ProcessorAddress=0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9
TeeAuthenticatorAddress=0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
AuthorityServiceURL=http://localhost:8081
SubgraphURL=http://localhost:8000/subgraphs/name/hcce
```
