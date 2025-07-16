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

func TestExportImportPrivateKeyP521Hex(t *testing.T) {
	key, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	hexKey := ExportPrivateKeyP521ToHex(key)

	loadedKey, err := ImportPrivateKeyP521FromHex(hexKey)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	if !bytes.Equal(key.Bytes(), loadedKey.Bytes()) {
		t.Fatal("loaded key does not match original key")
	}
}

func TestExportPublicKeyP521Hex(t *testing.T) {
	key, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	hexPubKey := ExportPublicKeyP521ToHex(key.PublicKey())

	// The public key should start with 04, followed by 66 bytes for X and 66 bytes for Y coordinates
	if len(hexPubKey) != 266 {
		t.Fatalf("public key has wrong length: got %d, want 266", len(hexPubKey))
	}
}
