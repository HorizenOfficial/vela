# Guided installer

`vela-install.sh` walks step by step through preparing a Vela installation on a
server, in **without-TEE** mode (Nitro attestation disabled, executor running
with a fixed keyset).

The script runs on the **operator workstation**, not on the server: at the end
it produces a bundle (`docker-compose.yml` + `.env` + instructions) to upload
manually to the target machine.

Everything it generates is derived from the repository files at the selected
tag: no docker-compose or `.env` written from scratch.

## Prerequisites

| Requirement | Notes |
|---|---|
| bash, curl, tar, sed | present on any Linux/macOS |
| Node.js >= 20 and npm | needed for hardhat, for the subgraph and to generate the keys |
| docker | **optional**, only used to validate the generated compose locally |

The manual steps also need the [`graph`](https://thegraph.com/docs/en/subgraphs/developing/creating/install-the-cli/)
and [`goldsky`](https://docs.goldsky.com/subgraphs/deploying-subgraphs) CLIs: the
script checks whether they are present and, if not, prints the installation
commands.

## Usage

```bash
./install-helper/vela-install.sh
```

Options:

| Option | Effect |
|---|---|
| `--workdir DIR` | resume an installation already started in `DIR` |
| `--local PATH` | use a local repo checkout instead of downloading the tag from GitHub |
| `-h`, `--help` | show the help |

The script is resumable: its state lives in `DIR/state.env` and every completed
step is skipped. If you interrupt it (or a manual step goes wrong), run it again
with `--workdir`.

## What it asks at startup

1. **Version** to install (only `0.2.0` for now).
2. **TEE mode** (only *without TEE* for now).
3. **Target network**: `horizen-l3-testnet` or `base-sepolia`.
4. **Deployer Ethereum address**: it checks the format and that the address has
   funds on the selected network.
5. **Temporary working directory**: it must be empty.

## What it does

| Step | Automatic | Manual |
|---|---|---|
| 1. Download the tag sources | ✔ | |
| 2. `npm ci` in the contracts | ✔ | |
| 3. Generate the keys (manager, executor, state) | ✔ | |
| 4. Fund the manager address | check | send funds |
| 5. Deploy the contracts | writes `.env`, verifies the result on-chain | `npx hardhat run scripts/deploy/all.ts` |
| 6. Allowlist the ERC-20 token | patches the script, verifies on-chain | `npx hardhat run scripts/management/addAllowedToken.ts` |
| 7. Deploy the subgraph to Goldsky | patches `subgraph.yaml`, validates the endpoint | `graph build` + `goldsky subgraph deploy` |
| 8. Build the server bundle | ✔ | |

The manager and executor keys are randomly generated and never asked for. The
**deployer** private key, on the other hand, never goes through the script:
`contracts/.env` is prepared with an empty `PRIVATE_KEY` field, which the
operator fills in before running hardhat.

### Supported networks

| | `horizen-l3-testnet` | `base-sepolia` |
|---|---|---|
| chainId | 2651420 | 84532 |
| RPC | `https://horizen-testnet.rpc.caldera.xyz` | `https://sepolia.base.org` |
| Goldsky network | `horizen-testnet` | `base-sepolia` |
| suggested ZEN token | `0xb06EC4ce262D8dbDc24Fac87479A49A7DC4cFb87` | `0x107fdE93838e3404934877935993782F977324BB` |
| docker subnet | `10.10.40.0/24` | `10.20.40.0/24` |
| manager admin port | `4002` | `4102` |
| authority service port | `8081` | `8181` |

Everything that is scoped to the host differs per network — container, volume
and network names, the docker subnet, and the two published ports — so two
stacks can run side by side on one server. The subnet and the ports are asked in
the wizard, since a collision can only be detected on the target machine: docker
refuses overlapping address pools (`Pool overlaps with other one on this address
space`) and already-bound ports (`port is already allocated`).

## What it produces

In `DIR/out/deploy/`:

| File | Contents |
|---|---|
| `docker-compose.yml` | the compose from `dockerfiles/` without the local-only containers (`chain`, `deployer`, `subgraph-*`, `kms-proxy`), with the images pinned to the tag and names prefixed per network |
| `.env` | derived from `dockerfiles/.env.template`, keeping only the variables the generated compose actually uses |
| `README-DEPLOY.md` | addresses, upload instructions, ports, backup notes |

The compose transformations are guarded by assertions: if the upstream file
changes shape the script stops with an explicit error instead of silently
producing a wrong bundle.

## Security

- `DIR/secrets/keys.env` and `DIR/out/deploy/.env` contain **private keys**
  (mode `600`). Store them in a password manager and delete the working
  directory when you are done.
- The manager address must hold ETH on the selected network: it is the account
  that signs the state update transactions.
- The fixed executor keys correspond to the `TEE_SIGNER_ADDRESS` and
  `TEE_PUB_P521` registered on-chain: changing them makes already-stored
  encrypted state unreadable.

## Known limitations

- The `facilitator` service (x402 gasless payments) is **not** included: it is
  not part of this repository's docker-compose. The script reminds you of this
  in the final summary.
- Release tags older than the one that introduced a network in
  `contracts/hardhat.config.ts` do not contain it: in that case the script adds
  it to the downloaded copy and says so.
