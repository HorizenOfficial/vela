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

func TestExportImportPublicKeySecp256k1Hex(t *testing.T) {
	key, err := GeneratePrivateKeySecp256k1()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	pubKey := key.PublicKey()

	hexKey := ExportPublicKeySecp256k1ToHex(pubKey)

	loadedKey, err := ImportPublicKeySecp256k1FromHex(hexKey)
	if err != nil {
		t.Fatalf("failed to load public key: %v", err)
	}

	if !bytes.Equal(crypto.FromECDSAPub(pubKey.PublicKey), crypto.FromECDSAPub(loadedKey.PublicKey)) {
		t.Fatal("loaded key does not match original key")
	}
}


func TestPublicKeySecp256k1Address(t *testing.T) {
	dummykey := "4504532170fc47ec70aaae37b9a69a0fc8efda59dbaf36b2545651aa51151b97"
	dummyAddress := "0xDe98c1BCf65d928514C0a6e800b499054328314e" //ethereum address generated from the previous key

	myKey, _ := ImportPrivateKeySecp256k1FromHex(dummykey)
	myAddress := myKey.PublicKey().Address()

	if dummyAddress != myAddress {
		t.Fatal("Address not matching")
	}

}

func TestSignEthereum(t *testing.T) {
	key, err := GeneratePrivateKeySecp256k1()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	// Create a dummy message and hash it
	msg := []byte("test message")
	hash := crypto.Keccak256(msg)

	// Sign the hash
	sig, err := key.Sign(hash)
	if err != nil {
		t.Fatalf("failed to sign message: %v", err)
	}

	// The signature should be 65 bytes long (R + S + V)
	if len(sig) != 65 {
		t.Fatalf("signature has wrong length: got %d, want 65", len(sig))
	}

	// Recover the public key from the signature
	recoveredPubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		t.Fatalf("failed to recover public key: %v", err)
	}

	// The recovered public key should match the original public key
	if !bytes.Equal(crypto.FromECDSAPub(key.PublicKey().PublicKey), crypto.FromECDSAPub(recoveredPubKey)) {
		t.Fatal("recovered public key does not match original key")
	}
}

