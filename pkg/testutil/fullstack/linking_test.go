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
// The assertion is the link itself: if this test binary builds and runs, the two
// zstd copies did not collide. It used to demonstrate that the
// -Wl,--allow-multiple-definition workaround worked; since that workaround was
// removed it demonstrates the stronger property that none is needed. See
// cgo_ldflags.go for why, and for what would make the collision return.
func TestLinking(t *testing.T) {
	t.Log("fullstack + wasmtime linked successfully")
}
