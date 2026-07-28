# Admin Package

The `admin` package provides administrative command interfaces and servers for managing the Horizen Vela system at runtime.

## Overview

This package enables runtime administration of the Manager and Executor components through a command server interface. Clients connect to the Manager's admin server, which can handle commands locally or forward them to the Executor via the communication channel. Supported administrative commands include checking version information, changing logging levels, and requesting key attestations.

## Architecture

The admin package defines:
- **AdminCommandServer**: Server interface for accepting admin connections
- **AdminCmdHandler**: Generic handler interface for admin commands
- **AdminServer**: Concrete implementation of the admin server

### Communication Protocol

Admin commands use a JSON-based protocol. The admin server listens on **TCP**; commands targeting the Executor are forwarded over the Manager↔Executor communication channel, which supports both **TCP** and **V-Socket** (used in AWS Nitro Enclaves).

```
                         TCP                          TCP or V-Socket
Admin Client ─────────────────> Manager Admin Server ─────────────────> Executor
              (admin protocol)                        (comm channel:
                                                       ForwardAdminCommand)
```

1. Client connects to the admin server (TCP)
2. Client sends an `AdminMessage` with:
   - `Type`: The command type (e.g., `"set_log_level"`)
   - `Target`: Routing target (`"manager"`, `"executor"`, `"all"`, or omitted for default `"all"`)
   - `Data`: Command-specific payload
3. Server processes the command and responds with:
   - `"response"`: Success response with data
   - `"error"`: Error response with code and message

Messages are newline-delimited JSON (`\n`).

## Commands

### Target Field

All admin commands accept an optional `target` field on the `AdminMessage` envelope (not inside `data`). The Manager admin server (the single admin entry point) validates the target and routes the command accordingly:

| `target` value | Result |
|---|---|
| `"manager"` | Applied to Manager only |
| `"executor"` | Forwarded to Executor only (Manager is not affected) |
| `"all"` | Applied to Manager locally **and** forwarded to Executor; returns an aggregated response |
| `""` / omitted | Defaults to `"all"` (a warning is logged); see `"all"` above |
| anything else | **Rejected** (unknown target) |

> **Note:** The Executor does not have its own admin server. Admin commands reach the Executor exclusively through the Manager's proxy forwarding via the communication channel (`ForwardAdminCommand`). When `target` is omitted, the command defaults to `"all"` so that both components are affected.

### Partial Failure Semantics (target="all")

When `target` is `"all"`, the Manager always attempts both the local operation and the Executor forward, regardless of whether one fails. The response is always an aggregated `"type": "response"` with per-component result and error fields:

- **One fails:** The successful component's result is returned normally; the failed component's error is in `managerError` or `executorError`.
- **Both fail:** Both `managerError` and `executorError` are populated; the result fields are empty.

Example — SetLogLevel succeeds on Manager but Executor is unreachable:

```json
{
  "type": "response",
  "data": {
    "manager": "debug",
    "executorError": "connection refused"
  }
}
```

To retry the failed component, send a follow-up command with `target: "executor"`.

### Common Commands

These commands can be sent to both the Manager and the Executor using the `target` field.

#### GetVersion

Retrieves the version of the targeted component(s). The version is the git tag injected at build time via `-ldflags`. Falls back to `"dev"` when built without `-ldflags`.

**Request:**
```json
{"type": "get_version", "data": null}
```

Or with an explicit target:
```json
{"type": "get_version", "target": "manager", "data": null}
```

**Response (single target):**
```json
{
  "type": "response",
  "data": "v1.2.3"
}
```

**Response (target="all"):**
```json
{
  "type": "response",
  "data": {
    "manager": "v1.2.3",
    "executor": "v1.2.3",
    "executorPcr0": "a1b2c3..."
  }
}
```

`executorPcr0` is the hex-encoded PCR0 of the running enclave image, as reported by the Executor during the handshake. It is omitted (empty) in non-Nitro (TCP/dev) mode where no NSM device exists.

#### GetLogLevel

Retrieves the current logging level of the targeted component(s). Only supported when the component uses ZeroNetworkLogger (the production logger).

**Request:**
```json
{"type": "get_log_level", "data": null}
```

Or with an explicit target:
```json
{"type": "get_log_level", "target": "executor", "data": null}
```

