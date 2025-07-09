package crypto

import (
	"bytes"
	"os"
	"testing"
)

func TestGeneratePrivateKey25519(t *testing.T) {
	key, err := GeneratePrivateKey25519()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	if key == nil {
		t.Fatal("key is nil")
	}
}

func TestSaveLoadPrivateKey25519DER(t *testing.T) {
	key, err := GeneratePrivateKey25519()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	filename := "test_key_25519.der"
	defer os.Remove(filename)

	err = SavePrivateKey25519ToFileDER(key, filename)
	if err != nil {
		t.Fatalf("failed to save private key: %v", err)
	}

	loadedKey, err := LoadPrivateKey25519FromFileDER(filename)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	if !bytes.Equal(key.Bytes(), loadedKey.Bytes()) {
		t.Fatal("loaded key does not match original key")
	}
}

func TestSaveLoadPrivateKey25519PEM(t *testing.T) {
	key, err := GeneratePrivateKey25519()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	filename := "test_key_25519.pem"
	defer os.Remove(filename)

	err = SavePrivateKey25519ToFilePEM(key, filename)
	if err != nil {
		t.Fatalf("failed to save private key: %v", err)
	}

	loadedKey, err := LoadPrivateKey25519FromFilePEM(filename)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	if !bytes.Equal(key.Bytes(), loadedKey.Bytes()) {
		t.Fatal("loaded key does not match original key")
	}
}
