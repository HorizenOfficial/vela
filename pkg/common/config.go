package common

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/magiconair/properties"
)

/*
* Get the string vale of a configuration variable, by using the following order (the first value found is returned):
* - ENV variables
* - properties passed in the fileProperties arguments (usually loaded via config file)
* - defualtValue passed
 */
func GetConfigVar(name string, defaultValue string, fileProperties *properties.Properties) string {
	valFromEnv := os.Getenv(name)
	if valFromEnv != "" {
		return valFromEnv
	} else {
		return fileProperties.GetString(name, defaultValue)
	}
}

/*
* Same as GetConfigVar method, but the value is converted in int64.
* In case of conversion errors, the default value is returned.
 */
func GetConfigVarInt64(name string, defaultValue int64, fileProperties *properties.Properties) int64 {
	var confVar = GetConfigVar(name, "", fileProperties)
	if confVar == "" {
		return defaultValue
	} else {
		// Parse the full int64 range: values above 2^31 must reach the caller
		// (e.g. so config validation can reject an out-of-range value) instead
		// of silently falling back to the default.
		var parsed, err = strconv.ParseInt(confVar, 10, 64)
		if err != nil {
			fmt.Printf("Failed to convert %v for error %v, using default value\n", name, err)
			return defaultValue
		}
		return parsed
	}
}

// MaxGuestMemoryCeilingBytes is both the default and the maximum per-guest linear
// memory cap. It is a hard ABI constraint rather than a tuning choice: the WASM host
// exchanges guest pointers as signed int32 offsets, so no guest offset may reach
// 2 GiB (see pkg/wasm writeToMemory / extractResultBytes).
//
// It lives here, in a package with no heavy dependencies, so that pkg/wasm (which
// links libwasmtime) and pkg/executor (which must not) can share one definition
// instead of keeping copies in step by hand.
const MaxGuestMemoryCeilingBytes = 2 * 1024 * 1024 * 1024

// Wall-clock bounds on a single WASM guest operation, in milliseconds. Guest
// execution is interrupted once it exceeds the configured timeout (see pkg/wasm,
// epoch interruption), so that a stuck or malicious module cannot hang the
// executor. These live here, alongside MaxGuestMemoryCeilingBytes and for the same
// reason: pkg/executor validates the configured value at start-up without linking
// libwasmtime.
const (
	// DefaultGuestExecutionTimeoutMs is used when EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS
	// is unset or 0. It must stay well under the manager's request timeout
	// (MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC, 30s by default), or the
	// manager gives up on a request while the executor is still running it and the
	// two end up permanently out of step; the executor enforces this at handshake,
	// where the manager reports that timeout. It is also orders of magnitude above what
	// a real request needs — the reference guests in app/simple return in single-digit
	// milliseconds — which is what makes exceeding it evidence of a pathological
	// guest rather than of a slow one.
	DefaultGuestExecutionTimeoutMs = 10_000

	// MaxGuestExecutionTimeoutMs is an absurdity guard, not an ABI constraint (unlike
	// MaxGuestMemoryCeilingBytes): any value this large already defeats the purpose of
	// having a bound, since the manager stops waiting long before it. A configured
	// value above it is a misconfiguration worth failing at start-up for.
	//
	// There is deliberately no "unlimited" setting: 0 selects the default. An
	// unbounded guest is the failure mode this exists to remove.
	MaxGuestExecutionTimeoutMs = 300_000
)

// GuestExecutionEpochTick is the granularity of the guest execution bound: the
// period at which pkg/wasm advances the engine epoch, and therefore the resolution
// of every deadline armed from it. pkg/wasm derives its own epochTickInterval from
// this, so the two cannot drift.
//
// It is exported because the bound is not enforced to the millisecond and callers
// have to budget for that. A window armed with T is interrupted somewhere in
// [T, T+2*GuestExecutionEpochTick]: the tick count is rounded UP, and one further
// tick is added so that a store armed just before a tick fires still gets its full
// allowance. pkg/executor sizes the safety margin in its handshake check from this
// value rather than from a copied constant.
const GuestExecutionEpochTick = 100 * time.Millisecond

// MaxTCPPort is the highest valid TCP port number. GetConfigVarUint32 keeps an
// out-of-range value from truncating to 0, but a value that fits in uint32 and is
// still not a port (70000, say) would only fail later at bind time with a generic
// error — so configs that carry TCP ports should range-check them in Validate().
// vsock is unaffected: it legitimately uses the full uint32 range for CIDs and ports.
const MaxTCPPort = 65535

/*
* Same as GetConfigVar method, but the value is converted to uint32 — the type
* used for TCP ports and vsock CIDs. In case of conversion errors, including a
* value that does not fit in 32 unsigned bits, the default value is returned.
*
* Use this rather than narrowing GetConfigVarInt64 with a uint32 cast: the cast
* truncates instead of failing, and a truncated port silently becomes 0, which
* makes a listener bind an arbitrary OS-assigned port instead of reporting a bad
* configuration.
 */
func GetConfigVarUint32(name string, defaultValue uint32, fileProperties *properties.Properties) uint32 {
	var confVar = GetConfigVar(name, "", fileProperties)
	if confVar == "" {
		return defaultValue
	}
	var parsed, err = strconv.ParseUint(confVar, 10, 32)
	if err != nil {
		fmt.Printf("Failed to convert %v for error %v, using default value\n", name, err)
		return defaultValue
	}
	return uint32(parsed)
}

/*
* Same as GetConfigVar method, but the value is converted in boolean.
* In case of conversion errors, the default value is returned.
 */
func GetConfigVarBool(name string, defaultValue bool, fileProperties *properties.Properties) bool {
	var confVar = GetConfigVar(name, "", fileProperties)
	if confVar == "" {
		return defaultValue
	} else {
		var parsed, err = strconv.ParseBool(confVar)
		if err != nil {
			fmt.Printf("Failed to convert %v for error %v, using default value\n", name, err)
			return defaultValue
		}
		return parsed
	}
}
