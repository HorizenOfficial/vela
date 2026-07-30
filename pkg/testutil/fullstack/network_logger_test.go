package fullstack_test

import (
	"os"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/executor"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/manager"
	"github.com/HorizenOfficial/vela/pkg/testutil/fullstack"
	"github.com/stretchr/testify/require"
)

// TestFullStackSuiteWithNetworkLogger builds the fullstack suite with a
// "zeronetwork" logger config, the way consumers of this test helper do when they
// want manager/executor logs to reach the log server.
//
// Regression test: the suite used to build a logger from the manager's log config
// before NewTestSuiteCore had injected the ephemeral log-server address into it
// (RemoteLogParams), so a network logger panicked during construction with
// "Invalid RemoteLogParams for tcp: <nil>". Every test in this package passes a
// console config, which ignores RemoteLogParams harmlessly, so the bug was
// invisible here and only surfaced downstream.
//
// Constructing the suite is the whole assertion — it panics if the ordering
// regresses.
func TestFullStackSuiteWithNetworkLogger(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("skipping fullstack test under CI_FLAG")
	}

	// manager.LoadConfig requires MANAGER_ARTIFACTS_PATH; TestSuiteCore overrides
	// it with a temp dir during construction.
	t.Setenv("MANAGER_ARTIFACTS_PATH", t.TempDir())
	// Unsafe keyset recovery (type 0) so GenerateEnclaveKeySet needs no KMS.
	t.Setenv("EXECUTOR_KEYSET_RECOVERY_TYPE", "0")

	mgrCfg, err := manager.LoadConfig()
	require.NoError(t, err)
	execCfg, err := executor.LoadConfig()
	require.NoError(t, err)
	keySet, recovery, err := executor.GenerateEnclaveKeySet(t.Context(), execCfg.KeySetRecoveryType, nil, nil, "")
	require.NoError(t, err)

	// Deliberately no RemoteLogParams: the suite is responsible for injecting the
	// ephemeral log-server address, since only it knows the port.
	netLogCfg := func() *logger.Config {
		return &logger.Config{
			Kind:             "zeronetwork",
			Console:          true,
			ConsoleColor:     false,
			ConsoleLevel:     "error",
			FileLevel:        "error",
			NetworkLevel:     "error",
			RemoteLogNetwork: "tcp",
		}
	}

	mgrLogCfg, excLogCfg := netLogCfg(), netLogCfg()
	require.Nil(t, mgrLogCfg.RemoteLogParams, "precondition: the caller does not know the log-server port")

	suite := fullstack.NewFullStackSystemTestSuiteWithConfigs(
		t, "mock-runtime", mgrCfg, execCfg, keySet, recovery, mgrLogCfg, excLogCfg)
	require.NotNil(t, suite)
	defer suite.Cleanup()

	// Construction must have injected the log-server address into both configs.
	// This is the contract ManagerLogger() depends on, and what any logger built
	// from these configs needs in order not to panic.
	require.NotNil(t, mgrLogCfg.RemoteLogParams, "suite must inject the log-server address into the manager log config")
	require.NotNil(t, excLogCfg.RemoteLogParams, "suite must inject the log-server address into the executor log config")
}
