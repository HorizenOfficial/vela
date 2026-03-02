# Admin Package

The `admin` package provides administrative command interfaces and servers for managing the Horizen PES system at runtime.

## Overview

This package enables runtime administration of the Manager component through a command server interface. Clients can connect to the Manager's admin server to execute administrative commands such as checking version information, changing logging levels, or requesting key attestations.

## Architecture

The admin package defines:
- **AdminCommandServer**: Server interface for accepting admin connections
- **AdminCmdHandler**: Generic handler interface for Manager admin commands
- **AdminServer**: Concrete implementation of the admin server

### Communication Protocol

Admin commands use a JSON-based protocol over TCP or V-Socket connections:

1. Client connects to the admin server
2. Client sends an `AdminMessage` with:
   - `Type`: The command type (e.g., `"set_log_level"`)
   - `Target`: Routing target (`"manager"`, `"executor"`, `"all"`, or omitted for default `"all"`)
   - `Data`: Command-specific payload
3. Server processes the command and responds with:
   - `"response"`: Success response with data
   - `"error"`: Error response with code and message

Messages are newline-delimited JSON (`\n`).

## Manager Commands

### GetVersion

Retrieves the version of the Manager and/or Executor. The version is the git tag injected at build time via `-ldflags`. Falls back to `"dev"` when built without `-ldflags`.

**Request:**
```json
{
  "type": "get_version",
  "data": null
}
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
    "executor": "v1.2.3"
  }
}
```

The `target` field is set on the `AdminMessage` (not in `data`). Valid values: `"manager"`, `"executor"`, `"all"`, or omitted (defaults to `"all"`).

### GetLogLevel

Retrieves the current logging level of the Manager. Only supported when the Manager uses ZeroNetworkLogger (the production logger).

**Request:**
```json
{"type": "get_log_level", "data": null}
```

Or with an explicit target:
```json
{"type": "get_log_level", "target": "manager", "data": null}
```

**Response:**
```json
{
  "type": "response",
  "data": "debug"
}
```

The response contains a string with the current log level. Possible values are: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`, `disabled`.

### SetLogLevel

Changes the logging level at runtime without requiring a restart. Only supported when the component uses ZeroNetworkLogger (the production logger). This is especially useful for the Executor running inside an AWS Nitro Enclave, where restarts are not acceptable.

**Request:**
```json
{"type": "set_log_level", "data": {"level": "debug"}}
```

Or with an explicit target:
```json
{"type": "set_log_level", "target": "manager", "data": {"level": "debug"}}
```

**Response:**
```json
{
  "type": "response",
  "data": {
    "success": true,
    "level": "debug"
  }
}
```

#### Target Field

All admin commands accept an optional `target` field on the `AdminMessage` envelope (not inside `data`). The Manager admin server (the single admin entry point) validates the target and routes the command accordingly:

| `target` value | Result |
|---|---|
| `"manager"` | Applied to Manager only |
| `"executor"` | Forwarded to Executor only (Manager is not affected) |
| `"all"` | Applied to Manager locally **and** forwarded to Executor; returns an aggregated response |
| `""` / omitted | Defaults to `"all"` (a warning is logged); see `"all"` above |
| anything else | **Rejected** (unknown target) |

> **Note:** The Executor does not have its own admin server. Admin commands reach the Executor exclusively through the Manager's proxy forwarding via the communication channel (`ForwardAdminCommand`). When `target` is omitted, the command defaults to `"all"` so that both components are affected.

#### Partial Failure Semantics (target="all")

When `target` is `"all"`, the Manager always attempts both the local operation and the Executor forward, regardless of whether one fails. The response includes error fields for any component that failed:

- **One fails:** The response is still a success (`"type": "response"`). The successful component's result is returned normally; the failed component's error is in `managerError` or `executorError`.
- **Both fail:** The entire request returns an error (`"type": "error"`).

Example — SetLogLevel succeeds on Manager but Executor is unreachable:

```json
{
  "type": "response",
  "data": {
    "manager": {"success": true, "level": "debug"},
    "executorError": "connection refused"
  }
}
```

To retry the failed component, send a follow-up command with `target: "executor"`.

#### Supported Log Levels

The following log levels are supported (based on zerolog):

- `trace` - Most verbose, includes all debug information
- `debug` - Detailed debugging information
- `info` - General informational messages (default)
- `warn` - Warning messages for potentially harmful situations
- `error` - Error messages for failures that don't stop execution
- `fatal` - Critical errors that cause the application to exit
- `panic` - Severe errors that cause a panic
- `disabled` - Suppresses all log output; useful for silencing a component temporarily

#### Level Hierarchy

Log levels follow a hierarchy where setting a level enables all messages at that level and above:
- `trace` → enables all levels
- `debug` → enables debug, info, warn, error, fatal, panic
- `info` → enables info, warn, error, fatal, panic
- `warn` → enables warn, error, fatal, panic
- `error` → enables error, fatal, panic
- `fatal` → enables fatal, panic
- `panic` → enables panic only
- `disabled` → suppresses all log output

#### Error Handling

All command errors are returned with the `COMMAND_ERROR` code.

Unknown target:
```json
{
  "type": "error",
  "data": {
    "code": "COMMAND_ERROR",
    "message": "unknown target 'foo'; valid targets: 'manager', 'executor', 'all'"
  }
}
```

Invalid log level:
```json
{
  "type": "error",
  "data": {
    "code": "COMMAND_ERROR",
    "message": "invalid log level 'invalid'; supported levels: trace, debug, info, warn, error, fatal, panic, disabled"
  }
}
```

Component not using ZeroNetworkLogger:
```json
{
  "type": "error",
  "data": {
    "code": "COMMAND_ERROR",
    "message": "SetLogLevel is only supported with the ZeroNetworkLogger"
  }
}
```

#### Implementation Details

The `SetLogLevel` command:
1. Validates the `target` field matches the component (or is empty)
2. Validates the logger is a ZeroNetworkLogger (returns error otherwise)
3. Validates the level string is not empty
4. Calls `logger.SetLevel()` to update the logger configuration
5. Changes take effect immediately for all subsequent log statements

## Executor Commands (Forwarded via Manager)

All Executor commands are sent to the Manager admin server and forwarded to the Executor through the communication channel. The Executor does not have its own admin server.

### CreateKeyAttestation

Generates a key attestation document for the Executor's cryptographic keys. Always forwarded to the Executor.

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

### GetLogLevel (Executor)

Retrieves the current logging level of the Executor. Use `target: "executor"` to get only the Executor's level, or `target: "all"` to get both.

### SetLogLevel (Executor)

Changes the Executor's logging level at runtime. Use `target: "executor"` to change only the Executor's level, or `target: "all"` to change both.

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
    "manager": {"success": true, "level": "debug"},
    "executor": {"success": true, "level": "debug"}
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
                                              ├── get_version (target="all") ──> local + ForwardAdminCommand ──> aggregated response
                                              ├── set_log_level (target="all") ──> local + ForwardAdminCommand ──> aggregated response
                                              └── get_log_level (target="all") ──> local + ForwardAdminCommand ──> aggregated response
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
