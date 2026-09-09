#!/usr/bin/env bash
#
# vela-install.sh - Guided installer for Horizen Vela (non-TEE deployments).
#
# Runs on the operator workstation, NOT on the target server. It walks through
# key generation, the manual contract deploy, the ERC-20 allowlisting and the
# manual subgraph deploy, and finally produces a ready-to-upload bundle
# (docker-compose.yml + .env) in the chosen working directory.
#
# Everything it builds is derived from the upstream repository at the selected
# release tag: no hand-written compose or env files.
#
# Usage:
#   ./vela-install.sh                     start a new installation
#   ./vela-install.sh --workdir DIR       resume an installation started in DIR
#   ./vela-install.sh --local PATH        use a local repo checkout as source
#
set -euo pipefail

#-----------------------------------------------------------------------------
# Constants
#-----------------------------------------------------------------------------
readonly SCRIPT_NAME="$(basename "$0")"
readonly REPO_SLUG="HorizenOfficial/vela"
readonly SUPPORTED_VERSIONS=("0.2.0")
readonly SUPPORTED_TEE_MODES=("without TEE (no-tee)")
readonly SUPPORTED_NETWORKS=("horizen-l3-testnet" "base-sepolia")

# Services and volumes that only exist for the local dev stack and must not end
# up in a server deployment.
readonly DROP_SERVICES=(kms-proxy chain deployer subgraph-postgres subgraph-ipfs subgraph-node subgraph-deployer)
readonly KEEP_SERVICES=(executor manager authorityservice)
readonly DROP_VOLUMES=(horizen-cce-chain-data horizen-cce-deploy-data horizen-cce-subgraph-postgres horizen-cce-subgraph-ipfs)
readonly KEEP_VOLUMES=(horizen-cce-manager-data horizen-cce-shared-data horizen-cce-logs)
readonly UPSTREAM_NETWORK="pes_network"

#-----------------------------------------------------------------------------
# Output helpers
#-----------------------------------------------------------------------------
if [ -t 1 ] && command -v tput >/dev/null 2>&1 && [ "$(tput colors 2>/dev/null || echo 0)" -ge 8 ]; then
    C_RESET="$(tput sgr0)"; C_BOLD="$(tput bold)"
    C_RED="$(tput setaf 1)"; C_GREEN="$(tput setaf 2)"
    C_YELLOW="$(tput setaf 3)"; C_BLUE="$(tput setaf 4)"
else
    C_RESET=""; C_BOLD=""; C_RED=""; C_GREEN=""; C_YELLOW=""; C_BLUE=""
fi

info()  { printf '%s\n' "$*"; }
ok()    { printf '%s✓%s %s\n' "$C_GREEN" "$C_RESET" "$*"; }
warn()  { printf '%s!%s %s\n' "$C_YELLOW" "$C_RESET" "$*" >&2; }
err()   { printf '%s✗%s %s\n' "$C_RED" "$C_RESET" "$*" >&2; }
die()   { err "$*"; exit 1; }

step() {
    printf '\n%s%s──────────────────────────────────────────────────────────────%s\n' \
        "$C_BOLD" "$C_BLUE" "$C_RESET"
    printf '%s%s %s%s\n' "$C_BOLD" "$C_BLUE" "$*" "$C_RESET"
    printf '%s%s──────────────────────────────────────────────────────────────%s\n' \
        "$C_BOLD" "$C_BLUE" "$C_RESET"
}

cmd() { printf '    %s%s%s\n' "$C_BOLD" "$*" "$C_RESET"; }

pause_for_user() {
    local prompt="${1:-Press ENTER once you have completed the step above}"
    printf '\n%s%s...%s ' "$C_BOLD" "$prompt" "$C_RESET"
    read -r _ </dev/tty
    printf '\n'
}

ask() {
    # ask VARNAME "prompt" [default]
    local __var="$1" __prompt="$2" __default="${3:-}" __answer
    if [ -n "$__default" ]; then
        printf '%s [%s]: ' "$__prompt" "$__default" >&2
    else
        printf '%s: ' "$__prompt" >&2
    fi
    read -r __answer </dev/tty || die "input aborted"
    [ -z "$__answer" ] && __answer="$__default"
    printf -v "$__var" '%s' "$__answer"
}

ask_menu() {
    # ask_menu VARNAME "prompt" option...
    local __var="$1" __prompt="$2"; shift 2
    local __opts=("$@") __i __answer
    printf '\n%s\n' "$__prompt" >&2
    for __i in "${!__opts[@]}"; do
        printf '  %d) %s\n' "$((__i + 1))" "${__opts[$__i]}" >&2
    done
    while true; do
        printf 'Choice [1]: ' >&2
        read -r __answer </dev/tty || die "input aborted"
        [ -z "$__answer" ] && __answer=1
        if [[ "$__answer" =~ ^[0-9]+$ ]] && [ "$__answer" -ge 1 ] && [ "$__answer" -le "${#__opts[@]}" ]; then
            printf -v "$__var" '%s' "${__opts[$((__answer - 1))]}"
            return 0
        fi
        warn "Invalid choice."
    done
}

confirm() {
    local prompt="$1" answer
    while true; do
        printf '%s [y/n]: ' "$prompt" >&2
        read -r answer </dev/tty || die "input aborted"
        case "$answer" in
            [yY]|[yY][eE][sS]) return 0 ;;
            [nN]|[nN][oO]) return 1 ;;
            *) warn "Please answer 'y' or 'n'." ;;
        esac
    done
}

#-----------------------------------------------------------------------------
# Network profiles
#-----------------------------------------------------------------------------
apply_network_profile() {
    case "$1" in
        horizen-l3-testnet)
            NET_LABEL="Horizen L3 Testnet"
            NET_CHAIN_ID="2651420"
            NET_RPC_PROTOCOL="https"
            NET_RPC_HOST="horizen-testnet.rpc.caldera.xyz"
            NET_RPC_PORT="443"
            NET_RPC_URL="https://horizen-testnet.rpc.caldera.xyz"
            NET_HARDHAT="horizen-l3-testnet"
            NET_SUBGRAPH_SLUG="horizen-testnet"
            NET_ZEN_TOKEN="0xb06EC4ce262D8dbDc24Fac87479A49A7DC4cFb87"
            NET_EXPLORER="https://horizen-testnet.explorer.caldera.xyz"
            # Distinct per network so two stacks can coexist on one host.
            NET_SUBNET_PREFIX="10.10.40"
            NET_MANAGER_ADMIN_PORT="4002"
            NET_AUTHORITY_PORT="8081"
            ;;
        base-sepolia)
            NET_LABEL="Base Sepolia Testnet"
            NET_CHAIN_ID="84532"
            NET_RPC_PROTOCOL="https"
            NET_RPC_HOST="sepolia.base.org"
            NET_RPC_PORT="443"
            NET_RPC_URL="https://sepolia.base.org"
            NET_HARDHAT="base-sepolia"
            NET_SUBGRAPH_SLUG="base-sepolia"
            NET_ZEN_TOKEN="0x107fdE93838e3404934877935993782F977324BB"
            NET_EXPLORER="https://sepolia.basescan.org"
            NET_SUBNET_PREFIX="10.20.40"
            NET_MANAGER_ADMIN_PORT="4102"
            NET_AUTHORITY_PORT="8181"
            ;;
        *) die "unknown network: $1" ;;
    esac
}

#-----------------------------------------------------------------------------
# JSON-RPC helpers
#-----------------------------------------------------------------------------
rpc_result() {
    # rpc_result METHOD [PARAMS_JSON] -> prints the "result" string
    local method="$1" params="${2:-[]}" resp
    resp="$(curl -sS --max-time 25 -X POST "$NET_RPC_URL" \
        -H 'content-type: application/json' \
        -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"$method\",\"params\":$params}" 2>/dev/null)" \
        || { warn "RPC call $method failed (network unreachable?)"; return 1; }
    if printf '%s' "$resp" | grep -q '"error"'; then
        warn "RPC $method returned an error: $resp"
        return 1
    fi
    printf '%s' "$resp" | sed -n 's/.*"result"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p'
}

check_rpc_reachable() {
    local chain_id_hex chain_id_dec
    chain_id_hex="$(rpc_result eth_chainId)" || die "RPC $NET_RPC_URL is unreachable"
    [ -n "$chain_id_hex" ] || die "RPC $NET_RPC_URL did not return a chainId"
    chain_id_dec="$(printf '%d' "$chain_id_hex")"
    [ "$chain_id_dec" = "$NET_CHAIN_ID" ] \
        || die "unexpected chainId from $NET_RPC_URL: got $chain_id_dec, expected $NET_CHAIN_ID"
    ok "RPC reachable: $NET_RPC_URL (chainId $NET_CHAIN_ID)"
}

balance_wei_hex() {
    rpc_result eth_getBalance "[\"$1\",\"latest\"]"
}

format_eth() {
    # format_eth 0x<hex-wei> -> human readable, 6 decimals
    node -e '
        const wei = BigInt(process.argv[1]);
        const whole = wei / 10n ** 18n;
        const frac = (wei % 10n ** 18n).toString().padStart(18, "0").slice(0, 6);
        console.log(`${whole}.${frac}`);
    ' "$1"
}

