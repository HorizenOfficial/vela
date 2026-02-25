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
   - `Data`: Command-specific payload
3. Server processes the command and responds with:
   - `"response"`: Success response with data
   - `"error"`: Error response with code and message

Messages are newline-delimited JSON (`\n`).

## Manager Commands

### GetVersion

Retrieves the current version of the Manager.

**Request:**
```json
{
  "type": "get_version",
  "data": null
}
```

**Response:**
```json
{
  "type": "response",
  "data": "0.1.0"
}

### GetLogLevel

Retrieves the current logging level of the Manager. This is useful for verifying the active log level or when building administrative tools.

**Request:**
```json
{
  "type": "get_log_level",
  "data": null
}
```

**Response:**
```json
{
  "type": "response",
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
  "type": "set_log_level",
  "data": {
    "level": "debug"
  }
}
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
  "type": "error",
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

### CreateKeyAttestation

Requests a key attestation document for the Executor's cryptographic keys. The Manager receives this command and forwards it to the Executor over the manager-executor communication channel. The Executor produces the attestation and sends it back to the Manager, which returns it to the admin client.

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

## Message Types

| Type | Name | Description |
|------|------|-------------|
| `"response"` | `AdminResponseMessage` | Success response |
| `"error"` | `AdminErrorMessage` | Error response |
| `"key_attestation"` | `KeyAttestationRequestMessage` | Request key attestation (forwarded to Executor) |
| `"get_version"` | `GetVersionRequestMessage` | Request version info |
| `"set_log_level"` | `SetLogLevelRequestMessage` | Change log level |
| `"get_log_level"` | `GetLogLevelRequestMessage` | Get current log level |

## Usage Examples

### Change the log level

To connect to a running Manager and change the log level:

```bash
# Using netcat to send a SetLogLevel command
echo '{"type":"set_log_level","data":{"level":"debug"}}' | nc <manager-host> <admin-port>
```

### Get the current log level

To retrieve the current logging level:

```bash
# Using netcat to send a GetLogLevel command
echo '{"type":"get_log_level","data":null}' | nc <manager-host> <admin-port>
```

### Get the Manager version

To retrieve the Manager version:

```bash
# Using netcat to send a GetVersion command
echo '{"type":"get_version","data":null}' | nc <manager-host> <admin-port>
```

### Request a key attestation

To request a key attestation from the Executor (via the Manager):

```bash
# Using netcat to send a KeyAttestation command
echo '{"type":"key_attestation","data":null}' | nc <manager-host> <admin-port>
```

## Server Configuration

The admin server runs on the Manager and uses TCP.
- Client timeout controls how long a client connection can remain active
- Only one client can be connected at a time
- Commands that target the Executor (e.g., KeyAttestation) are forwarded over the manager-executor communication channel

## Testing

The package includes comprehensive tests:
- [adminserver_manager_test.go](adminserver_manager_test.go) - Tests for Manager admin commands
- [adminserver_test.go](adminserver_test.go) - Tests for admin server infrastructure

To run the tests:
```bash
go test ./pkg/admin/...
```
