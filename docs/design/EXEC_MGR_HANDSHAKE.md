# Handshake Protocol between Manager and Executor

The handshake process is initiated by the Executor every time a new connection is established with the Manager. The primary goal of the handshake is to ensure the Executor has a valid and consistent set of cryptographic keys (`EnclaveKeySet`) to perform its duties.

The Manager's first reply (`GetKeysetRecoveryResponse`) also carries `requestTimeoutMs`, the Manager's own request timeout (`MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC`, in milliseconds). It is the value every later per-request execution budget is derived from, and the Executor validates its guest execution bound (`EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS`) against it **before** restoring or generating any keyset: the bound must be strictly below the timeout minus the 2 s reply margin, otherwise the Manager would give up on a runaway guest before the Executor could sign its failure, and the request would be retried forever (see the sizing notes in `WASM_HOST_ABI.md`). The enclave cannot read the Manager's configuration, which is why the value travels in the handshake. The field is optional: an older Manager omits it, and the Executor logs a warning and proceeds unchecked.

There are two main scenarios for the handshake:

1.  **First Connection**: This occurs when the Manager connects to the Executor for the very first time, and no keyset recovery data exists in the Manager's persistent storage.
2.  **Reconnection**: This occurs on subsequent connections, where the Manager already holds recovery data from a previous session.

---

### Scenario 1: First Connection (No Keyset Data)

In this scenario, the Executor generates a new keyset and sends the corresponding recovery data to the Manager for storage.

```
Manager                                     Executor
   |                                            |
   |----------- TCP/VSock Connection ---------->|
   |                                            |
   |                                            | (Connection established)
   |<--- GetKeysetRecoveryRequest (ID: 1) <-----| (Executor initiates handshake)
   |                                            |
   | (Queries data layer, finds nothing)        |
   |                                            |
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (DataFound: false, RequestTimeoutMs)
   |                                            |
   |                                            | (Checks guest bound fits RequestTimeoutMs)
   |                                            | (Generates new keyset)
   |<--- SetKeysetRecoveryRequest (ID: 2) <-----| (Contains new recovery data)
   |                                            |
   | (Stores recovery data in data layer)       |
   |                                            |
   |---- SetKeysetRecoveryResponse (ID: 2) ---->|
   |                                            |
   o (Handshake complete)                       o (Handshake complete)
   |                                            |
```

---

### Scenario 2: Reconnection (Successful Keyset Recovery)

Here, the Manager provides the existing recovery data to the Executor, which then restores its keyset.

```
Manager                                     Executor
   |                                            |
   |----------- TCP/VSock Connection ---------->|
   |                                            |
   |                                            | (Connection established)
   |<--- GetKeysetRecoveryRequest (ID: 1) <-----| (Executor initiates handshake)
   |                                            |
   | (Queries data layer, finds recovery data)  |
   |                                            |
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (DataFound: true, with recovery data, RequestTimeoutMs)
   |                                            |
   |                                            | (Checks guest bound fits RequestTimeoutMs)
   |                                            | (Restores keyset from data)
   |<--- KeysetRecoveryResult (ID: 2) <---------| (Result: success)
   |                                            |
   o (Handshake complete)                       o (Handshake complete)
   |                                            |
```

---

### Scenario 4: Guest Execution Bound Does Not Fit the Manager's Timeout

If `EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS` is not strictly below `RequestTimeoutMs` minus the reply margin, the Executor rejects the connection before touching any keyset. The Manager blocks on the handshake at start-up and exits with the Executor's error, so the misconfiguration is visible on both sides and the Executor never serves a request under it.

```
Manager                                     Executor
   |                                            |
   |----------- TCP/VSock Connection ---------->|
   |                                            |
   |<--- GetKeysetRecoveryRequest (ID: 1) <-----|
   |                                            |
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (…, RequestTimeoutMs: 5000)
   |                                            |
   |                                            | (Guest bound 10000 ms >= 5000 - 2000: reject)
   |<--- KeysetRecoveryResult (ID: 2) <---------| (Result: failure naming both settings)
   |                                            |
   x (Start returns the error, Manager exits)   x (Executor closes connection, no keyset touched)
   |                                            |
```

---

### Scenario 3: Reconnection (Keyset Recovery Failure)

If the Executor fails to restore its keyset from the data provided by the Manager (e.g., due to data corruption), it notifies the Manager of the failure.

```
Manager                                     Executor
   |                                            |
   |----------- TCP/VSock Connection ---------->|
   |                                            |
   |                                            | (Connection established)
   |<--- GetKeysetRecoveryRequest (ID: 1) <-----|
   |                                            |
   | (Queries data layer, finds corrupted data) |
   |                                            |
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (DataFound: true, with corrupted data)
   |                                            |
   |                                            | (Fails to restore keyset)
   |<--- KeysetRecoveryResult (ID: 2) <---------| (Result: failure with error message)
   |                                            |
   x (Handshake fails, Manager may terminate)   x (Handshake fails, Executor closes connection)
   |                                            |
```
