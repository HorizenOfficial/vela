# admincli

Interactive CLI for sending admin commands to the Vela Manager's admin server.
All commands go through the Manager (TCP port 4002 by default); commands
targeting the Executor are forwarded automatically via the Manager-Executor
communication channel.

## Build

```bash
# Via Makefile (outputs to ../bin/)
cd cmd && make admincli

# Or directly
go build ./cmd/admincli
```

## Usage

```bash
# Default: localhost:4002
./admincli

# Custom host/port
./admincli --host 10.0.0.5 --port 5000
```

The tool presents a numbered menu and prompts for any required fields:

```
Vela Admin CLI
Server: localhost:4002

Commands:
  1) Get Version
  2) Get Log Level
  3) Set Log Level
  4) Key Attestation
  5) Quit
Select command>
```

For choices (log level, target) you can enter either the number or the literal
value (e.g. `3` or `info`).

## Commands

### 1) Get Version

Returns the build version string for the selected target (e.g. `"0.1.0"`).
When `target=all`, the response is an aggregated object with both components:

```json
{"manager": "0.1.0", "executor": "0.1.0"}
```

When targeting a single component, the response is a plain string (`"0.1.0"`).

### 2) Get Log Level

Returns the current zerolog level for the selected target. When `target=all`,
the response is an aggregated object with both components:

```json
{"manager": "info", "executor": "debug"}
```

When targeting a single component, the response is a plain string (`"info"`).

### 3) Set Log Level

Changes the zerolog level at runtime for the selected target without restarting
the process. Accepts: `trace`, `debug`, `info`, `warn`, `error`, `fatal`,
`panic`, `disabled`. Each level enables itself and all higher-priority levels
(e.g. setting `warn` also enables `error`, `fatal`, and `panic`).

When `target=all`, both Manager and Executor are updated and the response
contains the new level from each:

```json
{
  "manager":  "warn",
  "executor": "warn"
}
```

### 4) Key Attestation

Requests an AWS Nitro Enclave attestation document from the Executor. The
attestation document cryptographically binds the Executor's keyset (P521
communication key + secp256k1 signing address) to the enclave's identity,
proving to third parties that these keys were generated inside a genuine Nitro
Enclave running specific code. The document is a CBOR-encoded COSE_Sign1
structure returned as base64.

This command only works when the Executor is running inside a Nitro Enclave
(requires `/dev/nsm`). Outside of an enclave it returns an error.

### Targets

| Value      | Behaviour                                              |
|------------|--------------------------------------------------------|
| `manager`  | Command is handled locally by the Manager              |
| `executor` | Command is forwarded to the Executor via the Manager   |
| `all`      | Applied to the Manager first, then forwarded to the Executor; response is aggregated |

## Protocol

The admin server uses a simple TCP protocol:

1. Client connects to the Manager admin server (default port 4002)
2. Client sends a JSON message terminated by a newline (`\n`)
3. Server processes the command and sends a JSON response terminated by a newline
4. Connection closes

Message format:
```json
{"type": "<string>", "target": "<string>", "data": <object or null>}
```

Response type `"response"` indicates success; type `"error"` indicates an error
with `{"code": "...", "message": "..."}` in the data field.

## Architecture: zero internal dependencies

The binary intentionally duplicates message type constants and request structs
from `pkg/admin` rather than importing them. This keeps the CLI free of internal
module dependencies so it can be built and distributed independently of the main
project (e.g. copied to a remote host without the full source tree).

## Contract tests

Because the CLI duplicates values from `pkg/admin`, there is a risk of **drift**:
if someone changes a message type number, renames a JSON field, or adds a log
level in `pkg/admin`, the CLI silently breaks — it still compiles but sends
wrong data on the wire.

`main_test.go` guards against this. It imports `pkg/admin` (test files can
import internal packages without affecting the production binary) and
cross-checks every duplicated value:

- **Message type constants** — each local constant (e.g. `setLogLevelRequest = "set_log_level"`)
  is asserted equal to the canonical `admin.SetLogLevelRequestMessage`.
- **Valid targets** — the local list is compared against `admin.TargetManager`,
  `admin.TargetExecutor`, `admin.TargetAll`.
- **Valid log levels** — the local list is compared against
  `admin.SupportedLogLevels` (parsed from the canonical comma-separated string).
- **JSON field names** — both local and upstream structs (`setLogLevelReq` vs
  `admin.SetLogLevelRequest`, etc.) are marshaled and the resulting JSON keys
  are compared to catch renamed or missing fields.

These tests run in milliseconds with no network or servers required:

```bash
go test -v ./cmd/admincli/
```

If a new command type is added to `pkg/admin`, the existing tests won't fail
(they only check the values the CLI already knows about). Adding the new command
to the CLI menu and extending the contract tests is a separate, conscious step.
