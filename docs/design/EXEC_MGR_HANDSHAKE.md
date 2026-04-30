# Handshake Protocol between Manager and Executor

The handshake process is initiated by the Executor every time a new connection is established with the Manager. The primary goal of the handshake is to ensure the Executor has a valid and consistent set of cryptographic keys (`EnclaveKeySet`) to perform its duties.

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
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (DataFound: false)
   |                                            |
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
   |---- GetKeysetRecoveryResponse (ID: 1) ---->| (DataFound: true, with recovery data)
   |                                            |
   |                                            | (Restores keyset from data)
   |<--- KeysetRecoveryResult (ID: 2) <---------| (Result: success)
   |                                            |
   o (Handshake complete)                       o (Handshake complete)
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