has_code() {
    local code
    code="$(rpc_result eth_getCode "[\"$1\",\"latest\"]")" || return 1
    [ -n "$code" ] && [ "$code" != "0x" ]
}

is_address() {
    [[ "$1" =~ ^0x[0-9a-fA-F]{40}$ ]]
}

is_port() {
    [[ "$1" =~ ^[0-9]{1,5}$ ]] && [ "$1" -ge 1 ] && [ "$1" -le 65535 ]
}

# Accepts the first three octets of a /24, e.g. 10.20.40
is_subnet_prefix() {
    [[ "$1" =~ ^([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})$ ]] || return 1
    local octet
    for octet in "${BASH_REMATCH[@]:1:3}"; do
        [ "$octet" -le 255 ] || return 1
    done
    return 0
}

#-----------------------------------------------------------------------------
# State (resume support)
#-----------------------------------------------------------------------------
STATE_KEYS=(
    VELA_VERSION TEE_MODE NETWORK DEPLOYER_ADDRESS PREFIX COMPOSE_NETWORK_NAME
    MANAGER_ADDRESS TEE_SIGNER_ADDRESS TEE_PUB_P521
    CHAIN_PROCESSOR_ADDRESS CHAIN_TEEAUTHENTICATOR_ADDRESS CHAIN_TOKEN_ALLOWLIST_ADDRESS
    TOKEN_ADDRESS SUBGRAPH_START_BLOCK SUBGRAPH_NAME SUBGRAPH_VERSION SUBGRAPH_URL
    SUBNET_PREFIX MANAGER_ADMIN_PORT AUTHORITY_HOST_PORT MANAGER_FUNDED DONE_STEPS
)

state_save() {
    local key value tmp
    tmp="$(mktemp "${STATE_FILE}.XXXXXX")"
    {
        echo "# vela-install.sh state file - do not edit by hand"
        echo "# generated: $(date -u '+%Y-%m-%dT%H:%M:%SZ')"
        for key in "${STATE_KEYS[@]}"; do
            value="${!key:-}"
            case "$value" in
                *"'"*) die "value not allowed for $key (contains a single quote)" ;;
            esac
            printf "%s='%s'\n" "$key" "$value"
        done
    } >"$tmp"
    mv "$tmp" "$STATE_FILE"
    chmod 600 "$STATE_FILE"
}

state_load() {
    # shellcheck disable=SC1090
    set -a; . "$STATE_FILE"; set +a
}

mark_done() {
    case " ${DONE_STEPS:-} " in
        *" $1 "*) ;;
        *) DONE_STEPS="${DONE_STEPS:-} $1" ;;
    esac
    state_save
}

is_done() {
    case " ${DONE_STEPS:-} " in
        *" $1 "*) return 0 ;;
        *) return 1 ;;
    esac
}

#-----------------------------------------------------------------------------
# Prerequisites
#-----------------------------------------------------------------------------
check_prereqs() {
    local missing=() node_major
    for bin in curl tar node npm sed; do
        command -v "$bin" >/dev/null 2>&1 || missing+=("$bin")
    done
    if [ "${#missing[@]}" -gt 0 ]; then
        die "missing commands: ${missing[*]} - install them and run the script again"
    fi
    node_major="$(node -p 'process.versions.node.split(".")[0]')"
    [ "$node_major" -ge 20 ] || die "Node.js >= 20 is required (found $(node -v))"
    ok "prerequisites present (node $(node -v), npm $(npm -v))"
    if ! command -v docker >/dev/null 2>&1; then
        warn "docker not found: the generated docker-compose will not be validated locally"
    fi
}

