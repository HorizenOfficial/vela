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
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	"github.com/HorizenOfficial/vela/pkg/logger"
	wasmtime "github.com/bytecodealliance/wasmtime-go/v47"
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

// recycleTestWat is a guest that behaves like a TinyGo guest on large inputs: its allocate
// never frees (bump allocator growing memory on demand, returning 0 when the store refuses
// to grow). Every export counts its calls in a global, so a result tells whether the
// instance had already served a call: process_request returns state [1] on the first
// export call of an instance and [2] afterwards, deposit [7] and [8]. deploy only counts.
const recycleTestWatTemplate = `(module
  (memory (export "memory") 1)
  (global $top (mut i32) (i32.const 4096))
  (global $calls (mut i32) (i32.const 0))
  %s
  %s
  %s
  %s
  %s
  (func $first (result i32)
    (global.set $calls (i32.add (global.get $calls) (i32.const 1)))
    (i32.eq (global.get $calls) (i32.const 1)))
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
  (func (export "deploy") (param i64 i32 i32) (result i32)
    (drop (call $first))
    (i32.const 300))
  (func (export "deposit") (param i64 i32 i32 i32 i32 i32 i32 i32 i32) (result i32)
    (if (result i32) (call $first) (then (i32.const 400)) (else (i32.const 500))))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    (if (result i32) (call $first) (then (i32.const 100)) (else (i32.const 200))))
)`

func recycleTestWasm(t *testing.T) []byte {
	t.Helper()
	wat := fmt.Sprintf(recycleTestWatTemplate,
		watResult(100, `{"state":[1],"events":[],"appEvents":[],"withdrawals":[],"fuel":"0x1"}`),
		watResult(200, `{"state":[2],"events":[],"appEvents":[],"withdrawals":[],"fuel":"0x1"}`),
		watResult(300, `{"state":[],"fuel":"0x1"}`),
		watResult(400, `{"state":[7],"events":[],"appEvents":[],"fuel":"0x1"}`),
		watResult(500, `{"state":[8],"events":[],"appEvents":[],"fuel":"0x1"}`),
	)
	wasmBytes, err := wasmtime.Wat2Wasm(wat)
	require.NoError(t, err)
	return wasmBytes
}

// recycleTestMemoryCap fits two 6 MiB states in one never-freeing guest, not three.
const (
	recycleTestMemoryCap = 16 << 20
	recycleTestStateSize = 6 << 20
)

func processRecycleTest(r *WasmtimeRuntime, appId common.ApplicationIdType, state, wasm []byte) ([]byte, *apperrors.RequestFailure) {
	newState, _, _, _, _, _, failure := r.ProcessRequest(context.Background(), appId, ethCommon.Address{}, common.Process, []byte("{}"), state, wasm)
	return newState, failure
}

func TestProcessRequest_EachCallRunsOnAFreshInstance(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t)
	appId := common.NewApplicationId(1)

	for i := 0; i < 3; i++ {
		state, failure := processRecycleTest(runtime, appId, []byte("{}"), wasmBytes)
		require.Nil(t, failure)
		require.Equal(t, []byte{1}, state, "call %d must not see the guest globals of a previous call", i)
	}
}

func TestProcessRequest_GuestMemoryDoesNotAccumulateAcrossCalls(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	runtime.SetMaxGuestMemoryBytes(recycleTestMemoryCap)
	wasmBytes := recycleTestWasm(t)
	appId := common.NewApplicationId(1)
	state := make([]byte, recycleTestStateSize)

	for i := 0; i < 6; i++ {
		_, failure := processRecycleTest(runtime, appId, state, wasmBytes)
		require.Nil(t, failure, "call %d: guest memory left by previous calls must not count against this one", i)
	}
}

func TestDeposit_GuestMemoryDoesNotAccumulateAcrossCalls(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	runtime.SetMaxGuestMemoryBytes(recycleTestMemoryCap)
	wasmBytes := recycleTestWasm(t)
	appId := common.NewApplicationId(1)
	state := make([]byte, recycleTestStateSize)

	for i := 0; i < 6; i++ {
		newState, _, _, _, failure := runtime.Deposit(context.Background(), appId, ethCommon.Address{}, ethCommon.Address{}, big.NewInt(1), state, wasmBytes)
		require.Nil(t, failure, "call %d", i)
		require.Equal(t, []byte{7}, newState, "call %d must run on a fresh instance", i)
	}
}

