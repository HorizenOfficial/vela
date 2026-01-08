# Key Attestation CMD

The `key_attestation_cmd` is a command-line tool designed to request and retrieve a key attestation from the executor. It can communicate with the server over a TCP socket or a VSOCK for virtualized environments.

## How it Works

The tool establishes a connection with the server and sends an `AdminMessage` of type `KeyAttestationRequestMessage`. It then waits for a response.

- If the operation is successful, it receives an `AdminResponseMessage` containing the attestation data, which is then printed to the console in JSON format.
- If an error occurs, it receives an `AdminErrorMessage`, and the error details are printed to `stderr`.

A timeout is used to prevent the client from waiting indefinitely for a response.

## Usage

You can run the tool from the command line, specifying the server connection details.

```bash
go run cmd/key_attestation_cmd/main.go [flags]
```

### Flags

| Flag          | Type   | Description                                                     | Default       |
|---------------|--------|-----------------------------------------------------------------|---------------|
| `-servertype` | string | The type of server connection. Can be `tcp` or `vsocket`.       | `vsocket`     |
| `-addr`       | string | The server address for TCP connections (e.g., `host:port`).     | `localhost:12345` |
| `-cid`        | uint   | The context ID (CID) for VSOCK connections.                     | `3`           |
| `-port`       | uint   | The port number for VSOCK connections.                          | `12345`       |
| `-timeout`    | uint   | The timeout in seconds to wait for a response from the server.  | `60`          |

## Examples

### Using VSOCK (Default)

To connect to a server listening on VSOCK with CID `3` and port `12345`:
```bash
go run cmd/key_attestation_cmd/main.go
```

To connect to a different VSOCK endpoint:
```bash
go run cmd/key_attestation_cmd/main.go -cid 4 -port 54321
```

### Using TCP

To connect to a server listening on a TCP socket at `localhost:9999`:
```bash
go run cmd/key_attestation_cmd/main.go -servertype tcp -addr localhost:9999
```
