package versioned_leveldb

// This file is used to export unexported identifiers for testing purposes.
// In this way we can achieve selective exposure of symbols (they are private to the world but public for the tests)

var (
	// TestAppStatePrefix is the prefix for application state keys.
	TestAppStatePrefix = appStatePrefix
	// TestWasmPrefix is the prefix for wasm bytecode keys.
	TestWasmPrefix = wasmPrefix
)

// GetAdapter_ForTest returns the underlying VersionedLevelDbStorageAdapter instance for testing purposes.
func (vdl *LevelDBDataLayer) GetAdapter_ForTest() *VersionedLevelDbStorageAdapter {
	return vdl.VersionedLevelDBAppStateStore.getAdapter()
}

// TestGenerateVersionID returns the internal method for testing purposes.
func GenerateVersionID_ForTest(key []byte, data []byte) []byte {
	return generateVersionID(key, data)
}
