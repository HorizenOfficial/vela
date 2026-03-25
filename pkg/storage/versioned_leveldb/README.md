# Versioned LevelDB Storage

The versioned LevelDB storage system provides a key-value store with the ability to version data and roll back to previous versions. This is achieved through a layered architecture that uses LevelDB as the underlying storage engine.

## Architecture

The system is composed of several layers:

1.  **`LevelDBDataLayer`**: This is the highest-level abstraction, implementing the `storage.DataLayer` interface. It acts as a unified entry point for accessing different types of data. It combines both versioned and non-versioned storage by composing multiple specialized stores. Data is segregated into different subdirectories to avoid key collisions.

2.  **`VersionedLevelDBAppStateStore`**: This component implements the `storage.ApplicationStateStore` interface and is responsible for handling data that requires versioning, such as application states and WASM bytecode. It uses the `VersionedLevelDbStorageAdapter` to interact with the underlying versioned database. It also handles the serialization and deserialization of application-specific data structures.

3.  **`LevelDBUserKeyStore`**: Non-versioned store for user keys using a simple `LevelDbStorageAdapter` to interact directly with LevelDB.

4.  **`VersionedLevelDbStorageAdapter`**: This layer adapts the `VersionedLDBKVStore` to the `storage.VersionedStorage` interface. It provides methods for getting, setting, and updating data, as well as for managing versions.

5.  **`VersionedLDBKVStore`**: This is the core of the versioning system. It directly interacts with the LevelDB database and implements the versioning logic.

### Versioning Mechanism

The versioning mechanism is based on the concept of **change sets**. Each update to the database is treated as a new version, identified by a unique `versionID`.

-   **`VersionsKey`**: Each application has its own versions key (`sha256("versions_<appID>")`), which stores a list of all active `versionID`s for that application. This list is ordered from newest to oldest. Version chains are independent per application — storing, pruning, and rolling back one app does not affect any other app.

-   **`ChangeSet`**: For each `versionID`, a corresponding `ChangeSet` is stored. The `ChangeSet` is a JSON-serialized object that records the changes made in that version. It contains three lists:
    -   `InsertedKeys`: A list of keys that were added in this version.
    -   `Removed`: A list of key-value pairs that were removed in this version. The old values are stored to allow for rollback.
    -   `Altered`: A list of key-value pairs that were modified in this version. The old values are stored to allow for rollback.

### Update Operation

When an `Update` operation is performed:

1.  A new `versionID` is provided.
2.  The system checks if the `versionID` already exists. If it does, an error is returned.
3.  A `ChangeSet` is created by comparing the keys to be updated and removed with their current values in the database.
4.  The new `versionID` is prepended to the application's version list (stored under its per-app versions key).
5.  If the number of versions for that application exceeds the `versionsToKeep` limit, the oldest versions are pruned. This limit is applied independently per application.
6.  The `ChangeSet` is serialized and stored with the `versionID` as the key.
7.  The actual data is updated in the database.
8.  All these operations are performed within a single LevelDB transaction to ensure atomicity.

### Rollback Operation

When a `RollbackTo` operation is performed with a target `versionID`:

1.  The system finds the target `versionID` in the list of versions.
2.  It then iterates through all the versions that are *newer* than the target version.
3.  For each of these newer versions, it retrieves the corresponding `ChangeSet`.
4.  It applies the *inverse* of the operations in the `ChangeSet`:
    -   Keys in `InsertedKeys` are deleted.
    -   Key-value pairs in `Removed` are re-inserted.
    -   Key-value pairs in `Altered` are restored to their old values.
5.  The `versionID`s of the rolled-back versions are removed from the list of versions.
6.  All these operations are performed within a single LevelDB transaction.

This architecture ensures that the database can be reliably rolled back to any previous state while maintaining data integrity.
