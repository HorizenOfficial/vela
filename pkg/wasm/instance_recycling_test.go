package wasm

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/bytecodealliance/wasmtime-go"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// watResult returns a WAT data segment holding a length-prefixed guest result at offset.
func watResult(offset int, json string) string {
	var prefix [4]byte
	binary.LittleEndian.PutUint32(prefix[:], uint32(len(json)))
	var b strings.Builder
	fmt.Fprintf(&b, "(data (i32.const %d) \"", offset)
	for _, c := range append(prefix[:], json...) {
		fmt.Fprintf(&b, "\\%02x", c)
	}
	b.WriteString("\")")
	return b.String()
}

// recycleTestWat builds a guest whose allocate never frees (bump allocator, memory capped at
// maxPages) and whose process_request counts its calls in a global: the first call on an
// instance returns state [1], later calls state [2]. spin is the busy-loop length of process_request.
func recycleTestWat(maxPages, spin int) string {
	return fmt.Sprintf(`(module
  (memory (export "memory") 1 %d)
  (global $top (mut i32) (i32.const 4096))
  (global $calls (mut i32) (i32.const 0))
  %s
  %s
  %s
  %s
  (func (export "allocate") (param $n i32) (result i32)
    (local $old i32) (local $need i32)
    (local.set $old (global.get $top))
    (local.set $need (i32.shr_u (i32.add (i32.add (local.get $old) (local.get $n)) (i32.const 65535)) (i32.const 16)))
    (if (i32.gt_u (local.get $need) (memory.size))
      (then
        (if (i32.eq (memory.grow (i32.sub (local.get $need) (memory.size))) (i32.const -1))
          (then (return (i32.const 0))))))
    (global.set $top (i32.add (local.get $old) (local.get $n)))
    (local.get $old))
  (func (export "deallocate") (param i32 i32))
  (func (export "load_module") (param i64) (result i32) (i32.const 300))
  (func (export "deploy") (param i64 i32 i32) (result i32) (i32.const 300))
  (func (export "deposit") (param i64 i32 i32 i32 i32 i32 i32 i32 i32) (result i32) (i32.const 400))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    (local $i i32)
    (loop $spin
      (local.set $i (i32.add (local.get $i) (i32.const 1)))
      (br_if $spin (i32.lt_u (local.get $i) (i32.const %d))))
    (global.set $calls (i32.add (global.get $calls) (i32.const 1)))
    (if (result i32) (i32.eq (global.get $calls) (i32.const 1))
      (then (i32.const 100))
      (else (i32.const 200))))
)`, maxPages,
		watResult(100, `{"state":[1],"events":[],"appEvents":[],"withdrawals":[],"fuel":"0x1"}`),
		watResult(200, `{"state":[2],"events":[],"appEvents":[],"withdrawals":[],"fuel":"0x1"}`),
		watResult(300, `{"state":[],"fuel":"0x1"}`),
		watResult(400, `{"state":[7],"events":[],"appEvents":[],"fuel":"0x1"}`),
		spin)
}

func recycleTestWasm(t *testing.T, maxPages, spin int) []byte {
	t.Helper()
	wasmBytes, err := wasmtime.Wat2Wasm(recycleTestWat(maxPages, spin))
	require.NoError(t, err)
	return wasmBytes
}

func TestProcessRequest_EachCallRunsOnAFreshInstance(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t, 16, 1)
	appId := common.NewApplicationId(1)

	for i := 0; i < 3; i++ {
		state, _, _, _, _, _, failure := runtime.ProcessRequest(context.Background(), appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
		require.Nil(t, failure)
		require.Equal(t, []byte{1}, state, "call %d must not see guest globals of a previous call", i)
	}
}

func TestProcessRequest_GuestMemoryDoesNotAccumulateAcrossCalls(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	// Guest memory capped at 16 MiB and never freed: a single instance fits only two 6 MiB states
	wasmBytes := recycleTestWasm(t, 256, 1)
	appId := common.NewApplicationId(1)
	state := make([]byte, 6<<20)

	for i := 0; i < 20; i++ {
		_, _, _, _, _, _, failure := runtime.ProcessRequest(context.Background(), appId, ethCommon.Address{}, common.Process, []byte("{}"), state, wasmBytes)
		require.Nil(t, failure, "call %d", i)
	}
}

func TestDeposit_GuestMemoryDoesNotAccumulateAcrossCalls(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t, 256, 1)
	appId := common.NewApplicationId(1)
	state := make([]byte, 6<<20)

	for i := 0; i < 20; i++ {
		newState, _, _, _, failure := runtime.Deposit(context.Background(), appId, ethCommon.Address{}, ethCommon.Address{}, big.NewInt(1), state, wasmBytes)
		require.Nil(t, failure, "call %d", i)
		require.Equal(t, []byte{7}, newState)
	}
}

func openFds(t *testing.T) int {
	t.Helper()
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		t.Skip("/proc/self/fd not available")
	}
	return len(entries)
}

