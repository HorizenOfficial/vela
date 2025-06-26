package main

// Config defines the configuration for the executor application
type Config struct {
	// ServerType is the type of server to use (http or vsock)
	ServerType string
	// ServerAddr is the address for the HTTP server
	ServerAddr string
	// ServerPort is the port for the v-socket server
	ServerPort uint32
	// SigningKey is the key to use for signing update payloads
	SigningKey []byte
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		ServerType: "tcp",
		ServerAddr: "localhost:8080",
		ServerPort: 5000,
		SigningKey: []byte("dummy-signing-key"),
	}
}

func main() {

}
