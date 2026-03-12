package manager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig_UsesConfiguredArtifactAndReportPaths(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir := t.TempDir()
	reportsDir := filepath.Join(tmpDir, "reports")
	artifactsDir := filepath.Join(tmpDir, "artifacts")

	configContents := "MANAGER_REPORTS_FOLDER=" + reportsDir + "\n" +
		"MANAGER_ARTIFACTS_PATH=" + artifactsDir + "\n"
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, confFileName), []byte(configContents), 0o600))

	require.NoError(t, os.Chdir(tmpDir))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(wd))
	})

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.Equal(t, artifactsDir, cfg.ArtifactsPath)
	require.Equal(t, reportsDir, cfg.DeanonymizationReportPath)
}
