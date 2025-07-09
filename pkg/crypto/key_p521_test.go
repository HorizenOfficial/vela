package crypto

import (
	"bytes"
	"os"
	"testing"
)

func TestGeneratePrivateKeyP521(t *testing.T) {
	key, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	if key == nil {
		t.Fatal("key is nil")
	}
}

func TestSaveLoadPrivateKeyP521DER(t *testing.T) {
	key, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	filename := "test_key.der"
	defer os.Remove(filename)

	err = SavePrivateKeyP521ToFileDER(key, filename)
	if err != nil {
		t.Fatalf("failed to save private key: %v", err)
	}

	loadedKey, err := LoadPrivateKeyP521FromFileDER(filename)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	if !bytes.Equal(key.Bytes(), loadedKey.Bytes()) {
		t.Fatal("loaded key does not match original key")
	}
}

func TestSaveLoadPrivateKeyP521PEM(t *testing.T) {
	key, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	filename := "test_key.pem"
	defer os.Remove(filename)

	err = SavePrivateKeyP521ToFilePEM(key, filename)
	if err != nil {
		t.Fatalf("failed to save private key: %v", err)
	}

	loadedKey, err := LoadPrivateKeyP521FromFilePEM(filename)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	if !bytes.Equal(key.Bytes(), loadedKey.Bytes()) {
		t.Fatal("loaded key does not match original key")
	}
}
