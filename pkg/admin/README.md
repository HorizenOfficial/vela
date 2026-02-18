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

Both `SetLogLevel` and `GetLogLevel` accept an optional `target` field. Each admin server validates the target matches itself:

| Connected to | `target` value | Result |
|---|---|---|
| Manager `:4002` | `"manager"` | Handled |
| Manager `:4002` | `""` / omitted | Handled (applies to self) |
| Manager `:4002` | `"all"` | Handled (applies to self; see note below) |
| Manager `:4002` | `"executor"` | **Rejected** |
| Manager `:4002` | anything else | **Rejected** (unknown target) |
| Executor `:4001` | `"executor"` | Handled |
| Executor `:4001` | `""` / omitted | Handled (applies to self) |
| Executor `:4001` | `"all"` | Handled (applies to self; see note below) |
| Executor `:4001` | `"manager"` | **Rejected** |
| Executor `:4001` | anything else | **Rejected** (unknown target) |

> **Note on `"all"`:** Currently `"all"` is accepted and applies to the receiving component only (same as omitted). When the Manager gains proxy capability, `"all"` on the Manager will apply the command locally **and** forward it to the Executor, returning an aggregated response.

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

Wrong target:
```json
{
  "type": 1,
  "data": {
    "code": "COMMAND_ERROR",
    "message": "this is the manager admin server, target 'executor' is not supported; connect to the executor admin server directly"
  }
}
```

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

## Executor Commands

### CreateKeyAttestation

Generates a key attestation document for the Executor's cryptographic keys.

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

Retrieves the current logging level of the Executor. Same behavior as the Manager command (type 5), but accepts target `"executor"` instead of `"manager"`.

### SetLogLevel (Executor)

Changes the Executor's logging level at runtime. Same behavior as the Manager command (type 4), but accepts target `"executor"` instead of `"manager"`.

## Message Types

| Type | Name | Description |
|------|------|-------------|
| 0 | `AdminResponseMessage` | Success response |
| 1 | `AdminErrorMessage` | Error response |
| 2 | `KeyAttestationRequestMessage` | Request key attestation (Executor) |
| 3 | `GetVersionRequestMessage` | Request version info (Manager) |
| 4 | `SetLogLevelRequestMessage` | Change log level (Manager and Executor) |
| 5 | `GetLogLevelRequestMessage` | Get current log level (Manager and Executor) |
| 6 | `ExecutorKeyAttestationRequestMessage` | Reserved: proxy key attestation via Manager |
| 7 | `ExecutorSetLogLevelRequestMessage` | Reserved: proxy SetLogLevel via Manager |
| 8 | `ExecutorGetLogLevelRequestMessage` | Reserved: proxy GetLogLevel via Manager |

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
echo '{"type":4,"data":{"level":"debug","target":"executor"}}' | nc <executor-host> <admin-port>
```

### Disable all logging on the Executor

```bash
echo '{"type":4,"data":{"level":"disabled","target":"executor"}}' | nc <executor-host> <admin-port>
```

## Server Configuration

There are two admin servers: one in the Manager and one in the Executor.
- The Manager admin server always uses TCP
- The Executor admin server uses a connection factory that determines TCP or V-Socket transport
- Client timeout controls how long a client connection can remain active
- Only one client can be connected at a time

## Testing

The package includes comprehensive tests:
- [adminserver_manager_test.go](adminserver_manager_test.go) - Tests for Manager admin server
- [adminserver_test.go](adminserver_test.go) - Tests for Executor admin server
- [admin_integration_test.go](admin_integration_test.go) - Integration test with real AdminServer over TCP and ZeroNetworkLogger

To run the tests:
```bash
go test ./pkg/admin/...
```
