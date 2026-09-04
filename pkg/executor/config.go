// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/communication"
	"github.com/magiconair/properties"
)

// Config defines the configuration for the executor application
type Config struct {
	// ChannelType is the type of communication channel to use between manager and executor (tcp / vsock)
	ChannelType string
	// ChannelParams are the parameters for the connection server
	ChannelParams common.ChannelConnectionParams

	// KeySetRecoveryType is the type of recovery mechanism to use for the keyset
	// Type 0: Unsafe (development only) - masterKey stored in plaintext
	// Type 1: AWS KMS - masterKey encrypted with KMS using Nitro attestation
	KeySetRecoveryType common.RecoveryType

	// KMS Configuration (used when KeySetRecoveryType == RecoveryTypeKMS)
	// KMSKeyARN is the ARN of the AWS KMS key used for envelope encryption.
	// Required when KeySetRecoveryType is RecoveryTypeKMS.
	KMSKeyARN string
	// KMSRegion is the AWS region where the KMS key is located.
	// Defaults to "eu-west-1" if not specified.
	KMSRegion string
	// KMSProxyPort is the vsock port where the KMS proxy listens on the parent EC2.
	// Inside a Nitro Enclave, KMS requests are forwarded through this proxy.
	// Defaults to 8000 if not specified.
	KMSProxyPort uint32
	// FuelPricePerUnit is the price of fuel per unit
	FuelPricePerUnit *big.Int
	// MinFeePerRquest is the minimum fee to be paid for each request
	MinFeePerRequest *big.Int

	// LogKind is the type of logger to use (e.g., "zeronetwork", "tcplog", "zerolog").
	LogKind string
	// LogConsole is true if we want output on console
	LogConsole bool
	// LogConsoleLevel is the level of logging for the console
	LogConsoleLevel string
	// LogConsoleColor is true if the log output should be colored
	LogConsoleColor bool
	// LogFileName is the path to the log file. An empty string means no output on log file
	LogFileName string
	// LogFIleLevel is the level of logging for the console
	LogFileLevel string

	// LogChannelParams are the parameters for sending logging to a remote server.
	LogChannelParams common.ChannelConnectionParams

	// LogNetworkLevel is the level of logging for the network
	LogNetworkLevel string

	// CommunicationParams holds parameters for communication between manager and executor
	CommunicationParams common.CommunicationParams

	// MaxCachedModules is the maximum number of WASM modules to keep in the LRU cache.
	// 0 means unlimited.
	MaxCachedModules int

	// MaxGuestMemoryBytes caps the linear memory a single WASM guest can grow to.
	// 0 means the 2 GiB default, which is also the maximum allowed value (the
	// host ABI exchanges guest pointers as signed 32-bit offsets). Note that the
	// worst-case enclave RAM usage is roughly MaxCachedModules * MaxGuestMemoryBytes,
	// so size the two together against the enclave memory budget. Table storage sits
	// outside linear memory and is bounded separately in pkg/wasm (a few tens of MiB
	// per module at most), so it does not change the sizing arithmetic.
	MaxGuestMemoryBytes int64

	// GuestExecutionTimeoutMs bounds the wall-clock duration of a single WASM guest
	// operation. A guest that exceeds it is interrupted and the request fails with a
	// signed, on-chain error, so a stuck or malicious module cannot hang the executor.
	// 0 means the default (common.DefaultGuestExecutionTimeoutMs); there is no
	// "unlimited" setting. It must be below the manager's request timeout minus the
	// reply margin, or the manager gives up while the executor is still running the
	// request; the manager reports that timeout in the handshake and the executor
	// refuses the connection if this does not fit (checkGuestBoundFitsCallerBudget).
	GuestExecutionTimeoutMs int64
}

const confFileName = "executor.conf"

