package wasm

import (
	"fmt"
	"sort"
	"strings"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v47"
)

// The set of host functions a guest is allowed to declare, and why it is a closed
// list rather than an open one.
//
// A host call runs OUTSIDE the execution bound. Epoch interruption works by checks
// the compiler emits into guest code, so it is only reached while wasm is
// executing; a guest that calls into the host and blocks there is off in native
// code, where bumping the epoch does nothing until the call returns. Every lever
// the host has is that same bump — the ticker, the watcher acting on a cancelled
// context or an expired caller budget, and Close. So a blocking import outlasts the
// guest bound, the caller's execution budget and shutdown alike, all while holding
// execLock, which is global: one guest waiting inside poll_oneoff stalls every
// other application for as long as it named, and costs nothing, because a request
// that is never completed is never signed and never charged. TinyGo lowers
// time.Sleep to exactly that import, so it is reachable without any ill intent.
//
// That cannot be enforced from inside the blocked thread, so it is enforced before
// the guest ever runs: the module's import section is checked at load, and anything
// not listed here is refused.
//
// Why refuse a DECLARED import rather than neuter the call. Shadowing the function
// with one that returns an errno would let the module load and change its behaviour
// silently — a guest tested against stock wasmtime would do something different
// inside the enclave, and what "something" is depends on how its toolchain reacts
// to an unexpected errno. Refusing at load is unambiguous, and it is cheap: the
// check runs after compilation but before instantiation, so it lands before a start
// section can execute, on the deploy path and on every cache-miss reload alike. The
// refusal surfaces as CodeFailedLoadingOrGettingModule, which is SIGNED, so the
// request settles on-chain and the queue advances rather than stalling.
//
// The list is the exact import set that TinyGo's wasip1 target emits for the guests
// in this system — app/simple, app/trigger and vela-nova's payment app all declare
// these eight and nothing else. Notably absent are poll_oneoff and fd_read, the two
// calls that can park the thread.
//
// Kept deliberately: clock_time_get and random_get return immediately and every
// guest needs them. They are, however, real sources of non-determinism in a signed
// state transition — worth knowing about if re-execution or multi-executor
// attestation is ever added, since newPinnedEngine canonicalises NaN payloads for
// reproducibility that these two do not preserve.
//
// CHANGING THIS LIST IS AN ABI CHANGE. Adding an entry admits a new host function:
// it must return without blocking and without unbounded work, or it reintroduces
// exactly the hole above — see docs/design/WASM_HOST_ABI.md. Removing an entry
// bricks every already-deployed app that declares it, whose requests then fail
// on-chain on the next cache-miss reload.
const wasiModuleName = "wasi_snapshot_preview1"

var allowedGuestImports = map[string]map[string]struct{}{
	wasiModuleName: {
		"args_get":          {},
		"args_sizes_get":    {},
		"clock_time_get":    {},
		"environ_get":       {},
		"environ_sizes_get": {},
		"fd_write":          {},
		"proc_exit":         {},
		"random_get":        {},
	},
}

// checkGuestImportsAllowed rejects a module that declares any host function outside
// allowedGuestImports. Every offending import is reported at once, so an author
// rebuilding a guest is not fed them one deploy at a time.
//
// Must be called before instantiating: instantiation runs the start section, which
// is guest code, and the point is to refuse the module before it can run anything.
func checkGuestImportsAllowed(module *wasmtime.Module) error {
	var rejected []string
	seen := make(map[string]struct{})

	for _, imp := range module.Imports() {
		namespace := imp.Module()

		// A nil name is a module-level (rather than function-level) import. Nothing in
		// the allowed set has that shape, so it is refused like any other unknown.
		name := "<module import>"
		if n := imp.Name(); n != nil {
			name = *n
		}

		qualified := namespace + "." + name
		if _, ok := allowedGuestImports[namespace][name]; ok {
			continue
		}
		if _, dup := seen[qualified]; dup {
			continue
		}
		seen[qualified] = struct{}{}
		rejected = append(rejected, qualified)
	}

	if len(rejected) == 0 {
		return nil
	}

	// Sorted so the message is stable regardless of the order the module lists them.
	sort.Strings(rejected)
	return fmt.Errorf(
		"module declares host import(s) that are not allowed: %s; a guest may only import %s. "+
			"A host call cannot be interrupted by the guest execution bound, so an import that "+
			"blocks would stall every application; see docs/design/WASM_HOST_ABI.md",
		strings.Join(rejected, ", "), describeAllowedGuestImports())
}

// describeAllowedGuestImports renders the allowed set for the rejection message, so
// the author is told what IS permitted rather than only what is not.
func describeAllowedGuestImports() string {
	var all []string
	for namespace, names := range allowedGuestImports {
		for name := range names {
			all = append(all, namespace+"."+name)
		}
	}
	sort.Strings(all)
	return strings.Join(all, ", ")
}