**Response (single target):**
```json
{
  "type": "response",
  "data": "debug"
}
```

**Response (target="all"):**
```json
{
  "type": "response",
  "data": {
    "manager": "info",
    "executor": "debug"
  }
}
```

The response contains a string with the current log level. Possible values are: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`, `disabled`.

#### SetLogLevel

Changes the logging level of the targeted component(s) at runtime without requiring a restart. Only supported when the component uses ZeroNetworkLogger (the production logger). This is especially useful for the Executor running inside an AWS Nitro Enclave, where restarts are not acceptable.

**Request:**
```json
{"type": "set_log_level", "data": {"level": "debug"}}
```

Or with an explicit target:
```json
{"type": "set_log_level", "target": "all", "data": {"level": "debug"}}
```

**Response (single target):**
```json
{
  "type": "response",
  "data": {
    "level": "debug"
  }
}
```

**Response (target="all"):**
```json
{
  "type": "response",
  "data": {
    "manager": "debug",
    "executor": "debug"
  }
}
```

##### Supported Log Levels

The following log levels are supported (based on zerolog):

- `trace` - Most verbose, includes all debug information
- `debug` - Detailed debugging information
- `info` - General informational messages (default)
- `warn` - Warning messages for potentially harmful situations
- `error` - Error messages for failures that don't stop execution
- `fatal` - Critical errors that cause the application to exit
- `panic` - Severe errors that cause a panic
- `disabled` - Suppresses all log output; useful for silencing a component temporarily

##### Level Hierarchy

Log levels follow a hierarchy where setting a level enables all messages at that level and above:
- `trace` → enables all levels
- `debug` → enables debug, info, warn, error, fatal, panic
- `info` → enables info, warn, error, fatal, panic
- `warn` → enables warn, error, fatal, panic
- `error` → enables error, fatal, panic
- `fatal` → enables fatal, panic
- `panic` → enables panic only
- `disabled` → suppresses all log output

##### Implementation Details

The `SetLogLevel` command:
1. Validates the logger is a ZeroNetworkLogger (returns error otherwise)
2. Validates the level string is not empty
3. Calls `logger.SetLevel()` to update the logger configuration
4. Changes take effect immediately for all subsequent log statements

### Executor Commands

These commands are always forwarded to the Executor through the communication channel. They do not apply to the Manager.

#### CreateKeyAttestation

Generates a key attestation document for the Executor's cryptographic keys. Always forwarded to the Executor. Only accepts `target: "executor"` or `target: "all"` (the default); `target: "manager"` is rejected.

**Request:**
```json
{
  "type": "key_attestation",
  "data": null
}
```

**Response:**
```json
{
  "type": "response",
  "data": "<base64-encoded-attestation>"
}
```

#### Shutdown (internal)

`shutdown` (`AdminCmdShutdown`) asks the Executor to exit cleanly. It is **not** an
operator-facing verb — the Manager admin server does not route it. It is forwarded over the
communication channel by the Manager's swap-drain sequence (drain step 3): after the
in-flight request finishes, the Manager sends `shutdown`, the Executor acks with
`{"stopping": true}` and then exits. The Manager treats the ensuing disconnect as a normal
outcome, not a fatal error.

### Error Handling

Errors are returned as `"type": "error"` messages with a `code` and `message` field:

```json
{
  "type": "error",
  "data": {
    "code": "<error_code>",
    "message": "<human-readable description>"
  }
}
```

| Code | When |
|---|---|
| `COMMAND_ERROR` | Command execution failed (unknown target, invalid log level, unsupported logger, unsupported command type, invalid request data, etc.) |
| `INVALID_REQUEST` | Malformed message or server is busy (only one client at a time) |
| `INTERNAL_ERROR` | Internal server failure (e.g., no handler set, response marshaling error) |

## Message Types

| Type | Name | Description |
|------|------|-------------|
| `"response"` | `AdminResponseMessage` | Success response |
| `"error"` | `AdminErrorMessage` | Error response |
| `"key_attestation"` | `KeyAttestationRequestMessage` | Request key attestation (forwarded to Executor) |
| `"get_version"` | `GetVersionRequestMessage` | Request version info |
| `"set_log_level"` | `SetLogLevelRequestMessage` | Change log level (Manager and Executor) |
| `"get_log_level"` | `GetLogLevelRequestMessage` | Get current log level (Manager and Executor) |

> **Note:** Proxy forwarding (Manager → Executor) uses the existing communication channel via `ForwardAdminCommand`, not separate admin message types. See the "Admin Proxy Forwarding" section below.

## Usage Examples

### Change the Manager log level

```bash
echo '{"type":"set_log_level","target":"manager","data":{"level":"debug"}}' | nc <manager-host> <admin-port>
```

### Get the Manager log level

```bash
echo '{"type":"get_log_level","target":"manager"}' | nc <manager-host> <admin-port>
```

### Get the version of both Manager and Executor

```bash
echo '{"type":"get_version"}' | nc <manager-host> <admin-port>
```

### Change the Executor log level

```bash
echo '{"type":"set_log_level","target":"executor","data":{"level":"debug"}}' | nc <manager-host> <admin-port>
```

### Disable all logging on the Executor

```bash
echo '{"type":"set_log_level","target":"executor","data":{"level":"disabled"}}' | nc <manager-host> <admin-port>
```

### Set log level on both Manager and Executor

Use `target: "all"` to apply to both components at once. The Manager applies locally first, then forwards to the Executor:

```bash
echo '{"type":"set_log_level","target":"all","data":{"level":"debug"}}' | nc <manager-host> <admin-port>
```

**Aggregated response:**
```json
{
  "type": "response",
  "data": {
    "manager": "debug",
    "executor": "debug"
  }
}
```

### Get log level from both Manager and Executor

```bash
echo '{"type":"get_log_level","target":"all"}' | nc <manager-host> <admin-port>
```

**Aggregated response:**
```json
{
  "type": "response",
  "data": {
    "manager": "info",
    "executor": "debug"
  }
}
```

### Request a key attestation

```bash
echo '{"type":"key_attestation","data":null}' | nc <manager-host> <admin-port>
```

## Admin Proxy Forwarding (Manager → Executor)

In production (AWS Nitro Enclave), the Executor is not directly accessible — only the Manager is. The Manager can forward admin commands to the Executor through the existing bidirectional communication channel (TCP/V-Socket) using `ForwardAdminCommand`.

No additional configuration is needed — the proxy uses the same communication channel that carries `ProcessRequest` and `DeployApp` messages.

### Communication Channel Protocol

Admin commands are forwarded using `AdminCommandRequestMessage` and `AdminCommandResponseMessage` in the communication layer. Each request contains:
- `commandType`: A string identifier (e.g., `"get_version"`, `"set_log_level"`, `"get_log_level"`, `"key_attestation"`)
- `data`: Command-specific payload as `json.RawMessage`

### Proxy Flow

```
Admin Client ──> Manager Admin Server ──> Manager ExecuteCommand
                                              │
                                              ├── target="manager"  ──> local only
                                              ├── target="executor" ──> ForwardAdminCommand only
                                              ├── target="all"      ──> local + ForwardAdminCommand ──> aggregated response
                                              └── key_attestation   ──> ForwardAdminCommand only (always)
                                                                           │
                                                              communication.Client.ForwardAdminCommand
                                                                           │
                                                                           ▼
                                                              communication.Server (Executor)
                                                                           │
                                                                           ▼
                                                              executor.HandleAdminCommand
```

## Server Configuration

There is one admin server, running in the Manager. The Executor does not have a standalone admin server — admin commands are forwarded to it through the communication channel.

- The Manager admin server uses TCP (configured via `MANAGER_ADMIN_PORT`)
- Client timeout controls how long a client connection can remain active
- Only one client can be connected at a time
- Commands that target the Executor (e.g., KeyAttestation) are forwarded over the manager-executor communication channel

## Testing

The package includes comprehensive tests:
- [adminserver_manager_test.go](adminserver_manager_test.go) - Tests for Manager admin server
- [adminserver_test.go](adminserver_test.go) - Tests for admin server core (shared infrastructure)
- [admin_integration_test.go](admin_integration_test.go) - Integration tests with real AdminServer over TCP, ZeroNetworkLogger, and admin command forwarding through the communication channel

To run the tests:
```bash
go test ./pkg/admin/...
```