func TestRecycledInstancesAreReleased(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t, 16, 1)
	appId := common.NewApplicationId(1)
	call := func() {
		_, _, _, _, _, _, failure := runtime.ProcessRequest(context.Background(), appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
		require.Nil(t, failure)
	}
	call() // module compiled and cached, log pipe open
	baseline := openFds(t)

	// Every store holds the WASI stdout/stderr of the log pipe: unreleased stores would add 2 fds per call
	for i := 0; i < 200; i++ {
		call()
	}
	require.Eventually(t, func() bool { return openFds(t) <= baseline+10 }, 5*time.Second, 50*time.Millisecond,
		"open fds grew from %d to %d", baseline, openFds(t))
}

// TestGuestCallsWaitForTheExecutionLock checks that guest calls are serialized: a request
// retried after a timeout must wait for the previous attempt instead of running beside it.
func TestGuestCallsWaitForTheExecutionLock(t *testing.T) {
	wasmBytes := recycleTestWasm(t, 16, 1)
	ctx := context.Background()
	calls := map[string]func(r *WasmtimeRuntime) error{
		"ProcessRequest": func(r *WasmtimeRuntime) error {
			_, _, _, _, _, _, failure := r.ProcessRequest(ctx, common.NewApplicationId(1), ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
			if failure != nil {
				return failure
			}
			return nil
		},
		"Deposit": func(r *WasmtimeRuntime) error {
			_, _, _, _, failure := r.Deposit(ctx, common.NewApplicationId(1), ethCommon.Address{}, ethCommon.Address{}, big.NewInt(1), []byte("{}"), wasmBytes)
			if failure != nil {
				return failure
			}
			return nil
		},
		"Deploy": func(r *WasmtimeRuntime) error {
			_, _, err := r.Deploy(ctx, common.NewApplicationId(2), nil, wasmBytes)
			return err
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			runtime := NewWasmtimeRuntime(testLogger, 0)
			defer runtime.Close()

			runtime.execLock.Lock() // a call in progress
			done := make(chan error, 1)
			go func() { done <- call(runtime) }()
			select {
			case err := <-done:
				t.Fatalf("%s ran while another guest call held the execution lock (err: %v)", name, err)
			case <-time.After(200 * time.Millisecond):
			}
			runtime.execLock.Unlock()
			select {
			case err := <-done:
				require.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatalf("%s did not run after the execution lock was released", name)
			}
		})
	}
}

// recordingLogger keeps the formatted messages of every level.
type recordingLogger struct {
	mu   sync.Mutex
	msgs []string
}

func (l *recordingLogger) add(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.msgs = append(l.msgs, fmt.Sprintf(msg, args...))
}
func (l *recordingLogger) Trace(msg string, args ...any) { l.add(msg, args...) }
func (l *recordingLogger) Debug(msg string, args ...any) { l.add(msg, args...) }
func (l *recordingLogger) Info(msg string, args ...any)  { l.add(msg, args...) }
func (l *recordingLogger) Warn(msg string, args ...any)  { l.add(msg, args...) }
func (l *recordingLogger) Error(msg string, args ...any) { l.add(msg, args...) }
func (l *recordingLogger) Fatal(msg string, args ...any) { l.add(msg, args...) }
func (l *recordingLogger) Panic(msg string, args ...any) { l.add(msg, args...) }
func (l *recordingLogger) SetLevel(string) error         { return nil }
func (l *recordingLogger) GetLevel() string              { return "trace" }
func (l *recordingLogger) Close() error                  { return nil }

var _ logger.Logger = (*recordingLogger)(nil)

func TestDeposit_DoesNotLogTheGuestResult(t *testing.T) {
	log := &recordingLogger{}
	runtime := NewWasmtimeRuntime(log, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t, 16, 1)

	_, _, _, _, failure := runtime.Deposit(context.Background(), common.NewApplicationId(1), ethCommon.Address{}, ethCommon.Address{}, big.NewInt(1), []byte("{}"), wasmBytes)
	require.Nil(t, failure)

	log.mu.Lock()
	defer log.mu.Unlock()
	for _, m := range log.msgs {
		require.NotContains(t, m, `"state":`, "the deposit result (with the full app state) must not be logged")
	}
}

func TestDeploy_CachesTheModuleForLaterRequests(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t, 16, 1)
	appId := common.NewApplicationId(1)
	ctx := context.Background()

	state, _, err := runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, []byte{}, state)

	// The module is cached: requests work without the wasm bytes
	newState, _, _, _, _, _, failure := runtime.ProcessRequest(ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), nil)
	require.Nil(t, failure)
	require.Equal(t, []byte{1}, newState)

	_, _, err = runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.ErrorContains(t, err, "already deployed")
}

func TestLoadModule_ReturnsTheInitialState(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t, 16, 1)

	state, _, err := runtime.LoadModule(context.Background(), common.NewApplicationId(1), wasmBytes)
	require.NoError(t, err)
	require.Equal(t, []byte{}, state)
}