#-----------------------------------------------------------------------------
# Wizard
#-----------------------------------------------------------------------------
run_wizard() {
    step "Installation setup"

    ask_menu VELA_VERSION "Which Vela version do you want to install?" "${SUPPORTED_VERSIONS[@]}"
    ok "version: $VELA_VERSION"

    ask_menu TEE_MODE "With or without TEE?" "${SUPPORTED_TEE_MODES[@]}"
    TEE_MODE="no-tee"
    ok "mode: without TEE (attestation disabled, fixed executor keys)"

    ask_menu NETWORK "Target network?" "${SUPPORTED_NETWORKS[@]}"
    apply_network_profile "$NETWORK"
    ok "network: $NET_LABEL (chainId $NET_CHAIN_ID)"
    check_rpc_reachable

    # Deployer address: format + funding check.
    while true; do
        ask DEPLOYER_ADDRESS "Ethereum address of the contract deployer"
        if ! is_address "$DEPLOYER_ADDRESS"; then
            warn "invalid format: expected 0x followed by 40 hex characters"
            continue
        fi
        local bal_hex bal_eth
        if ! bal_hex="$(balance_wei_hex "$DEPLOYER_ADDRESS")"; then
            warn "could not read the balance, try again"
            continue
        fi
        bal_eth="$(format_eth "$bal_hex")"
        if [ "$bal_hex" = "0x0" ]; then
            warn "address $DEPLOYER_ADDRESS has no funds on $NET_LABEL"
            info "  ETH is needed to pay the gas for the contract deploy."
            info "  Explorer: $NET_EXPLORER/address/$DEPLOYER_ADDRESS"
            confirm "Do you want to enter a different address?" && continue
            die "the deployer must have funds in order to deploy the contracts"
        fi
        ok "deployer: $DEPLOYER_ADDRESS (balance $bal_eth ETH)"
        break
    done

    # Working directory: must be empty.
    while true; do
        ask WORKDIR "Temporary working directory"
        [ -n "$WORKDIR" ] || { warn "empty path"; continue; }
        # Expand a leading ~ manually: read -r does not do it.
        case "$WORKDIR" in "~"/*) WORKDIR="$HOME/${WORKDIR#\~/}" ;; esac
        if [ -e "$WORKDIR" ] && [ ! -d "$WORKDIR" ]; then
            warn "$WORKDIR exists and is not a directory"
            continue
        fi
        if [ -d "$WORKDIR" ] && [ -n "$(ls -A "$WORKDIR" 2>/dev/null)" ]; then
            warn "$WORKDIR is not empty"
            info "  To resume an installation that was already started, use:"
            cmd "$SCRIPT_NAME --workdir $WORKDIR"
            continue
        fi
        mkdir -p "$WORKDIR" || { warn "could not create $WORKDIR"; continue; }
        WORKDIR="$(cd "$WORKDIR" && pwd)"
        ok "working directory: $WORKDIR"
        break
    done

    # The docker network range lives on the target server, so a collision only
    # shows up at 'docker compose up'. Let the operator pick a free one now.
    info ""
    info "The stack creates a private docker network on the server. If that range is"
    info "already taken (another stack, a VPN, the host LAN), docker refuses to start"
    info "with 'Pool overlaps with other one on this address space'."
    while true; do
        ask SUBNET_PREFIX "Docker network range (first three octets of a /24)" "$NET_SUBNET_PREFIX"
        if is_subnet_prefix "$SUBNET_PREFIX"; then
            ok "subnet: ${SUBNET_PREFIX}.0/24 (executor .10, manager .20, authority .40)"
            break
        fi
        warn "invalid format: expected three octets, e.g. $NET_SUBNET_PREFIX"
    done

    # Published host ports collide the same way the subnet does.
    info ""
    info "The stack publishes two ports on the server. They must be free there."
    while true; do
        ask MANAGER_ADMIN_PORT "Host port for the manager admin server" "$NET_MANAGER_ADMIN_PORT"
        is_port "$MANAGER_ADMIN_PORT" && break
        warn "invalid port: expected a number between 1 and 65535"
    done
    while true; do
        ask AUTHORITY_HOST_PORT "Host port for the authority service" "$NET_AUTHORITY_PORT"
        if ! is_port "$AUTHORITY_HOST_PORT"; then
            warn "invalid port: expected a number between 1 and 65535"
            continue
        fi
        [ "$AUTHORITY_HOST_PORT" = "$MANAGER_ADMIN_PORT" ] \
            && { warn "must differ from the manager admin port"; continue; }
        break
    done
    ok "host ports: manager admin $MANAGER_ADMIN_PORT, authority service $AUTHORITY_HOST_PORT"

    PREFIX="vela-${NETWORK}"
    COMPOSE_NETWORK_NAME="$(printf '%s' "$PREFIX" | tr '-' '_')_network"
    SUBGRAPH_NAME="vela-${NETWORK}"
    SUBGRAPH_VERSION="$VELA_VERSION"
    MANAGER_FUNDED="no"
    DONE_STEPS=""

    set_paths
    mkdir -p "$SECRETS_DIR" "$OUT_DIR"
    chmod 700 "$SECRETS_DIR"
    state_save

    info ""
    info "Docker resource names: containers/volumes prefixed with ${C_BOLD}${PREFIX}${C_RESET}, network ${C_BOLD}${COMPOSE_NETWORK_NAME}${C_RESET}"
    info "If the script is interrupted, resume it with:"
    cmd "$SCRIPT_NAME --workdir $WORKDIR"
}

set_paths() {
    STATE_FILE="$WORKDIR/state.env"
    SRC_DIR="$WORKDIR/src"
    SECRETS_DIR="$WORKDIR/secrets"
    KEYS_FILE="$SECRETS_DIR/keys.env"
    OUT_DIR="$WORKDIR/out"
    BUNDLE_DIR="$OUT_DIR/deploy"
    SUBGRAPH_DIR="$WORKDIR/subgraph"
    CONTRACTS_DIR="$SRC_DIR/contracts"
    DEPLOYED_ADDRESSES="$OUT_DIR/deployed_addresses.env"
}

#-----------------------------------------------------------------------------
# Step: fetch sources
#-----------------------------------------------------------------------------
step_fetch_sources() {
    step "1/7 - Downloading the sources (tag v${VELA_VERSION})"
    if is_done fetch; then
        ok "already done: sources in $SRC_DIR"
        return
    fi

    mkdir -p "$SRC_DIR"

    if [ -n "$LOCAL_REPO" ]; then
        [ -d "$LOCAL_REPO/contracts" ] && [ -d "$LOCAL_REPO/subgraphs/hcce" ] \
            || die "$LOCAL_REPO does not look like a vela repo checkout"
        info "Copying the sources from $LOCAL_REPO (--local mode)"
        mkdir -p "$SRC_DIR/dockerfiles" "$SRC_DIR/subgraphs"
        # Skip build artifacts and dependencies: they get regenerated locally.
        tar -C "$LOCAL_REPO" -cf - \
            --exclude=node_modules --exclude=artifacts --exclude=cache \
            --exclude=typechain-types --exclude=build --exclude=generated \
            contracts subgraphs/hcce | tar -C "$SRC_DIR" -xf -
        cp "$LOCAL_REPO/dockerfiles/docker-compose.yml" "$SRC_DIR/dockerfiles/"
        cp "$LOCAL_REPO/dockerfiles/.env.template" "$SRC_DIR/dockerfiles/"
        warn "sources taken from a local checkout: the bundle may not match tag v${VELA_VERSION}"
    else
        local url tarball
        url="https://codeload.github.com/${REPO_SLUG}/tar.gz/refs/tags/v${VELA_VERSION}"
        tarball="$WORKDIR/.vela-v${VELA_VERSION}.tar.gz"
        info "Downloading $url"
        curl -fsSL --max-time 180 -o "$tarball" "$url" \
            || die "download failed: check your connection and that tag v${VELA_VERSION} exists"
        # The archive root is vela-<version>/; strip it and keep only what we need.
        tar -C "$SRC_DIR" --strip-components=1 -xzf "$tarball" \
            "vela-${VELA_VERSION}/contracts" \
            "vela-${VELA_VERSION}/subgraphs/hcce" \
            "vela-${VELA_VERSION}/dockerfiles/docker-compose.yml" \
            "vela-${VELA_VERSION}/dockerfiles/.env.template" \
            || die "archive extraction failed"
        rm -f "$tarball"
    fi

    for f in contracts/hardhat.config.ts contracts/scripts/deploy/all.ts \
             contracts/scripts/management/addAllowedToken.ts \
             subgraphs/hcce/subgraph.yaml \
             dockerfiles/docker-compose.yml dockerfiles/.env.template; do
        [ -f "$SRC_DIR/$f" ] || die "expected file missing from the sources: $f"
    done

    ensure_hardhat_network

    ok "sources ready in $SRC_DIR"
    mark_done fetch
}

# Releases older than the one that introduced a given network do not define it
# in hardhat.config.ts. Add it rather than forcing the operator to patch by hand.
ensure_hardhat_network() {
    local config="$SRC_DIR/contracts/hardhat.config.ts"
    if grep -q "'${NET_HARDHAT}':" "$config"; then
        ok "network '${NET_HARDHAT}' already defined in hardhat.config.ts"
        return
    fi
    warn "network '${NET_HARDHAT}' is not defined in hardhat.config.ts of tag v${VELA_VERSION}: adding it"
    CFG_FILE="$config" CFG_NET="$NET_HARDHAT" CFG_URL="$NET_RPC_URL" node - <<'NODE' || die "could not add the network to hardhat.config.ts"
const fs = require('fs');
const file = process.env.CFG_FILE;
const src = fs.readFileSync(file, 'utf8');
const anchor = /^( *)local: \{$/m.exec(src);
if (!anchor) {
    console.error("hardhat.config.ts does not contain the expected 'local: {' anchor");
    process.exit(1);
}
const indent = anchor[1];
const block = [
    `${indent}'${process.env.CFG_NET}': {`,
    `${indent}  url: '${process.env.CFG_URL}',`,
    `${indent}  accounts,`,
    `${indent}},`,
    '',
].join('\n');
fs.writeFileSync(file, src.replace(anchor[0], block + anchor[0]));
NODE
    grep -q "'${NET_HARDHAT}':" "$config" \
        || die "adding network '${NET_HARDHAT}' to hardhat.config.ts did not succeed"
    ok "network '${NET_HARDHAT}' added to hardhat.config.ts"
}

#-----------------------------------------------------------------------------
# Step: install contract dependencies
#-----------------------------------------------------------------------------
step_npm_install() {
    step "2/7 - Installing the contract dependencies"
    if is_done npm; then
        ok "already done"
        return
    fi
    info "Running 'npm ci' in $CONTRACTS_DIR (needed for the deploy and to generate the keys)"
    ( cd "$CONTRACTS_DIR" && npm ci --no-audit --no-fund ) \
        || die "npm ci failed in $CONTRACTS_DIR"
    [ -d "$CONTRACTS_DIR/node_modules/ethers" ] \
        || die "ethers not found after npm ci"
    ok "dependencies installed"
    mark_done npm
}

# Run node with ethers resolvable from the contracts checkout.
node_with_ethers() {
    NODE_PATH="$CONTRACTS_DIR/node_modules" node "$@"
}

#-----------------------------------------------------------------------------
# Step: generate keys
#-----------------------------------------------------------------------------
step_generate_keys() {
    step "3/7 - Generating the keys"
    if is_done keys; then
        if [ -f "$KEYS_FILE" ]; then
            ok "already done: keys in $KEYS_FILE"
            return
        fi
        # Regenerating would invalidate anything already registered on-chain.
        if is_done contracts; then
            die "$KEYS_FILE is gone but the contracts were already deployed with those keys: restore them from a backup"
        fi
        warn "$KEYS_FILE not found: regenerating the keys"
    fi

    info "Generating random manager and executor keys..."
    local tmp
    tmp="$(mktemp "${SECRETS_DIR}/keys.XXXXXX")"
    chmod 600 "$tmp"

    node_with_ethers - >"$tmp" <<'NODE' || { rm -f "$tmp"; die "key generation failed"; }
const crypto = require('crypto');
const { ethers } = require('ethers');

// JWK EC coordinates are already left-padded, but be defensive.
const pad = (buf, n) => (buf.length === n ? buf : Buffer.concat([Buffer.alloc(n - buf.length), buf]));

// secp256k1 keys: the Go side imports them with hex.DecodeString, so no 0x prefix.
const newSecp = () => {
    const w = ethers.Wallet.createRandom();
    return { priv: w.privateKey.slice(2), addr: w.address };
};

const manager = newSecp();
const signing = newSecp();

// P-521 communication key: 66-byte scalar, 133-byte uncompressed public key.
const { privateKey } = crypto.generateKeyPairSync('ec', { namedCurve: 'secp521r1' });
const jwk = privateKey.export({ format: 'jwk' });
const d = pad(Buffer.from(jwk.d, 'base64url'), 66);
const x = pad(Buffer.from(jwk.x, 'base64url'), 66);
const y = pad(Buffer.from(jwk.y, 'base64url'), 66);
const commPriv = d.toString('hex');
const commPub = '04' + x.toString('hex') + y.toString('hex');

const stateKey = crypto.randomBytes(32).toString('hex');

const assertLen = (name, value, len) => {
    if (value.length !== len) throw new Error(`${name}: expected ${len} hex characters, got ${value.length}`);
    if (!/^[0-9a-f]+$/.test(value)) throw new Error(`${name}: not a lowercase hex string`);
};
assertLen('MANAGER_KEY_SECP256', manager.priv, 64);
assertLen('EXECUTOR_FIXED_SIGNING_KEY', signing.priv, 64);
assertLen('EXECUTOR_FIXED_COMMUNICATION_KEY', commPriv, 132);
assertLen('EXECUTOR_FIXED_STATE_KEY', stateKey, 64);
assertLen('TEE_PUB_P521', commPub, 266);

process.stdout.write([
    '# Keys generated by vela-install.sh - SECRET MATERIAL, do not commit',
    `MANAGER_KEY_SECP256='${manager.priv}'`,
    `MANAGER_ADDRESS='${manager.addr}'`,
    `EXECUTOR_FIXED_SIGNING_KEY='${signing.priv}'`,
    `TEE_SIGNER_ADDRESS='${signing.addr}'`,
    `EXECUTOR_FIXED_COMMUNICATION_KEY='${commPriv}'`,
    `TEE_PUB_P521='0x${commPub}'`,
    `EXECUTOR_FIXED_STATE_KEY='${stateKey}'`,
    '',
].join('\n'));
NODE

    mv "$tmp" "$KEYS_FILE"
    chmod 600 "$KEYS_FILE"

    # shellcheck disable=SC1090
    set -a; . "$KEYS_FILE"; set +a
    state_save

    ok "keys generated in $KEYS_FILE (mode 600)"
    info ""
    info "  Public parts:"
    info "    manager address (UPDATE_STATUS_OPERATOR):  $MANAGER_ADDRESS"
    info "    TEE signer address (TEE_SIGNER_ADDRESS):   $TEE_SIGNER_ADDRESS"
    warn "$SECRETS_DIR holds private keys: protect it and remove it when you are done"
    mark_done keys
}

load_keys() {
    [ -f "$KEYS_FILE" ] || die "key file not found: $KEYS_FILE"
    # shellcheck disable=SC1090
    set -a; . "$KEYS_FILE"; set +a
}

#-----------------------------------------------------------------------------
# Step: fund the manager account
#-----------------------------------------------------------------------------
step_fund_manager() {
    step "4/7 - Funding the manager address"
    if is_done fund; then
        ok "already done"
        return
    fi

    info "The manager signs the transactions that publish execution results"
    info "on-chain: without ETH the system stalls."
    info ""
    info "  Address to fund: ${C_BOLD}${MANAGER_ADDRESS}${C_RESET}"
    info "  Network:         $NET_LABEL"
    info "  Explorer:        $NET_EXPLORER/address/$MANAGER_ADDRESS"
    info ""
    info "Send funds to this address. Type 'skip' to continue anyway."

    local bal_hex bal_eth answer
    while true; do
        printf '\nPress ENTER to check the balance (or type skip): ' >&2
        read -r answer </dev/tty || die "input aborted"
        if [ "$answer" = "skip" ]; then
            MANAGER_FUNDED="no"
            warn "check skipped: remember to fund $MANAGER_ADDRESS before starting the containers"
            break
        fi
        if ! bal_hex="$(balance_wei_hex "$MANAGER_ADDRESS")"; then
            continue
        fi
        if [ "$bal_hex" != "0x0" ]; then
            bal_eth="$(format_eth "$bal_hex")"
            MANAGER_FUNDED="yes"
            ok "manager balance: $bal_eth ETH"
            break
        fi
        warn "balance is still zero"
    done
    mark_done fund
}

#-----------------------------------------------------------------------------
# Step: contract deploy (manual)
#-----------------------------------------------------------------------------
write_contracts_env() {
    local target="$CONTRACTS_DIR/.env"
    cat >"$target" <<EOF
# Generated by vela-install.sh - $(date -u '+%Y-%m-%dT%H:%M:%SZ')
# Vela ${VELA_VERSION} installation without TEE on ${NET_LABEL}.
#
# NOTE: fill in PRIVATE_KEY (or MNEMONIC) with the deployer private key before
# running hardhat. vela-install.sh never asks for it and never stores it.
#   expected deployer: ${DEPLOYER_ADDRESS}

NETWORK=${NET_HARDHAT}
PRIVATE_KEY=
MNEMONIC=

# Contract admin: manages ProcessorEndpoint and TokenAllowlist, and is the first
# address granted the DEPLOYAPP role.
ADMIN=${DEPLOYER_ADDRESS}
RESET_OPERATOR=${DEPLOYER_ADDRESS}

# Manager address, the only one allowed to publish state updates.
UPDATE_STATUS_OPERATOR=${MANAGER_ADDRESS}

# Without-TEE mode: NoAttestationTeeAuthenticator is deployed with the executor
# keys generated by this script.
TEE_NO_ATTESTATION=true
TEE_SIGNER_ADDRESS=${TEE_SIGNER_ADDRESS}
TEE_PUB_P521=${TEE_PUB_P521}
TEE_PCR0=
TEE_MAX_VERIFICATION_AGE=315360000

# Minimum fee per request (in wei). Must match EXECUTOR_MIN_FEE_PER_REQUEST in
# the server .env.
MIN_FEE_PER_REQUEST=10

# The deployed addresses are written here and read back by the script.
DEPLOY_OUTPUT_DIR=${OUT_DIR}
EOF
    chmod 600 "$target"
}

step_deploy_contracts() {
    step "5/7 - Contract deploy (manual)"
    if is_done contracts; then
        ok "already done"
        info "  ProcessorEndpoint:  $CHAIN_PROCESSOR_ADDRESS"
        info "  TeeAuthenticator:   $CHAIN_TEEAUTHENTICATOR_ADDRESS"
        info "  TokenAllowlist:     $CHAIN_TOKEN_ALLOWLIST_ADDRESS"
        return
    fi

    write_contracts_env
    ok "wrote $CONTRACTS_DIR/.env"

    # The subgraph must start indexing from before the contracts existed.
    if [ -z "${SUBGRAPH_START_BLOCK:-}" ]; then
        local blk_hex
        blk_hex="$(rpc_result eth_blockNumber)" || die "could not read the current block"
        SUBGRAPH_START_BLOCK="$(printf '%d' "$blk_hex")"
        state_save
        ok "subgraph start block: $SUBGRAPH_START_BLOCK"
    fi

    info ""
    info "Your turn. Open another terminal and run:"
    info ""
    cmd "cd $CONTRACTS_DIR"
    cmd "\$EDITOR .env      # set PRIVATE_KEY of deployer $DEPLOYER_ADDRESS"
    cmd "npx hardhat run scripts/deploy/all.ts"
    info ""
    info "The deploy script will write the addresses to:"
    info "  $DEPLOYED_ADDRESSES"
    pause_for_user "Press ENTER once the deploy has finished"

    [ -f "$DEPLOYED_ADDRESSES" ] \
        || die "$DEPLOYED_ADDRESSES not found: the deploy did not succeed"

    # shellcheck disable=SC1090
    set -a; . "$DEPLOYED_ADDRESSES"; set +a

    local addr name
    for name in CHAIN_PROCESSOR_ADDRESS CHAIN_TEEAUTHENTICATOR_ADDRESS CHAIN_TOKEN_ALLOWLIST_ADDRESS; do
        addr="${!name:-}"
        [ -n "$addr" ] || die "$name missing from $DEPLOYED_ADDRESSES"
        is_address "$addr" || die "$name is not a valid address: $addr"
        has_code "$addr" || die "$name ($addr) has no bytecode on $NET_LABEL: did you deploy to the right network?"
        ok "$name = $addr"
    done

    state_save
    mark_done contracts
}

#-----------------------------------------------------------------------------
# Step: ERC-20 allowlist (mandatory)
#-----------------------------------------------------------------------------
token_is_allowed() {
    NODE_PATH="$CONTRACTS_DIR/node_modules" \
    RPC_URL="$NET_RPC_URL" ALLOWLIST="$CHAIN_TOKEN_ALLOWLIST_ADDRESS" TOKEN="$1" \
    node - <<'NODE'
const { ethers } = require('ethers');
const provider = new ethers.JsonRpcProvider(process.env.RPC_URL);
const abi = ['function isAllowedToken(address) view returns (bool)'];
const c = new ethers.Contract(process.env.ALLOWLIST, abi, provider);
c.isAllowedToken(process.env.TOKEN)
    .then((allowed) => process.exit(allowed ? 0 : 1))
    .catch((e) => { console.error(String(e.message || e)); process.exit(2); });
NODE
}

patch_add_allowed_token_script() {
    local f="$CONTRACTS_DIR/scripts/management/addAllowedToken.ts"
    grep -q "^const TOKEN_ALLOWLIST = " "$f" || die "addAllowedToken.ts has an unexpected shape"
    grep -q "^const TOKEN_TO_ALLOW = " "$f" || die "addAllowedToken.ts has an unexpected shape"
    sed -i.bak \
        -e "s|^const TOKEN_ALLOWLIST = .*|const TOKEN_ALLOWLIST = '${CHAIN_TOKEN_ALLOWLIST_ADDRESS}';|" \
        -e "s|^const TOKEN_TO_ALLOW = .*|const TOKEN_TO_ALLOW = '${TOKEN_ADDRESS}';|" \
        "$f"
    rm -f "${f}.bak"
}

step_allowlist_token() {
    step "6/7 - ERC-20 token allowlisting"
    if is_done token; then
        ok "already done: $TOKEN_ADDRESS is allowlisted"
        return
    fi

    info "Without at least one allowlisted token nobody can deposit funds."
    info "The suggested address is ZEN on $NET_LABEL."
    info ""

    while true; do
        ask TOKEN_ADDRESS "Address of the ERC-20 token to allowlist" "$NET_ZEN_TOKEN"
        if ! is_address "$TOKEN_ADDRESS"; then
            warn "invalid format"
            continue
        fi
        if ! has_code "$TOKEN_ADDRESS"; then
            warn "$TOKEN_ADDRESS is not a contract on $NET_LABEL"
            continue
        fi
        break
    done
    state_save
    ok "token: $TOKEN_ADDRESS"

    if token_is_allowed "$TOKEN_ADDRESS"; then
        ok "the token is already allowlisted"
        mark_done token
        return
    fi

    patch_add_allowed_token_script
    ok "updated scripts/management/addAllowedToken.ts with the right addresses"

    info ""
    info "Run the allowlisting with the admin account ($DEPLOYER_ADDRESS):"
    info ""
    cmd "cd $CONTRACTS_DIR"
    cmd "npx hardhat run scripts/management/addAllowedToken.ts"
    pause_for_user "Press ENTER once the transaction is confirmed"

    local rc
    while true; do
        set +e
        token_is_allowed "$TOKEN_ADDRESS"; rc=$?
        set -e
        case "$rc" in
            0) ok "verified on-chain: $TOKEN_ADDRESS is allowlisted"; break ;;
            1) warn "the token is not allowlisted yet" ;;
            *) warn "check failed (node connectivity problem?)" ;;
        esac
        confirm "Retry the check?" || die "the token must be allowlisted to complete the installation"
    done

    mark_done token
}

#-----------------------------------------------------------------------------
# Step: subgraph (manual)
#-----------------------------------------------------------------------------
prepare_subgraph_dir() {
    rm -rf "$SUBGRAPH_DIR"
    mkdir -p "$SUBGRAPH_DIR"
    tar -C "$SRC_DIR/subgraphs/hcce" -cf - . | tar -C "$SUBGRAPH_DIR" -xf -

    local manifest="$SUBGRAPH_DIR/subgraph.yaml"
    grep -q '0x<processor_address>' "$manifest" \
        || die "subgraph.yaml does not contain the 0x<processor_address> placeholder"
    grep -q '0x<token_allowlist_address>' "$manifest" \
        || die "subgraph.yaml does not contain the 0x<token_allowlist_address> placeholder"

    sed -i.bak \
        -e "s|0x<processor_address>|${CHAIN_PROCESSOR_ADDRESS}|g" \
        -e "s|0x<token_allowlist_address>|${CHAIN_TOKEN_ALLOWLIST_ADDRESS}|g" \
        -e "s|^\( *network: \).*|\1${NET_SUBGRAPH_SLUG}|" \
        -e "s|^\( *startBlock: \).*|\1${SUBGRAPH_START_BLOCK}|" \
        "$manifest"
    rm -f "${manifest}.bak"

    grep -q "0x<" "$manifest" && die "placeholders are still present in subgraph.yaml"
    return 0
}

check_subgraph_cli() {
    local missing=0
    if ! command -v graph >/dev/null 2>&1; then
        warn "'graph' CLI not found. Install it with:"
        cmd "npm install -g @graphprotocol/graph-cli@latest"
        missing=1
    else
        ok "graph-cli present ($(graph --version 2>/dev/null | head -1))"
    fi
    if ! command -v goldsky >/dev/null 2>&1; then
        warn "'goldsky' CLI not found. Install it with:"
        cmd "curl https://goldsky.com | sh"
        missing=1
    else
        ok "goldsky CLI present"
        info "  if you have not logged in yet: goldsky login"
    fi
    return "$missing"
}

step_deploy_subgraph() {
    step "7/7 - Subgraph deploy on Goldsky (manual)"
    if is_done subgraph; then
        ok "already done: $SUBGRAPH_URL"
        return
    fi

    ask SUBGRAPH_NAME "Subgraph name on Goldsky" "${SUBGRAPH_NAME:-vela-${NETWORK}}"
    ask SUBGRAPH_VERSION "Subgraph version" "${SUBGRAPH_VERSION:-$VELA_VERSION}"
    state_save

    prepare_subgraph_dir
    ok "subgraph.yaml prepared in $SUBGRAPH_DIR"
    info "  network:    $NET_SUBGRAPH_SLUG"
    info "  processor:  $CHAIN_PROCESSOR_ADDRESS"
    info "  allowlist:  $CHAIN_TOKEN_ALLOWLIST_ADDRESS"
    info "  startBlock: $SUBGRAPH_START_BLOCK"
    info ""

    check_subgraph_cli || true

    info ""
    info "Build and deploy it:"
    info ""
    cmd "cd $SUBGRAPH_DIR"
    cmd "npm install"
    cmd "graph codegen"
    cmd "graph build"
    cmd "goldsky login    # if you have not already"
    cmd "goldsky subgraph deploy ${SUBGRAPH_NAME}/${SUBGRAPH_VERSION} --path ."
    cmd "goldsky subgraph tag create ${SUBGRAPH_NAME}/${SUBGRAPH_VERSION} --tag prod"
    info ""
    info "The 'prod' tag gives a stable endpoint across deploys:"
    info "  https://api.goldsky.com/api/public/<PROJECT_ID>/subgraphs/${SUBGRAPH_NAME}/prod/gn"
    pause_for_user "Press ENTER once the subgraph has been deployed"

    while true; do
        ask SUBGRAPH_URL "Paste the subgraph GraphQL endpoint"
        case "$SUBGRAPH_URL" in
            http://*|https://*) ;;
            *) warn "must be an http(s) URL"; continue ;;
        esac
        info "Checking the endpoint..."
        local resp
        resp="$(curl -sS --max-time 30 -X POST "$SUBGRAPH_URL" \
            -H 'content-type: application/json' \
            -d '{"query":"{ _meta { block { number } } }"}' 2>/dev/null)" || resp=""
        if printf '%s' "$resp" | grep -q '"number"'; then
            local indexed
            indexed="$(printf '%s' "$resp" | sed -n 's/.*"number"[[:space:]]*:[[:space:]]*\([0-9]*\).*/\1/p')"
            ok "subgraph reachable, last indexed block: ${indexed:-?}"
            break
        fi
        warn "the endpoint did not answer as expected: ${resp:-no response}"
        confirm "Do you want to enter a different URL?" && continue
        confirm "Continue anyway with this URL?" && break
    done

    state_save
    mark_done subgraph
}