// TestGuestMemoryCapAppliesToTheNextCall is the exclusion half of the per-call instance:
// a fresh instance must still be created under the memory cap, including a cap lowered
// after the app's module was compiled and cached.
func TestGuestMemoryCapAppliesToTheNextCall(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t)
	appId := common.NewApplicationId(1)

	_, failure := processRecycleTest(runtime, appId, []byte("{}"), wasmBytes)
	require.Nil(t, failure)
	require.True(t, isCached(runtime, appId))

	runtime.SetMaxGuestMemoryBytes(recycleTestMemoryCap)
	_, failure = processRecycleTest(runtime, appId, make([]byte, recycleTestMemoryCap+1), wasmBytes)
	require.NotNil(t, failure, "a state larger than the guest memory cap must not fit in a fresh instance")
	require.Equal(t, apperrors.CodeMemoryWriteError, failure.RequestError)

	state, failure := processRecycleTest(runtime, appId, []byte("{}"), wasmBytes)
	require.Nil(t, failure, "the app keeps working after a call refused by the cap")
	require.Equal(t, []byte{1}, state)
}

func TestDeploy_CachesTheModuleButNotTheInstance(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t)
	appId := common.NewApplicationId(1)
	ctx := context.Background()

	state, _, err := runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, []byte{}, state)

	// The module is cached: requests work without the wasm bytes, on a fresh instance
	newState, failure := processRecycleTest(runtime, appId, []byte("{}"), nil)
	require.Nil(t, failure)
	require.Equal(t, []byte{1}, newState, "the request must not run on the instance that ran deploy")

	_, _, err = runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.ErrorContains(t, err, "already deployed")
}

func openFds(t *testing.T) int {
	t.Helper()
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		t.Skip("/proc/self/fd not available")
	}
	return len(entries)
}

// TestRecycledInstancesAreReleased checks that per-call stores are closed: each one opens
// the app's log pipe for stdout and stderr, so an unreleased store adds 2 fds per call.
func TestRecycledInstancesAreReleased(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t)
	appId := common.NewApplicationId(1)

	call := func() {
		_, failure := processRecycleTest(runtime, appId, []byte("{}"), wasmBytes)
		require.Nil(t, failure)
	}
	call() // module compiled and cached, log pipe open
	baseline := openFds(t)

	for i := 0; i < 200; i++ {
		call()
	}
	require.LessOrEqual(t, openFds(t), baseline+10, "open fds grew from %d", baseline)
}

// TestGuestCallsWaitForTheExecutionLock checks that guest calls are serialized: a request
// retried after a timeout must wait for the previous attempt instead of running beside it.
func TestGuestCallsWaitForTheExecutionLock(t *testing.T) {
	wasmBytes := recycleTestWasm(t)
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

// TestGuestOutputStillReachesTheLog checks that every per-call instance writes to the
// app's log pipe, which outlives the instances.
func TestGuestOutputStillReachesTheLog(t *testing.T) {
	log := &recordingLogger{}
	runtime := NewWasmtimeRuntime(log, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(fmt.Sprintf(`(module
  (import "wasi_snapshot_preview1" "fd_write" (func $fd_write (param i32 i32 i32 i32) (result i32)))
  (memory (export "memory") 1)
  %s
  (data (i32.const 600) "INF recycled guest says hi\n")
  (data (i32.const 640) "\58\02\00\00\1b\00\00\00")
  (func (export "allocate") (param i32) (result i32) i32.const 4096)
  (func (export "deallocate") (param i32 i32))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    (drop (call $fd_write (i32.const 1) (i32.const 640) (i32.const 1) (i32.const 700)))
    (i32.const 100))
)`, watResult(100, `{"state":[1],"events":[],"appEvents":[],"withdrawals":[],"fuel":"0x1"}`)))
	require.NoError(t, err)

	appId := common.NewApplicationId(1)
	for i := 0; i < 2; i++ {
		_, failure := processRecycleTest(runtime, appId, []byte("{}"), wasmBytes)
		require.Nil(t, failure)
	}
	require.Eventually(t, func() bool { return log.count("recycled guest says hi") == 2 },
		5*time.Second, 20*time.Millisecond, "each call's guest output must reach the host log")
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

func (l *recordingLogger) count(substr string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := 0
	for _, m := range l.msgs {
		if strings.Contains(m, substr) {
			n++
		}
	}
	return n
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

// TestGuestResultsAreNotLogged checks that the deploy and deposit results, which carry
// the full application state, never reach the host log.
func TestGuestResultsAreNotLogged(t *testing.T) {
	log := &recordingLogger{}
	runtime := NewWasmtimeRuntime(log, 0)
	defer runtime.Close()
	wasmBytes := recycleTestWasm(t)
	appId := common.NewApplicationId(1)

	_, _, err := runtime.Deploy(context.Background(), appId, nil, wasmBytes)
	require.NoError(t, err)
	_, _, _, _, failure := runtime.Deposit(context.Background(), appId, ethCommon.Address{}, ethCommon.Address{}, big.NewInt(1), []byte("{}"), wasmBytes)
	require.Nil(t, failure)

	require.Zero(t, log.count(`"state":`), "a guest result (with the full app state) must not be logged")
}
