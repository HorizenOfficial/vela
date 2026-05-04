package fullstack_test

import (
	"testing"

	// This import pulls in go-ethereum's simulated.Backend → DataDog/zstd
	_ "github.com/HorizenOfficial/vela/pkg/testutil/fullstack"

	// This import pulls in wasmtime-go → libwasmtime.a → bundled zstd
	_ "github.com/HorizenOfficial/vela/pkg/wasm"
)

// TestLinking verifies that the fullstack package (simulated backend) and the
// wasm package (wasmtime-go) can coexist in the same test binary without
// duplicate zstd symbol errors at link time.
//
// If this test compiles and runs, the cgo LDFLAGS fix works.
func TestLinking(t *testing.T) {
	t.Log("fullstack + wasmtime linked successfully")
}