#-----------------------------------------------------------------------------
# Step: build the server bundle
#-----------------------------------------------------------------------------
build_compose() {
    CFG_VERSION="$VELA_VERSION" CFG_PREFIX="$PREFIX" \
    CFG_NETWORK_NAME="$COMPOSE_NETWORK_NAME" \
    CFG_DROP_SERVICES="${DROP_SERVICES[*]}" CFG_KEEP_SERVICES="${KEEP_SERVICES[*]}" \
    CFG_DROP_VOLUMES="${DROP_VOLUMES[*]}" CFG_KEEP_VOLUMES="${KEEP_VOLUMES[*]}" \
    CFG_UPSTREAM_NETWORK="$UPSTREAM_NETWORK" \
    CFG_AUTHORITY_PORT="${AUTHORITY_HOST_PORT:-$NET_AUTHORITY_PORT}" \
    CFG_SRC_COMPOSE="$SRC_DIR/dockerfiles/docker-compose.yml" \
    CFG_OUT_COMPOSE="$BUNDLE_DIR/docker-compose.yml" \
    node - <<'NODE'
const fs = require('fs');

const env = process.env;
const dropServices = new Set(env.CFG_DROP_SERVICES.split(/\s+/).filter(Boolean));
const keepServices = env.CFG_KEEP_SERVICES.split(/\s+/).filter(Boolean);
const dropVolumes = new Set(env.CFG_DROP_VOLUMES.split(/\s+/).filter(Boolean));
const keepVolumes = env.CFG_KEEP_VOLUMES.split(/\s+/).filter(Boolean);
const prefix = env.CFG_PREFIX;
const version = env.CFG_VERSION;
const netName = env.CFG_NETWORK_NAME;
const upstreamNet = env.CFG_UPSTREAM_NETWORK;
const authorityPort = env.CFG_AUTHORITY_PORT;

const fail = (msg) => { console.error(`unexpected upstream docker-compose: ${msg}`); process.exit(1); };

const src = fs.readFileSync(env.CFG_SRC_COMPOSE, 'utf8');
const lines = src.split('\n');

// Split the file into top-level sections keyed by their first-column name.
const sections = [];
let current = null;
for (const line of lines) {
    const m = /^([A-Za-z_][\w-]*):\s*$/.exec(line);
    if (m) {
        current = { name: m[1], header: line, body: [] };
        sections.push(current);
    } else if (current) {
        current.body.push(line);
    } else if (line.trim() !== '') {
        fail(`unexpected content before the first section: ${line}`);
    }
}
const section = (name) => sections.find((s) => s.name === name) || fail(`section '${name}' is missing`);

// Rename a volume: horizen-cce-manager-data -> <prefix>-manager-data
const renameVolume = (name) => `${prefix}-${name.replace(/^horizen-cce-/, '')}`;
// Rename a container: horizen-cce-manager -> <prefix>-manager
const renameContainer = (name) => `${prefix}-${name.replace(/^horizen-cce-/, '')}`;

// --- services -------------------------------------------------------------
const svcSection = section('services');
const blocks = [];
let block = null;
for (const line of svcSection.body) {
    const m = /^ {2}([A-Za-z_][\w.-]*):\s*$/.exec(line);
    if (m) {
        block = { name: m[1], lines: [line] };
        blocks.push(block);
    } else if (block) {
        block.lines.push(line);
    } else if (line.trim() !== '') {
        fail(`unexpected line in services: ${line}`);
    }
}

const found = blocks.map((b) => b.name);
for (const name of keepServices) {
    if (!found.includes(name)) fail(`service '${name}' not found`);
}
for (const name of found) {
    if (!keepServices.includes(name) && !dropServices.has(name)) {
        fail(`unknown service '${name}': vela-install.sh needs updating`);
    }
}

// Drop a "key:" mapping and all of its children from a block.
const dropKeyBlock = (blockLines, indent, key) => {
    const out = [];
    let skipping = false;
    for (const line of blockLines) {
        if (skipping) {
            const isChild = line.trim() === '' || /^\s*$/.test(line)
                ? false
                : line.search(/\S/) > indent;
            if (isChild) continue;
            skipping = false;
        }
        if (line === `${' '.repeat(indent)}${key}:`) { skipping = true; continue; }
        out.push(line);
    }
    return out;
};

// Filter the children of a "key:" mapping, dropping the key when nothing is left.
const filterKeyBlock = (blockLines, indent, key, keepEntry) => {
    const startIdx = blockLines.findIndex((l) => l === `${' '.repeat(indent)}${key}:`);
    if (startIdx === -1) return blockLines;
    let endIdx = startIdx + 1;
    while (endIdx < blockLines.length) {
        const line = blockLines[endIdx];
        if (line.trim() !== '' && line.search(/\S/) <= indent) break;
        endIdx += 1;
    }
    const children = blockLines.slice(startIdx + 1, endIdx);

    // Group children into entries: a list item, or a mapping key with its children.
    const entries = [];
    let entry = null;
    for (const line of children) {
        if (line.trim() === '') { if (entry) entry.lines.push(line); continue; }
        const listItem = /^\s*-\s*([\w.-]+)/.exec(line);
        const mapKey = /^\s*([\w.-]+):\s*$/.exec(line);
        const childIndent = line.search(/\S/);
        if (entry && childIndent > entry.indent) { entry.lines.push(line); continue; }
        if (listItem) entry = { name: listItem[1], indent: childIndent, lines: [line] };
        else if (mapKey) entry = { name: mapKey[1], indent: childIndent, lines: [line] };
        else fail(`unrecognized entry under ${key}: ${line}`);
        entries.push(entry);
    }

    const kept = entries.filter((e) => keepEntry(e.name)).flatMap((e) => e.lines);
    const head = blockLines.slice(0, startIdx);
    const tail = blockLines.slice(endIdx);
    if (kept.length === 0) return [...head, ...tail];
    return [...head, blockLines[startIdx], ...kept, ...tail];
};

const rendered = [];
for (const b of blocks) {
    if (dropServices.has(b.name)) continue;
    let bl = b.lines;

    // depends_on entries pointing at dropped services.
    bl = filterKeyBlock(bl, 4, 'depends_on', (name) => !dropServices.has(name));
    // volume mounts of dropped volumes (e.g. horizen-cce-deploy-data:/deploy-data).
    bl = filterKeyBlock(bl, 4, 'volumes', (name) => !dropVolumes.has(name));
    // "profiles: [...]" only exists on kms-proxy, but drop defensively.
    bl = dropKeyBlock(bl, 4, 'profiles');

    const text = bl.join('\n')
        // Pin the image to the selected release.
        .replace(/^(\s*image:\s*)(horizen\/[\w.-]+)\s*$/gm, (_, p, img) => `${p}${img}:v${version}`)
        .replace(/^(\s*container_name:\s*)([\w.-]+)\s*$/gm, (_, p, name) => `${p}${renameContainer(name)}`)
        .replace(new RegExp(`\\b${upstreamNet}\\b`, 'g'), netName)
        .replace(/\bhorizen-cce-([\w-]+)\b/g, (_, rest) => `${prefix}-${rest}`)
        // Host side of the published port becomes configurable; the container
        // side stays as upstream (the service listens on 8081).
        .replace(/^(\s*- ")8081:8081("\s*)$/gm, `$1\${AUTHORITY_HOST_PORT:-${authorityPort}}:8081$2`);
    rendered.push(text);
}

// --- top-level volumes ----------------------------------------------------
const volSection = section('volumes');
const volLines = volSection.body
    .filter((l) => l.trim() !== '')
    .map((l) => {
        const m = /^ {2}([\w.-]+):\s*$/.exec(l);
        if (!m) fail(`unexpected line in volumes: ${l}`);
        return m[1];
    });
for (const v of keepVolumes) {
    if (!volLines.includes(v)) fail(`volume '${v}' not found`);
}
for (const v of volLines) {
    if (!keepVolumes.includes(v) && !dropVolumes.has(v)) fail(`unknown volume '${v}'`);
}
const volumesOut = keepVolumes.map((v) => `  ${renameVolume(v)}:`).join('\n');

// --- top-level networks ---------------------------------------------------
const netSection = section('networks');
if (!netSection.body.some((l) => l.includes(`${upstreamNet}:`))) {
    fail(`network '${upstreamNet}' not found`);
}
const networksOut = netSection.body
    .join('\n')
    .replace(new RegExp(`\\b${upstreamNet}\\b`, 'g'), netName)
    .replace(/\s+$/, '');

const header = [
    `# Generated by install-helper/vela-install.sh from`,
    `# https://github.com/HorizenOfficial/vela/blob/v${version}/dockerfiles/docker-compose.yml`,
    `#`,
    `# Compared to the original, the local-development-only services have been`,
    `# removed (${[...dropServices].join(', ')}),`,
    `# the images are pinned to tag v${version}, and container, volume and network`,
    `# names use the '${prefix}' prefix.`,
    ``,
].join('\n');

// The body is what gets checked: the header intentionally names the services
// that were removed.
const body = [
    'networks:',
    networksOut,
    '',
    'volumes:',
    volumesOut,
    '',
    'services:',
    rendered.join('\n').replace(/\n{3,}/g, '\n\n').replace(/\s+$/, ''),
    '',
].join('\n');

// --- assertions on the produced file --------------------------------------
for (const name of dropServices) {
    if (new RegExp(`(^|[^\\w-])${name}([^\\w-]|$)`).test(body)) {
        fail(`removed service '${name}' still appears in the output`);
    }
}
for (const v of dropVolumes) {
    if (body.includes(v)) fail(`removed volume '${v}' still appears in the output`);
}
if (body.includes(upstreamNet)) fail(`network '${upstreamNet}' still appears in the output`);
if (body.includes('/deploy-data')) fail(`a reference to /deploy-data is left in the output`);
if (body.includes('horizen-cce-')) fail(`an unrenamed 'horizen-cce-*' name is left in the output`);
for (const name of keepServices) {
    if (!new RegExp(`^  ${name}:$`, 'm').test(body)) fail(`service '${name}' disappeared from the output`);
}
if (!body.includes('${AUTHORITY_HOST_PORT:-')) {
    fail(`the authorityservice published port was not parameterized (upstream '8081:8081' line changed?)`);
}
const publishedPorts = body.match(/^\s*- "[^"]*:[^"]*"\s*$/gm) || [];
for (const line of publishedPorts) {
    if (!line.includes('${')) fail(`published port with a fixed host side: ${line.trim()}`);
}
const images = body.match(/^\s*image:.*$/gm) || [];
if (images.length !== keepServices.length) fail(`expected ${keepServices.length} images, found ${images.length}`);
for (const img of images) {
    if (!img.includes(`:v${version}`)) fail(`image without a version tag: ${img.trim()}`);
}

fs.writeFileSync(env.CFG_OUT_COMPOSE, header + body);
console.log(`services: ${keepServices.join(', ')}`);
NODE
}

