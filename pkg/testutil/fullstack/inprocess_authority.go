package fullstack

import (
	"crypto/ecdsa"
	"fmt"
	"net/http/httptest"
	"os"
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
// that tests use when signing /getreport challenges, and owns a temp dir for
// the deploy-artifact store required by NewAuthorityService.
//
// The reports directory is supplied by the caller (the fullstack suite) so the
// manager and authority service share the same filesystem path for
// deanonymization reports.
type InProcessAuthority struct {
	server           *httptest.Server
	authorityKey     *ecdsa.PrivateKey
	authorityAddr    ethCommon.Address
	artifactsDir     string
}

// NewInProcessAuthority builds and starts an authority service on an ephemeral
// port. Call Close() to shut down and remove the artifacts temp dir.
//
//   - chainID      — the chain ID the service will enforce on /getreport requests
//     (must match what the wallet/test signs).
//   - reportsPath  — filesystem path where DeanonymizationReport JSON files live.
//     The manager writes them here; the authority service reads them.
//   - sg           — subgraph client used by /getreport to confirm on-chain
//     completion before serving the report.
func NewInProcessAuthority(t *testing.T, chainID uint64, reportsPath string, sg subgraph.Client) *InProcessAuthority {
	t.Helper()

	authorityKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err, "failed to generate authority key")

	artifactsDir, err := os.MkdirTemp("", "fullstack-authority-artifacts")
	require.NoError(t, err)

	testLogger := logger.NewLogger(&logger.Config{
		Kind:         "zerolog",
		Console:      true,
		ConsoleLevel: "info",
	})

	svc, err := authorityservice.NewAuthorityService(
		chainID,
		300*time.Second,
		reportsPath,
		artifactsDir,
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
		artifactsDir:  artifactsDir,
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

// Close stops the HTTP server and removes the artifacts temp dir.
func (a *InProcessAuthority) Close() error {
	if a.server != nil {
		a.server.Close()
	}
	if a.artifactsDir != "" {
		if err := os.RemoveAll(a.artifactsDir); err != nil {
			return fmt.Errorf("failed to remove artifacts dir: %w", err)
		}
	}
	return nil
}
