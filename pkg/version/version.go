package version

// Version is set at build time via:
//
//	-ldflags "-X github.com/HorizenOfficial/vela/pkg/version.Version=v1.2.3"
//
// Falls back to "dev" when built without -ldflags (local development).
var Version = "dev"