build_env() {
    VELA_VERSION="$VELA_VERSION" NET_LABEL="$NET_LABEL" NETWORK="$NETWORK" \
    TEMPLATE="$SRC_DIR/dockerfiles/.env.template" \
    COMPOSE="$BUNDLE_DIR/docker-compose.yml" \
    OUT_ENV="$BUNDLE_DIR/.env" \
    OV_CHAIN_ID="$NET_CHAIN_ID" OV_SUBNET_PREFIX="${SUBNET_PREFIX:-$NET_SUBNET_PREFIX}" \
    OV_ADMIN_PORT="${MANAGER_ADMIN_PORT:-$NET_MANAGER_ADMIN_PORT}" \
    OV_AUTHORITY_PORT="${AUTHORITY_HOST_PORT:-$NET_AUTHORITY_PORT}" \
    OV_RPC_PROTOCOL="$NET_RPC_PROTOCOL" OV_RPC_HOST="$NET_RPC_HOST" OV_RPC_PORT="$NET_RPC_PORT" \
    OV_PROCESSOR="$CHAIN_PROCESSOR_ADDRESS" OV_TEEAUTH="$CHAIN_TEEAUTHENTICATOR_ADDRESS" \
    OV_SUBGRAPH_URL="$SUBGRAPH_URL" \
    OV_MANAGER_KEY="$MANAGER_KEY_SECP256" \
    OV_SIGNING_KEY="$EXECUTOR_FIXED_SIGNING_KEY" \
    OV_COMM_KEY="$EXECUTOR_FIXED_COMMUNICATION_KEY" \
    OV_STATE_KEY="$EXECUTOR_FIXED_STATE_KEY" \
    node - <<'NODE'
const fs = require('fs');
const env = process.env;

const template = fs.readFileSync(env.TEMPLATE, 'utf8');
const compose = fs.readFileSync(env.COMPOSE, 'utf8');

// Every variable the produced compose actually reads: ${VAR}, ${VAR:-x} and
// bare "- VAR" entries in the environment lists.
const referenced = new Set();
for (const m of compose.matchAll(/\$\{([A-Za-z_][A-Za-z0-9_]*)/g)) referenced.add(m[1]);
for (const m of compose.matchAll(/^\s*-\s+([A-Z][A-Z0-9_]*)\s*$/gm)) referenced.add(m[1]);
// COMPOSE_PROFILES is read by docker compose itself, not by the file.
referenced.add('COMPOSE_PROFILES');

const overrides = {
    // No kms-proxy in this deployment.
    COMPOSE_PROFILES: '',
    CHANNEL_TYPE: 'tcp',
    // Fixed keyset: no KMS-based recovery.
    EXECUTOR_KEYSET_RECOVERY_TYPE: '0',
    EXECUTOR_KMS_KEY_ARN: '',
    EXECUTOR_KMS_REGION: '',
    EXECUTOR_FIXED_SIGNING_KEY: env.OV_SIGNING_KEY,
    EXECUTOR_FIXED_COMMUNICATION_KEY: env.OV_COMM_KEY,
    EXECUTOR_FIXED_STATE_KEY: env.OV_STATE_KEY,
    MANAGER_KEY_SECP256: env.OV_MANAGER_KEY,
    // The template default ('./manager_data') is a host path; inside the
    // container the data must live on the mounted volume.
    MANAGER_DATA_FOLDER: '/data',
    SHARED_DATA_FOLDER: '/shared-data',
    // Log server file lives on the logs volume, which the entrypoint chowns.
    LOG_SERVER_FILE_NAME: '/logs/log_server.log',
    LOG_SERVER_FILE_ROTATION: 'true',
    CHAIN_RPC_PROTOCOL: env.OV_RPC_PROTOCOL,
    CHAIN_RPC_ADDRESS: env.OV_RPC_HOST,
    CHAIN_RPC_PORT: env.OV_RPC_PORT,
    CHAIN_PROCESSOR_ADDRESS: env.OV_PROCESSOR,
    CHAIN_TEEAUTHENTICATOR_ADDRESS: env.OV_TEEAUTH,
    CHAIN_ID: env.OV_CHAIN_ID,
    AUTHORITY_SERVICE_SUBGRAPH_URL: env.OV_SUBGRAPH_URL,
    // Per-network subnet: docker refuses two networks with overlapping pools,
    // so a second Vela stack on the same host needs its own range. Host octets
    // keep the template layout (.10 executor, .20 manager, .40 authority).
    // Published host ports: they must be free on the target server, and must
    // differ between two stacks running on the same host.
    MANAGER_ADMIN_PORT: env.OV_ADMIN_PORT,
    AUTHORITY_HOST_PORT: env.OV_AUTHORITY_PORT,
    INTERNAL_NETWORK_SUBNET: `${env.OV_SUBNET_PREFIX}.0/24`,
    EXECUTOR_IP_HOST: `${env.OV_SUBNET_PREFIX}.10`,
    MANAGER_IP_HOST: `${env.OV_SUBNET_PREFIX}.20`,
    AUTHORITY_SERVICE_IP_ADDRESS: `${env.OV_SUBNET_PREFIX}.40`,
};

// Some template comments describe the local dev stack and would be misleading
// in a server bundle: replace them.
const commentOverrides = {
    COMPOSE_PROFILES: [
        '# No docker compose profile active: this installation does not use the',
        '# kms-proxy (EXECUTOR_KEYSET_RECOVERY_TYPE=0, fixed keyset).',
    ],
    EXECUTOR_KMS_KEY_ARN: [
        '# KMS parameters unused with EXECUTOR_KEYSET_RECOVERY_TYPE=0.',
    ],
    LOG_SERVER_IP_HOST: [
        '# The log server runs inside the manager: this is the address the executor',
        '# sends its logs to.',
    ],
    LOG_SERVER_FILE_NAME: [
        '# Log server file, on the logs volume mounted in the manager.',
    ],
    MANAGER_ADMIN_PORT: [
        '# Admin server port, published on the host (admincli connects here).',
        '# Must be free on the server and distinct per stack.',
    ],
    EXECUTOR_FIXED_SIGNING_KEY: [
        '# Fixed executor keyset, randomly generated by vela-install.sh.',
        '# In without-TEE mode the executor uses these keys instead of generating',
        '# new ones at every start: they must stay consistent with TEE_SIGNER_ADDRESS',
        '# and TEE_PUB_P521 registered on-chain on NoAttestationTeeAuthenticator.',
        '# Changing them makes already-stored encrypted state unreadable.',
    ],
};

// The executor reaches the log server on the manager, not on localhost.
overrides.LOG_SERVER_IP_HOST = overrides.MANAGER_IP_HOST;

const quote = (v) => `'${String(v).replace(/'/g, "'\\''")}'`;

const out = [];
const emitted = new Set();
let buffer = [];

for (const line of template.split('\n')) {
    const m = /^([A-Za-z_][A-Za-z0-9_]*)=(.*)$/.exec(line);
    if (!m) { buffer.push(line); continue; }
    const key = m[1];
    if (!referenced.has(key)) { buffer = []; continue; }  // drop the key and its comments
    if (key in commentOverrides) buffer = [...commentOverrides[key]];
    // Trim leading blank lines from the buffered comment block.
    while (buffer.length && buffer[0].trim() === '') buffer.shift();
    if (buffer.length && out.length) out.push('');
    out.push(...buffer);
    buffer = [];
    out.push(key in overrides ? `${key}=${quote(overrides[key])}` : line);
    emitted.add(key);
}

const extras = Object.keys(overrides).filter((k) => !emitted.has(k) && referenced.has(k));
if (extras.length) {
    out.push('', '# Variables not present in .env.template.',
             '# AUTHORITY_HOST_PORT is the host side of the authority service port',
             '# mapping; the container keeps listening on 8081.');
    for (const k of extras) out.push(`${k}=${quote(overrides[k])}`);
}

const header = [
    `# Generated by install-helper/vela-install.sh from`,
    `# https://github.com/HorizenOfficial/vela/blob/v${env.VELA_VERSION}/dockerfiles/.env.template`,
    `#`,
    `# Vela ${env.VELA_VERSION} without TEE on ${env.NET_LABEL}.`,
    `# The variables used only by the local development containers (chain,`,
    `# deployer, subgraph, kms-proxy) have been removed.`,
    `#`,
    `# WARNING: this file contains private keys. Keep it mode 600, do not commit it.`,
    ``,
    ``,
].join('\n');

// Sanity checks on the result.
const result = header + out.join('\n').replace(/\n{3,}/g, '\n\n').replace(/\s+$/, '') + '\n';
const missing = [...referenced].filter((k) => !new RegExp(`^${k}=`, 'm').test(result));
if (missing.length) {
    console.error(`variables used by the compose but missing from the .env: ${missing.join(', ')}`);
    process.exit(1);
}
for (const k of ['MANAGER_KEY_SECP256', 'CHAIN_PROCESSOR_ADDRESS', 'AUTHORITY_SERVICE_SUBGRAPH_URL']) {
    if (new RegExp(`^${k}=''$`, 'm').test(result)) { console.error(`${k} is empty`); process.exit(1); }
}

fs.writeFileSync(env.OUT_ENV, result);
console.log(`variables written: ${referenced.size}`);
NODE
}

write_deploy_readme() {
    cat >"$BUNDLE_DIR/README-DEPLOY.md" <<EOF
# Vela ${VELA_VERSION} - ${NET_LABEL}

Bundle generated by \`install-helper/vela-install.sh\` on $(date -u '+%Y-%m-%d %H:%M UTC').
**Without-TEE** configuration: Nitro attestation is disabled and the executor
uses a fixed keyset passed through \`.env\`.

## Contents

| File | Description |
|---|---|
| \`docker-compose.yml\` | runtime stack (executor, manager, authorityservice) |
| \`.env\` | full configuration, **contains private keys** |

## Addresses

| | |
|---|---|
| Network | ${NET_LABEL} (chainId ${NET_CHAIN_ID}) |
| RPC | ${NET_RPC_URL} |
| ProcessorEndpoint | \`${CHAIN_PROCESSOR_ADDRESS}\` |
| TeeAuthenticator | \`${CHAIN_TEEAUTHENTICATOR_ADDRESS}\` |
| TokenAllowlist | \`${CHAIN_TOKEN_ALLOWLIST_ADDRESS}\` |
| Allowlisted token | \`${TOKEN_ADDRESS}\` |
| Admin / deployer | \`${DEPLOYER_ADDRESS}\` |
| Manager (update status operator) | \`${MANAGER_ADDRESS}\` |
| TEE signer | \`${TEE_SIGNER_ADDRESS}\` |
| Subgraph | ${SUBGRAPH_URL} |

## Installing on the server

1. Copy \`docker-compose.yml\` and \`.env\` to a dedicated directory on the server:

   \`\`\`bash
   scp docker-compose.yml .env user@server:/opt/vela/
   \`\`\`

2. Secure the configuration file (it contains private keys):

   \`\`\`bash
   chmod 600 /opt/vela/.env
   \`\`\`

3. Start the stack:

   \`\`\`bash
   cd /opt/vela
   docker compose up -d
   docker compose logs -f
   \`\`\`

## Exposed ports

| Host port | Service | Notes |
|---|---|---|
| ${MANAGER_ADMIN_PORT} | manager (admin server) | used by \`admincli\`, **do not expose to the internet** |
| ${AUTHORITY_HOST_PORT} | authorityservice | public API for deanonymization reports (container listens on 8081) |

Both must be free on this host. If one is taken, \`docker compose up\` fails with
\`Bind for 0.0.0.0:<port> failed: port is already allocated\`. Check with
\`ss -lptn 'sport = :${MANAGER_ADMIN_PORT}'\` and change \`MANAGER_ADMIN_PORT\` or
\`AUTHORITY_HOST_PORT\` in \`.env\` — no other file needs editing.

Note that \`MANAGER_ADMIN_PORT\` changes the port on both sides of the mapping, so
point \`admincli\` at the value set here.

## Network range

The stack creates a private docker network on \`${SUBNET_PREFIX}.0/24\`. If that
range is already in use on this host, \`docker compose up\` fails with:

\`\`\`
failed to create network ...: Error response from daemon:
Pool overlaps with other one on this address space
\`\`\`

Find what is using it, checking both docker networks and host routes (a VPN or
the LAN can occupy a 10.x range too):

\`\`\`bash
docker network ls -q | xargs -r docker network inspect \\
  --format '{{.Name}}: {{range .IPAM.Config}}{{.Subnet}} {{end}}'
ip route | grep '^10\\.'
\`\`\`

To move the stack to a free range, change the four values together in \`.env\`
(they must all share the same first three octets):

\`\`\`bash
sed -i 's/${SUBNET_PREFIX//./\\.}\\./10.90.40./g' .env
\`\`\`

They are \`INTERNAL_NETWORK_SUBNET\`, \`EXECUTOR_IP_HOST\`, \`MANAGER_IP_HOST\`
(also used by \`LOG_SERVER_IP_HOST\`) and \`AUTHORITY_SERVICE_IP_ADDRESS\`.

## Operational prerequisites

- The manager address \`${MANAGER_ADDRESS}\` must hold ETH on ${NET_LABEL}:
  without funds it cannot publish results on-chain.
- To let other wallets deploy apps, grant them the \`DEPLOYAPP\` role with
  \`contracts/scripts/management/addAllowedDeployer.ts\` using the admin account.

## Backup

Persistent data lives in docker volumes:

- \`${PREFIX}-manager-data\` - manager database (**backup is essential**)
- \`${PREFIX}-shared-data\` - reports and WASM artifacts
- \`${PREFIX}-logs\` - logs

## Facilitator

This bundle does **not** include the \`facilitator\` service (x402 gasless
payments), which is not part of the upstream docker-compose of this repository.
If you need it, add it manually using the
\`horizen/vela-facilitator:v${VELA_VERSION}\` image.
EOF
}

step_build_bundle() {
    step "Building the server bundle"
    mkdir -p "$BUNDLE_DIR"

    info "Building docker-compose.yml from the upstream compose..."
    build_compose || die "docker-compose generation failed"
    ok "$BUNDLE_DIR/docker-compose.yml"

    info "Building .env from .env.template..."
    build_env || die ".env generation failed"
    chmod 600 "$BUNDLE_DIR/.env"
    ok "$BUNDLE_DIR/.env (mode 600)"

    write_deploy_readme
    ok "$BUNDLE_DIR/README-DEPLOY.md"

    if command -v docker >/dev/null 2>&1; then
        if ( cd "$BUNDLE_DIR" && docker compose config -q >/dev/null 2>&1 ); then
            ok "docker compose config: syntax valid"
        else
            warn "'docker compose config' reported problems:"
            ( cd "$BUNDLE_DIR" && docker compose config -q ) || true
        fi
    fi

    mark_done bundle
}

#-----------------------------------------------------------------------------
# Final summary
#-----------------------------------------------------------------------------
final_summary() {
    step "Installation prepared"

    info "Bundle to upload to the server: ${C_BOLD}${BUNDLE_DIR}${C_RESET}"
    info "    docker-compose.yml"
    info "    .env"
    info "    README-DEPLOY.md"
    info ""
    info "On the server:"
    cmd "scp $BUNDLE_DIR/docker-compose.yml $BUNDLE_DIR/.env user@server:/opt/vela/"
    cmd "ssh user@server 'chmod 600 /opt/vela/.env && cd /opt/vela && docker compose up -d'"
    info ""
    info "Summary:"
    info "  Network              $NET_LABEL (chainId $NET_CHAIN_ID)"
    info "  ProcessorEndpoint    $CHAIN_PROCESSOR_ADDRESS"
    info "  TeeAuthenticator     $CHAIN_TEEAUTHENTICATOR_ADDRESS"
    info "  TokenAllowlist       $CHAIN_TOKEN_ALLOWLIST_ADDRESS"
    info "  Allowlisted token    $TOKEN_ADDRESS"
    info "  Manager              $MANAGER_ADDRESS"
    info "  Subgraph             $SUBGRAPH_URL"
    info ""

    printf '%s%sDo not forget%s\n' "$C_BOLD" "$C_YELLOW" "$C_RESET"
    if [ "${MANAGER_FUNDED:-no}" != "yes" ]; then
        warn "the manager address $MANAGER_ADDRESS was not funded:"
        info "     without ETH on $NET_LABEL the manager cannot publish results on-chain."
    fi
    warn "the 'facilitator' service was NOT installed."
    info "     It is not part of this repository's docker-compose. If you need the"
    info "     gasless payment path (x402), add the service manually with the"
    info "     horizen/vela-facilitator:v${VELA_VERSION} image and the variables"
    info "     FACILITATOR_PRIVATE_KEY, FACILITATOR_VELA_VERSION,"
    info "     FACILITATOR_EXPLORER_BASEURL, VELA_NOVA_APPLICATION_ID."
    warn "$SECRETS_DIR and $BUNDLE_DIR/.env contain private keys:"
    info "     store them in a password manager and delete the working directory"
    info "     when you are done."
    info ""
    ok "done."
}

#-----------------------------------------------------------------------------
# Main
#-----------------------------------------------------------------------------
# Fires on any non-zero exit, including die(). An ERR trap is not usable here:
# it would also fire inside the deliberate 'set +e' blocks.
on_exit() {
    local rc=$?
    if [ "$rc" -ne 0 ]; then
        if [ -n "${WORKDIR:-}" ] && [ -f "${STATE_FILE:-}" ]; then
            printf '\n' >&2
            warn "interrupted. Resume where you left off with:"
            printf '    %s%s --workdir %s%s\n' "$C_BOLD" "$SCRIPT_NAME" "$WORKDIR" "$C_RESET" >&2
        fi
    fi
    return "$rc"
}

usage() {
    cat <<EOF
$SCRIPT_NAME - guided installer for Horizen Vela (without-TEE deployments)

Usage:
  $SCRIPT_NAME                    start a new installation
  $SCRIPT_NAME --workdir DIR      resume an installation started in DIR
  $SCRIPT_NAME --local PATH       build from a local repo checkout

Options:
  --workdir DIR   working directory of an installation already started
  --local PATH    path to a checkout of the vela repository, used instead of
                  downloading the release tag from GitHub
  -h, --help      show this message

Prerequisites: bash, curl, tar, node >= 20, npm. Docker is optional (only used
to validate the generated compose locally).
EOF
}

main() {
    WORKDIR=""
    LOCAL_REPO=""

    while [ $# -gt 0 ]; do
        case "$1" in
            --workdir) [ $# -ge 2 ] || die "--workdir requires a path"; WORKDIR="$2"; shift 2 ;;
            --local)   [ $# -ge 2 ] || die "--local requires a path"; LOCAL_REPO="$2"; shift 2 ;;
            -h|--help) usage; exit 0 ;;
            *) usage >&2; die "unknown option: $1" ;;
        esac
    done

    if [ -n "$LOCAL_REPO" ]; then
        [ -d "$LOCAL_REPO" ] || die "$LOCAL_REPO does not exist"
        LOCAL_REPO="$(cd "$LOCAL_REPO" && pwd)"
    fi

    [ -r /dev/tty ] || die "this script is interactive and requires a terminal"
    trap on_exit EXIT

    printf '%s%sHorizen Vela - guided installer%s\n' "$C_BOLD" "$C_BLUE" "$C_RESET"
    check_prereqs

    if [ -n "$WORKDIR" ]; then
        [ -d "$WORKDIR" ] || die "$WORKDIR does not exist"
        WORKDIR="$(cd "$WORKDIR" && pwd)"
        set_paths
        [ -f "$STATE_FILE" ] || die "$WORKDIR does not contain a started installation (state.env is missing)"
        state_load
        apply_network_profile "$NETWORK"
        SUBNET_PREFIX="${SUBNET_PREFIX:-$NET_SUBNET_PREFIX}"
        MANAGER_ADMIN_PORT="${MANAGER_ADMIN_PORT:-$NET_MANAGER_ADMIN_PORT}"
        AUTHORITY_HOST_PORT="${AUTHORITY_HOST_PORT:-$NET_AUTHORITY_PORT}"
        ok "resuming the installation in $WORKDIR"
        info "  version $VELA_VERSION, $NET_LABEL, completed steps:${DONE_STEPS:- none}"
    else
        run_wizard
    fi

    step_fetch_sources
    step_npm_install
    step_generate_keys
    load_keys
    step_fund_manager
    step_deploy_contracts
    step_allowlist_token
    step_deploy_subgraph
    step_build_bundle
    final_summary
}

main "$@"
