package version

// Version is set at build time via:
//
//	-ldflags "-X github.com/horizen-pes/pkg/version.Version=v1.2.3"
//
// Falls back to "dev" when built without -ldflags (local development).
var Version = "dev"
