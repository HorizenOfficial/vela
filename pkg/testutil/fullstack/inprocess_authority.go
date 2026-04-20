package fullstack

import (
	"crypto/ecdsa"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela-common-go/subgraph"
	"github.com/HorizenOfficial/vela/pkg/authorityservice"
	"github.com/HorizenOfficial/vela/pkg/logger"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// InProcessAuthority runs an AuthorityService HTTP server on an ephemeral
// port backed by an httptest.Server. It generates a fresh secp256k1 key pair
// that tests use when signing /getreport challenges.
//
// The reports and artifacts paths are supplied by the caller (the fullstack
// suite) so the manager and authority service share the same filesystem dirs.
// This matters for DeployApp: the wallet uploads WASM via the authority
// service's /deploy/upload endpoint, and the manager fetches it from the
// shared artifacts dir.
type InProcessAuthority struct {
	server        *httptest.Server
	authorityKey  *ecdsa.PrivateKey
	authorityAddr ethCommon.Address
}

// NewInProcessAuthority builds and starts an authority service on an ephemeral
// port. The server is shut down by Close().
//
//   - chainID       — enforced on /getreport requests (must match what the
//     wallet/test signs).
//   - reportsPath   — filesystem path shared with the manager; manager writes
//     DeanonymizationReport JSON files here, authority reads them.
//   - artifactsPath — filesystem path shared with the manager; /deploy/upload
//     writes WASM artifacts here, manager reads them on deploy processing.
//   - sg            — subgraph client used by /getreport to confirm on-chain
//     completion.
func NewInProcessAuthority(t *testing.T, chainID uint64, reportsPath, artifactsPath string, sg subgraph.Client) *InProcessAuthority {
	t.Helper()

	authorityKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err, "failed to generate authority key")

	testLogger := logger.NewLogger(&logger.Config{
		Kind:         "zerolog",
		Console:      true,
		ConsoleLevel: "info",
	})

	svc, err := authorityservice.NewAuthorityService(
		chainID,
		300*time.Second,
		reportsPath,
		artifactsPath,
		0, // unlimited upload size
		sg,
		testLogger,
	)
	require.NoError(t, err, "failed to build authority service")

	server := httptest.NewServer(svc.Handler())

	return &InProcessAuthority{
		server:        server,
		authorityKey:  authorityKey,
		authorityAddr: ethCrypto.PubkeyToAddress(authorityKey.PublicKey),
	}
}

// URL returns the base URL of the running authority service (e.g. http://127.0.0.1:xxxxx).
func (a *InProcessAuthority) URL() string {
	return a.server.URL
}

// AuthorityAddress returns the Ethereum address derived from the generated
// authority key. Reports stored for this service should list this address as
// their Authority, and /getreport requests must be signed with AuthorityKey().
func (a *InProcessAuthority) AuthorityAddress() ethCommon.Address {
	return a.authorityAddr
}

// AuthorityKey returns the secp256k1 private key used to sign /getreport
// challenges. Hand this to the test-side request builder.
func (a *InProcessAuthority) AuthorityKey() *ecdsa.PrivateKey {
	return a.authorityKey
}

// Close stops the HTTP server. Reports and artifacts paths are owned by the
// suite's core (manager), which cleans them up — Close() does not touch them.
func (a *InProcessAuthority) Close() error {
	if a.server != nil {
		a.server.Close()
	}
	return nil
}
