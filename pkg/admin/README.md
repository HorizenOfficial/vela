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
  "type": 1,
  "data": "0.1.0"
}

### GetLogLevel

Retrieves the current logging level of the Manager. This is useful for verifying the active log level or when building administrative tools.

**Request:**
```json
{
  "type": 5,
  "data": null
}
```

**Response:**
```json
{
  "type": 1,
  "data": "debug"
}
```

The response contains a string with the current log level. Possible values are: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`.

**Location in code:** [manager.go:295-305](../manager/manager.go#L295-L305)

### SetLogLevel

Changes the logging level at runtime without requiring a restart. This is useful for debugging or reducing verbosity in production.

**Request:**
```json
{
  "type": 4,
  "data": {
    "level": "debug"
  }
}
```

**Response:**
```json
{
  "type": 1,
  "data": {
    "success": true,
    "level": "debug"
  }
}
```

#### Supported Log Levels

The following log levels are supported (based on zerolog):

- `trace` - Most verbose, includes all debug information
- `debug` - Detailed debugging information
- `info` - General informational messages (default)
- `warn` - Warning messages for potentially harmful situations
- `error` - Error messages for failures that don't stop execution
- `fatal` - Critical errors that cause the application to exit
- `panic` - Severe errors that cause a panic

#### Level Hierarchy

Log levels follow a hierarchy where setting a level enables all messages at that level and above:
- `trace` → enables all levels
- `debug` → enables debug, info, warn, error, fatal, panic
- `info` → enables info, warn, error, fatal, panic
- `warn` → enables warn, error, fatal, panic
- `error` → enables error, fatal, panic
- `fatal` → enables fatal, panic
- `panic` → enables panic only

#### Error Handling

If an invalid log level is provided, the command returns an error:

```json
{
  "type": 2,
  "data": {
    "code": "ERROR_SETTING_LOG_LEVEL",
    "message": "invalid log level 'invalid': Unknown Level String: 'invalid', defaulting to NoLevel"
  }
}
```

#### Implementation Details

The `SetLogLevel` command:
1. Validates the level string is not empty
2. Calls the `ManagerCmdHandler.SetLogLevel()` method
3. The handler updates the logger configuration using `logger.SetLevel()`
4. For `ZeroLogger`, both console and file loggers are updated
5. For `ZeroNetworkLogger`, the network logger level is updated
6. Changes take effect immediately for all subsequent log statements

**Location in code:** [manager.go:282-292](../manager/manager.go#L282-L292)

## Executor Commands

### CreateKeyAttestation

Generates a key attestation document for the Executor's cryptographic keys.

**Request:**
```json
{
  "type": 0,
  "data": null
}
```

**Response:**
```json
{
  "type": 1,
  "data": "<base64-encoded-attestation>"
}
```

## Message Types

| Type | Name | Description |
|------|------|-------------|
| 0 | `KeyAttestationRequestMessage` | Request key attestation (Executor) |
| 1 | `AdminResponseMessage` | Success response |
| 2 | `AdminErrorMessage` | Error response |
| 3 | `GetVersionRequestMessage` | Request version info (Manager) |
| 4 | `SetLogLevelRequestMessage` | Change log level (Manager) |
| 5 | `GetLogLevelRequestMessage` | Get current log level (Manager) |

## Usage Examples

### Change the log level

To connect to a running Manager and change the log level:

```bash
# Using netcat to send a SetLogLevel command
echo '{"type":4,"data":{"level":"debug"}}' | nc <manager-host> <admin-port>
```

### Get the current log level

To retrieve the current logging level:

```bash
# Using netcat to send a GetLogLevel command
echo '{"type":5,"data":null}' | nc <manager-host> <admin-port>
```

### Get the Manager version

To retrieve the Manager version:

```bash
# Using netcat to send a GetVersion command
echo '{"type":3,"data":null}' | nc <manager-host> <admin-port>
```

## Server Configuration

The admin server is configured through the Manager's configuration:
- Connection factory determines TCP or V-Socket transport
- Client timeout controls how long a client connection can remain active
- Only one client can be connected at a time

## Testing

The package includes comprehensive tests:
- [adminserver_manager_test.go](adminserver_manager_test.go) - Tests for Manager admin server
- [adminserver_test.go](adminserver_test.go) - Tests for Executor admin server

To run the tests:
```bash
go test ./pkg/admin/...
```