func LoadConfig() (*Config, error) {
	fileProperties := properties.NewProperties()
	if common.FileExists(confFileName) {
		filePropertiesFromFile, err := properties.LoadFile(confFileName, properties.UTF8)
		if err != nil {
			return nil, err
		}
		fileProperties = filePropertiesFromFile
	}

	var channelType = common.GetConfigVar("CHANNEL_TYPE", "vsock", fileProperties)

	var channelServerConnectionParams common.ChannelConnectionParams
	executorServerPort := common.GetConfigVarUint32("EXECUTOR_PORT", 4000, fileProperties)
	var logClientConnectionParams common.ChannelConnectionParams
	logServerPort := common.GetConfigVarUint32("LOG_SERVER_PORT", 5000, fileProperties)
	if channelType == "vsock" {
		// CID 3 is reserved for the parent EC2 instance (where manager runs), CID >= 16 are available for enclaves (where executor runs)
		// CID is not used actually when creating a listening server
		channelServerConnectionParams = common.VSockChannelConnectionParams{Port: executorServerPort}
		// CID and port are both used when connecting to a server
		managerCid := common.GetConfigVarUint32("MANAGER_VSOCK_CID", 3, fileProperties)
		logClientConnectionParams = common.VSockChannelConnectionParams{CID: managerCid, Port: logServerPort}
	} else {
		executorIpHost := common.GetConfigVar("EXECUTOR_IP_HOST", "localhost", fileProperties)
		channelServerConnectionParams = common.TcpChannelConnectionParams{Ip: executorIpHost, Port: executorServerPort}
		logServerIpHost := common.GetConfigVar("LOG_SERVER_IP_HOST", "localhost", fileProperties)
		logClientConnectionParams = common.TcpChannelConnectionParams{Ip: logServerIpHost, Port: logServerPort}
	}

	communicationParams := common.CommunicationParams{
		RequestTimeoutSec: time.Duration(common.GetConfigVarInt64("EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC", 30, fileProperties)),
	}

	// KMS Configuration
	kmsKeyARN := common.GetConfigVar("EXECUTOR_KMS_KEY_ARN", "", fileProperties)
	kmsRegion := common.GetConfigVar("EXECUTOR_KMS_REGION", "eu-west-1", fileProperties)
	kmsProxyPort := common.GetConfigVarUint32("EXECUTOR_KMS_PROXY_PORT", 8000, fileProperties)

	keySetRecoveryType := common.RecoveryType(common.GetConfigVarInt64("EXECUTOR_KEYSET_RECOVERY_TYPE", int64(common.RecoveryTypeKMS), fileProperties))
	if keySetRecoveryType != common.RecoveryTypeUnsafe && keySetRecoveryType != common.RecoveryTypeKMS {
		return nil, fmt.Errorf("invalid EXECUTOR_KEYSET_RECOVERY_TYPE: %d", keySetRecoveryType)
	}

	return &Config{
		ChannelType:         channelType,
		ChannelParams:       channelServerConnectionParams,
		KeySetRecoveryType:  keySetRecoveryType,
		KMSKeyARN:           kmsKeyARN,
		KMSRegion:           kmsRegion,
		KMSProxyPort:        kmsProxyPort,
		FuelPricePerUnit:    big.NewInt(common.GetConfigVarInt64("EXECUTOR_FUEL_PRICE_PER_UNIT", 1, fileProperties)),
		MinFeePerRequest:    big.NewInt(common.GetConfigVarInt64("EXECUTOR_MIN_FEE_PER_REQUEST", 10, fileProperties)),
		LogKind:             common.GetConfigVar("EXECUTOR_LOG_KIND", "zeronetwork", fileProperties),
		LogConsole:          common.GetConfigVarBool("EXECUTOR_LOG_CONSOLE", true, fileProperties),
		LogConsoleLevel:     common.GetConfigVar("EXECUTOR_LOG_CONSOLE_LEVEL", "info", fileProperties),
		LogConsoleColor:     common.GetConfigVarBool("EXECUTOR_LOG_CONSOLE_COLOR", false, fileProperties),
		LogFileName:         common.GetConfigVar("EXECUTOR_LOG_FILE_NAME", "", fileProperties),
		LogFileLevel:        common.GetConfigVar("EXECUTOR_LOG_FILE_LEVEL", "info", fileProperties),
		LogChannelParams:    logClientConnectionParams,
		LogNetworkLevel:     common.GetConfigVar("EXECUTOR_LOG_NETWORK_LEVEL", "info", fileProperties),
		CommunicationParams: communicationParams,
		MaxCachedModules:    int(common.GetConfigVarInt64("EXECUTOR_MAX_CACHED_MODULES", 0, fileProperties)),
		MaxGuestMemoryBytes: common.GetConfigVarInt64("EXECUTOR_MAX_GUEST_MEMORY_BYTES", 0, fileProperties),

		GuestExecutionTimeoutMs: common.GetConfigVarInt64("EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS", 0, fileProperties),
	}, nil
}

