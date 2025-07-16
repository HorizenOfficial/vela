package crypto

import (
	"bytes"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestGeneratePrivateKeySecp256k1(t *testing.T) {
	key, err := GeneratePrivateKeySecp256k1()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	if key == nil {
		t.Fatal("key is nil")
	}
}

func TestSaveLoadPrivateKeySecp256k1DER(t *testing.T) {
	key, err := GeneratePrivateKeySecp256k1()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	filename := "test_key_secp256k1.der"
	defer os.Remove(filename)

	err = SavePrivateKeySecp256k1ToFileDER(key, filename)
	if err != nil {
		t.Fatalf("failed to save private key: %v", err)
	}

	loadedKey, err := LoadPrivateKeySecp256k1FromFileDER(filename)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	if !bytes.Equal(crypto.FromECDSA(key.PrivateKey), crypto.FromECDSA(loadedKey.PrivateKey)) {
		t.Fatal("loaded key does not match original key")
	}
}

func TestSaveLoadPrivateKeySecp256k1PEM(t *testing.T) {
	key, err := GeneratePrivateKeySecp256k1()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	filename := "test_key_secp256k1.pem"
	defer os.Remove(filename)

	err = SavePrivateKeySecp256k1ToFilePEM(key, filename)
	if err != nil {
		t.Fatalf("failed to save private key: %v", err)
	}

	loadedKey, err := LoadPrivateKeySecp256k1FromFilePEM(filename)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	if !bytes.Equal(crypto.FromECDSA(key.PrivateKey), crypto.FromECDSA(loadedKey.PrivateKey)) {
		t.Fatal("loaded key does not match original key")
	}
}

func TestExportImportPrivateKeySecp256k1Hex(t *testing.T) {
	key, err := GeneratePrivateKeySecp256k1()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	hexKey := ExportPrivateKeySecp256k1ToHex(key)

	loadedKey, err := ImportPrivateKeySecp256k1FromHex(hexKey)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	if !bytes.Equal(crypto.FromECDSA(key.PrivateKey), crypto.FromECDSA(loadedKey.PrivateKey)) {
		t.Fatal("loaded key does not match original key")
	}
}
