# Admin Package

The `admin` package provides administrative command interfaces and servers for managing the Horizen PES system at runtime.

## Overview

This package enables runtime administration of the Manager and Executor components through a command server interface. Clients can connect to the admin server to execute administrative commands such as checking version information or changing logging levels.

## Architecture

The admin package defines:
- **AdminCommandServer**: Server interface for accepting admin connections
- **AdminCmdHandler**: Generic handler interface for both Executor and Manager admin commands
- **AdminServer**: Concrete implementation of the admin server (shared by both Executor and Manager)

### Communication Protocol

Admin commands use a JSON-based protocol over TCP or V-Socket connections:

1. Client connects to the admin server
2. Client sends an `AdminMessage` with:
   - `Type`: The command type (e.g., `SetLogLevelRequestMessage`)
   - `Data`: Command-specific payload
3. Server processes the command and responds with:
   - `AdminResponseMessage`: Success response with data
   - `AdminErrorMessage`: Error response with code and message

Messages are newline-delimited JSON (`\n`).

## Manager Commands

### GetVersion

Retrieves the current version of the Manager.

**Request:**
```json
{
  "type": 3,
  "data": null
}
```

**Response:**
```json
{
  "type": 0,
  "data": "0.1.0"
}
```

### GetLogLevel

Retrieves the current logging level of the Manager. Only supported when the Manager uses ZeroNetworkLogger (the production logger).

**Request:**
```json
{"type": 5, "data": null}
```

Or with an explicit target:
```json
{"type": 5, "data": {"target": "manager"}}
```

**Response:**
```json
{
  "type": 0,
  "data": "debug"
}
```

The response contains a string with the current log level. Possible values are: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`, `disabled`.

### SetLogLevel

Changes the logging level at runtime without requiring a restart. Only supported when the component uses ZeroNetworkLogger (the production logger). This is especially useful for the Executor running inside an AWS Nitro Enclave, where restarts are not acceptable.

**Request:**
```json
{"type": 4, "data": {"level": "debug"}}
```

Or with an explicit target:
```json
{"type": 4, "data": {"level": "debug", "target": "manager"}}
```

**Response:**
```json
{
  "type": 0,
  "data": {
    "success": true,
    "level": "debug"
  }
}
```

#### Target Field

Both `SetLogLevel` and `GetLogLevel` accept an optional `target` field. The Manager admin server (the single admin entry point) validates the target and routes the command accordingly:

| `target` value | Result |
|---|---|
| `"manager"` | Applied to Manager only |
| `"executor"` | Forwarded to Executor only (Manager is not affected) |
| `"all"` | Applied to Manager locally **and** forwarded to Executor; returns an aggregated response |
| `""` / omitted | Defaults to `"all"` (a warning is logged); see `"all"` above |
| anything else | **Rejected** (unknown target) |

> **Note:** The Executor does not have its own admin server. Admin commands reach the Executor exclusively through the Manager's proxy forwarding via the communication channel (`ForwardAdminCommand`). When `target` is omitted, the command defaults to `"all"` so that both components are affected.

#### Partial Failure Semantics (target="all")

When `target` is `"all"`, the Manager applies the command locally first, then forwards it to the Executor. This is **not atomic** — if the Executor fails (e.g. communication timeout), the Manager's change is still applied. The error message reports the partial success:

```json
{
  "type": 1,
  "data": {
    "code": "COMMAND_ERROR",
    "message": "manager level set to 'debug' successfully, but executor failed: <error details>"
  }
}
```

To recover from a partial failure, send a follow-up command with `target: "executor"` to retry on the Executor alone.

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
  "type": 1,
  "data": {
    "code": "COMMAND_ERROR",
    "message": "unknown target 'foo'; valid targets: 'manager', 'executor', 'all'"
  }
}
```

Invalid log level:
```json
{
  "type": 1,
  "data": {
    "code": "COMMAND_ERROR",
    "message": "invalid log level 'invalid'; supported levels: trace, debug, info, warn, error, fatal, panic, disabled"
  }
}
```

Component not using ZeroNetworkLogger:
```json
{
  "type": 1,
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
  "type": 2,
  "data": null
}
```

