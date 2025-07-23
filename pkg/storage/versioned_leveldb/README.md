# Versioned LevelDB Storage

This document outlines the software architecture of the versioned LevelDB storage system.

## Overview

The versioned LevelDB storage provides a key-value store with the ability to version data and roll back to previous versions. This is achieved through a layered architecture that uses LevelDB as the underlying storage engine.

## Architecture

The system is composed of three main layers:

1.  **`VersionedLevelDBDataLayer`**: This is the highest-level abstraction, implementing the `storage.ApplicationStateStore` interface. It is responsible for handling application-specific data, such as application states, WASM bytecode, and user keys. It serializes and deserializes data and uses the `VersionedLevelDbStorageAdapter` to interact with the storage.

2.  **`VersionedLevelDbStorageAdapter`**: This layer adapts the `VersionedLDBKVStore` to the `storage.VersionedStorage` interface. It provides methods for getting, setting, and updating data, as well as for managing versions.

3.  **`VersionedLDBKVStore`**: This is the core of the versioning system. It directly interacts with the LevelDB database and implements the versioning logic.

### Versioning Mechanism

The versioning mechanism is based on the concept of **change sets**. Each update to the database is treated as a new version, identified by a unique `versionID`.

-   **`VersionsKey`**: A special key (`sha256("versions")`) is used to store a list of all active `versionID`s in the database. This list is ordered from newest to oldest.

-   **`ChangeSet`**: For each `versionID`, a corresponding `ChangeSet` is stored. The `ChangeSet` is a JSON-serialized object that records the changes made in that version. It contains three lists:
    -   `InsertedKeys`: A list of keys that were added in this version.
    -   `Removed`: A list of key-value pairs that were removed in this version. The old values are stored to allow for rollback.
    -   `Altered`: A list of key-value pairs that were modified in this version. The old values are stored to allow for rollback.

### Update Operation

When an `Update` operation is performed:

1.  A new `versionID` is provided.
2.  The system checks if the `versionID` already exists. If it does, an error is returned.
3.  A `ChangeSet` is created by comparing the keys to be updated and removed with their current values in the database.
4.  The new `versionID` is prepended to the list of versions stored under `VersionsKey`.
5.  If the number of versions exceeds the `versionsToKeep` limit, the oldest versions are pruned.
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

## TODO

-   Modify the `storage.ApplicationStateStore` interface in `pkg/storage/interface.go` to include a `Rollback` method. This will allow higher-level components to trigger a rollback of storage operations to a specific version, abstracting the underlying versioning mechanism.