// Validate checks the executor configuration for invalid or conflicting settings.
func (c *Config) Validate() error {
	var errs []string

	// --- Channel type ---
	// Only "tcp" and "vsock" are supported. Validate() is called before the
	// logger or WASM runtime are created, so catching this early avoids wasted
	// startup effort and prevents the unguarded type-assertions in the
	// ChannelType switch from panicking.
	if c.ChannelType != "tcp" && c.ChannelType != "vsock" {
		errs = append(errs, fmt.Sprintf(
			"CHANNEL_TYPE must be \"tcp\" or \"vsock\", got %q", c.ChannelType))
	}

	// --- TCP port ranges ---
	// Only meaningful for tcp; a vsock port shares the field but has no 16-bit
	// limit. Skipped when the params are absent or not TCP, so partially populated
	// configs in tests still validate.
	if c.ChannelType == "tcp" {
		if p, ok := c.ChannelParams.(common.TcpChannelConnectionParams); ok {
			if p.Port > common.MaxTCPPort {
				errs = append(errs, fmt.Sprintf(
					"EXECUTOR_PORT must be <= %d for tcp, got %d", common.MaxTCPPort, p.Port))
			}
			// 0 parses cleanly but is never usable here: the executor would listen on an
			// OS-assigned port while the manager dials the configured one from the same
			// variable, so it fails at connect time with nothing wrong at startup.
			// Deliberately not applied to the log server, where 0 means "no TCP log
			// listener" (see logserver.StartLogServer).
			if p.Port == 0 {
				errs = append(errs, "EXECUTOR_PORT must not be 0 for tcp: the manager dials this port")
			}
		}
		if p, ok := c.LogChannelParams.(common.TcpChannelConnectionParams); ok && p.Port > common.MaxTCPPort {
			errs = append(errs, fmt.Sprintf(
				"LOG_SERVER_PORT must be <= %d for tcp, got %d", common.MaxTCPPort, p.Port))
		}
	}

	// --- Fee parameters ---
	// FuelPricePerUnit is used in Mul operations (e.g. totalFuel * FuelPricePerUnit)
	// to compute application fees. A nil value would panic. A zero value effectively
	// disables fuel metering — every fee computation returns 0 and then falls back
	// to MinFeePerRequest, meaning all requests cost the same flat fee regardless
	// of how much computation they consume. A negative value would produce negative
	// fees, which breaks the economic model.
	if c.FuelPricePerUnit == nil || c.FuelPricePerUnit.Sign() <= 0 {
		errs = append(errs, "EXECUTOR_FUEL_PRICE_PER_UNIT must be > 0")
	}

	// MinFeePerRequest is compared against req.MaxFeeValue to reject underfunded
	// requests, and used as a floor for the application fee. A nil value would
	// panic. A negative value means the fee check always passes regardless of what
	// the user sends — a security concern since the system would process requests
	// without collecting fees. Zero is legal: it means there is no minimum fee.
	if c.MinFeePerRequest == nil || c.MinFeePerRequest.Sign() < 0 {
		errs = append(errs, "EXECUTOR_MIN_FEE_PER_REQUEST must be >= 0")
	}

	// --- Communication timeout ---
	// RequestTimeoutSec feeds time.Now().Add(Duration * time.Second) for pending
	// request deadlines. A zero value means every request's deadline is already
	// in the past at creation time — the response-routing loop evicts it
	// immediately, so every request fails with a timeout error.
	if c.CommunicationParams.RequestTimeoutSec <= 0 {
		errs = append(errs, fmt.Sprintf(
			"EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC must be > 0 (seconds), got %d",
			c.CommunicationParams.RequestTimeoutSec))
	}

	// --- Guest memory cap ---
	// The WASM host ABI exchanges guest pointers as signed 32-bit offsets, so a
	// guest may never grow past 2 GiB. The ceiling lives in pkg/common rather than
	// pkg/wasm so this check shares one definition with the runtime without linking
	// libwasmtime into every pkg/executor consumer. 0 selects the default; anything
	// else must be in range rather than silently clamped, so a misconfigured
	// operator finds out at startup.
	const maxGuestMemoryCeilingBytes = common.MaxGuestMemoryCeilingBytes
	if c.MaxGuestMemoryBytes < 0 || c.MaxGuestMemoryBytes > maxGuestMemoryCeilingBytes {
		errs = append(errs, fmt.Sprintf(
			"EXECUTOR_MAX_GUEST_MEMORY_BYTES must be between 0 (default: 2 GiB) and %d (2 GiB), got %d",
			int64(maxGuestMemoryCeilingBytes), c.MaxGuestMemoryBytes))
	}

	// --- Guest execution timeout ---
	// Bounds live in pkg/common so this check shares one definition with the runtime
	// without linking libwasmtime into every pkg/executor consumer, exactly as for the
	// guest memory cap. 0 selects the default; anything else must be in range rather
	// than silently clamped, so a misconfigured operator finds out at startup.
	//
	// Note the asymmetry with the runtime's SetGuestExecutionTimeout, which clamps
	// instead: failing here happens before the enclave is serving anything, whereas a
	// late admin-driven change must not be able to take a running enclave down.
	if c.GuestExecutionTimeoutMs < 0 || c.GuestExecutionTimeoutMs > common.MaxGuestExecutionTimeoutMs {
		errs = append(errs, fmt.Sprintf(
			"EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS must be between 0 (default: %d ms) and %d ms, got %d",
			int64(common.DefaultGuestExecutionTimeoutMs), int64(common.MaxGuestExecutionTimeoutMs),
			c.GuestExecutionTimeoutMs))
	}

	// The bound must also stay below the manager's execution budget, but the manager's
	// request timeout is configured on the parent EC2, out of the enclave's sight. The
	// manager reports it in the handshake, and checkGuestBoundFitsCallerBudget rejects
	// the connection there if the two do not fit; see that method for the invariant.

	// --- KMS configuration ---
	errs = append(errs, c.validateKMSConfig()...)

	if len(errs) > 0 {
		return fmt.Errorf("configuration errors:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

// validateKMSConfig validates KMS-related configuration.
// Returns all errors found so the caller can present them at once.
func (c *Config) validateKMSConfig() []string {
	if c.KeySetRecoveryType != common.RecoveryTypeKMS {
		return nil
	}
	var errs []string
	if c.KMSKeyARN == "" {
		errs = append(errs, "KMS key ARN is required when KMS is enabled (EXECUTOR_KMS_KEY_ARN)")
	}
	if c.KMSRegion == "" {
		errs = append(errs, "KMS region is required when KMS is enabled (EXECUTOR_KMS_REGION)")
	}
	return errs
}

// guestExecutionBound returns the configured per-operation guest execution bound,
// resolving 0 to the default. It is what the batch loop compares the remaining
// caller budget against, so it must reflect the value the runtime actually arms
// stores with rather than the raw config field.
func (c *Config) guestExecutionBound() time.Duration {
	ms := c.GuestExecutionTimeoutMs
	if ms <= 0 {
		ms = common.DefaultGuestExecutionTimeoutMs
	}
	return time.Duration(ms) * time.Millisecond
}

// guestBoundSafetyMargin is the headroom the guest execution bound must leave inside
// the caller's execution budget, on top of the bound itself.
//
// The bound is not enforced to the millisecond, and everything it fails to cover is
// time the caller is already counting. Three things have to fit in here:
//
//   - Epoch overshoot. A window armed with T is interrupted somewhere in
//     [T, T+2 ticks]: epochTicksFor rounds the tick count up and adds one more, so a
//     store armed just before a tick fires still gets its full allowance.
//   - Work before arming. The caller starts its clock when it sends the request; the
//     executor arms the store only after decoding the message, decrypting the
//     application state and hashing the module. For a multi-megabyte module and a
//     large state that is tens to hundreds of milliseconds.
//   - The floor on an exhausted budget. Once a request has spent its allowance,
//     epochTicksFor still arms the minimum two ticks, so a further operation can run
//     that long before it traps.
//
// One second covers all three with room to spare, and costs nothing: it only forbids
// configurations in which the two deadlines were going to race anyway. It is a fixed
// amount rather than a fraction of the timeout because every term above is
// independent of how patient the caller is.
const guestBoundSafetyMargin = time.Second

// checkGuestBoundFitsCallerBudget verifies that the guest execution bound fires before
// the caller's execution budget expires. callerTimeoutMs is the manager's request
// timeout (MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC) as reported in the
// handshake; the budget the executor will actually receive on every request is that
// value minus communication.ExecutionBudgetMargin, reserved for decrypting, signing
// and replying.
//
// This is a forward-progress invariant, not a tuning hint, which is why it is a hard
// error. Whichever deadline fires first decides how a runaway guest is classified. The
// enclave's own bound expiring is a signed on-chain failure, so the FIFO request queue
// advances and the request's author pays the minimum fee. The caller's budget expiring
// is host-side abandonment: nothing is signed and the request is retried unchanged. If
// the budget wins the race, a non-terminating guest is retried on every poll forever,
// the queue head never advances, and — since nothing is ever signed — nobody is charged
// for the executor time it keeps burning. That is the stall signing a timeout exists to
// prevent, reintroduced purely by configuration, and it turns a fee-paying nuisance
// into a free denial of service.
//
// It compares against the margin too, not just the raw timeout: a bare
// `>= requestTimeout` check leaves a silent window one margin wide (28000 to 29999 ms
// against the 30 s default) that passes and still wedges.
//
// It runs at handshake rather than at start-up because only the manager knows its own
// timeout. An earlier version compared against the executor's request timeout
// (EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC) as a stand-in; that setting
// governs calls in the other direction and gave false assurance whenever the two were
// configured apart.
func (c *Config) checkGuestBoundFitsCallerBudget(callerTimeoutMs int64) error {
	guestBoundMs := c.guestExecutionBound().Milliseconds()

	// callerTimeoutMs crosses the TEE boundary and, unlike the per-request
	// ExecutionBudgetMs, is not range-checked anywhere else. Guard both ends before the
	// subtraction below. A negative value would underflow budgetMs to a large positive
	// number and pass the comparison, silently disarming this invariant; a value above
	// communication.MaxExecutionBudgetMs would pass here only for every later request to
	// be rejected by validateExecutionBudget, stalling the queue forever. Refuse both at
	// the handshake, where the operator is told once, rather than on every request.
	if callerTimeoutMs <= 0 || callerTimeoutMs > communication.MaxExecutionBudgetMs {
		return fmt.Errorf(
			"the manager's request timeout (MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC = %d ms) must be "+
				"positive and at most %d ms; a value outside that range cannot yield a usable per-request "+
				"execution budget",
			callerTimeoutMs, int64(communication.MaxExecutionBudgetMs))
	}

	// One request costs at most one bound of guest execution, whatever it does —
	// loading the module, taking a deposit and processing all draw down a single
	// per-request budget (see common.WithGuestExecutionBudget). So the whole
	// requirement is that one bound, plus the slack the bound does not itself cover,
	// still fits inside what the caller will wait.
	budgetMs := callerTimeoutMs - communication.ExecutionBudgetMargin.Milliseconds()
	requiredMs := guestBoundMs + guestBoundSafetyMargin.Milliseconds()
	if requiredMs >= budgetMs {
		return fmt.Errorf(
			"EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS (%d ms, effective) plus the %d ms safety margin must be below "+
				"the manager's execution budget of %d ms (MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC = "+
				"%d ms minus the %d ms reply margin), otherwise a runaway guest is retried forever instead of "+
				"being settled on-chain and the request queue stalls",
			guestBoundMs, guestBoundSafetyMargin.Milliseconds(), budgetMs, callerTimeoutMs,
			communication.ExecutionBudgetMargin.Milliseconds())
	}
	return nil
}