**Response:**
```json
{
  "type": 0,
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
| 0 | `AdminResponseMessage` | Success response |
| 1 | `AdminErrorMessage` | Error response |
| 2 | `KeyAttestationRequestMessage` | Request key attestation (Executor) |
| 3 | `GetVersionRequestMessage` | Request version info (Manager) |
| 4 | `SetLogLevelRequestMessage` | Change log level (Manager and Executor) |
| 5 | `GetLogLevelRequestMessage` | Get current log level (Manager and Executor) |


> **Note:** Proxy forwarding (Manager → Executor) uses the existing communication channel via `ForwardAdminCommand`, not separate admin message types. See the "Admin Proxy Forwarding" section below.

## Usage Examples

### Change the Manager log level

```bash
echo '{"type":4,"data":{"level":"debug","target":"manager"}}' | nc <manager-host> <admin-port>
```

### Get the Manager log level

```bash
echo '{"type":5,"data":{"target":"manager"}}' | nc <manager-host> <admin-port>
```

### Get the Manager version

```bash
echo '{"type":3,"data":null}' | nc <manager-host> <admin-port>
```

### Change the Executor log level

```bash
echo '{"type":4,"data":{"level":"debug","target":"executor"}}' | nc <manager-host> <admin-port>
```

### Disable all logging on the Executor

```bash
echo '{"type":4,"data":{"level":"disabled","target":"executor"}}' | nc <manager-host> <admin-port>
```

### Set log level on both Manager and Executor

Use `target: "all"` to apply to both components at once. The Manager applies locally first, then forwards to the Executor:

```bash
echo '{"type":4,"data":{"level":"debug","target":"all"}}' | nc <manager-host> <admin-port>
```

**Aggregated response:**
```json
{
  "type": 0,
  "data": {
    "manager": {"success": true, "level": "debug"},
    "executor": {"success": true, "level": "debug"}
  }
}
```

### Get log level from both Manager and Executor

```bash
echo '{"type":5,"data":{"target":"all"}}' | nc <manager-host> <admin-port>
```

**Aggregated response:**
```json
{
  "type": 0,
  "data": {
    "manager": "info",
    "executor": "debug"
  }
}
```

## Admin Proxy Forwarding (Manager → Executor)

In production (AWS Nitro Enclave), the Executor is not directly accessible — only the Manager is. The Manager can forward admin commands to the Executor through the existing bidirectional communication channel (TCP/V-Socket) using `ForwardAdminCommand`.

No additional configuration is needed — the proxy uses the same communication channel that carries `ProcessRequest` and `DeployApp` messages.

### Communication Channel Protocol

Admin commands are forwarded using `AdminCommandRequestMessage` (type 10) and `AdminCommandResponseMessage` (type 11) in the communication layer. Each request contains:
- `commandType`: A string identifier (e.g., `"set_log_level"`, `"get_log_level"`)
- `data`: Command-specific payload as `json.RawMessage`

### Proxy Flow

```
Admin Client ──> Manager Admin Server ──> Manager ExecuteCommand
                                              │
                                              ├── type 4 (target="all") ──> local + ForwardAdminCommand ──> aggregated response
                                              └── type 5 (target="all") ──> local + ForwardAdminCommand ──> aggregated response
                                                                                     │
                                                                        communication.Client.ForwardAdminCommand
                                                                                     │
                                                                                     ▼
                                                                        communication.Server (Executor)
                                                                                     │
                                                                        executor.HandleAdminCommand
```

## Server Configuration

There is one admin server, running in the Manager. The Executor does not have a standalone admin server — admin commands are forwarded to it through the communication channel.

- The Manager admin server uses TCP (configured via `MANAGER_ADMIN_PORT`)
- Client timeout controls how long a client connection can remain active
- Only one client can be connected at a time

## Testing

The package includes comprehensive tests:
- [adminserver_manager_test.go](adminserver_manager_test.go) - Tests for Manager admin server
- [adminserver_test.go](adminserver_test.go) - Tests for admin server core (shared infrastructure)
- [admin_integration_test.go](admin_integration_test.go) - Integration tests with real AdminServer over TCP, ZeroNetworkLogger, and admin command forwarding through the communication channel

To run the tests:
```bash
go test ./pkg/admin/...
```
