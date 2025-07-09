package crypto

import (
	"bytes"
	"os"
	"testing"
)

func TestSaveLoadAESKey(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("failed to generate aes key: %v", err)
	}

	filename := "test_key.key"
	defer os.Remove(filename)

	err = SaveAESKeyToFile(&key, filename)
	if err != nil {
		t.Fatalf("failed to save aes key: %v", err)
	}

	loadedKey, err := LoadAESKeyFromFile(filename)
	if err != nil {
		t.Fatalf("failed to load aes key: %v", err)
	}

	if !bytes.Equal(key[:], loadedKey[:]) {
		t.Fatal("loaded key does not match original key")
	}
}

func TestSaveLoadAESKeyPEM(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("failed to generate aes key: %v", err)
	}

	filename := "test_key.pem"
	defer os.Remove(filename)

	err = SaveAESKeyToFilePEM(&key, filename)
	if err != nil {
		t.Fatalf("failed to save aes key: %v", err)
	}

	loadedKey, err := LoadAESKeyFromFilePEM(filename)
	if err != nil {
		t.Fatalf("failed to load aes key: %v", err)
	}

	if !bytes.Equal(key[:], loadedKey[:]) {
		t.Fatal("loaded key does not match original key")
	}
}
