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

func TestEncryptDecrypt(t *testing.T) {
	keyAlice, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate Alice's private key: %v", err)
	}
	keyBob, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate Bob's private key: %v", err)
	}

	message := []byte("this is a test message")

	// Encrypt
	encryptedMessage, err := Encrypt(keyAlice, keyBob.PublicKey(), message)
	if err != nil {
		t.Fatalf("failed to encrypt message: %v", err)
	}

	// Decrypt
	decryptedMessage, err := Decrypt(keyAlice.PublicKey(), keyBob, encryptedMessage)
	if err != nil {
		t.Fatalf("failed to decrypt message: %v", err)
	}

	if !bytes.Equal(message, decryptedMessage) {
		t.Fatalf("expected %s, got %s", message, decryptedMessage)
	}
}

func TestEncryptDecryptWithAES(t *testing.T) {
	// Generate a random AES key
	aesKey, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("failed to generate random AES key: %v", err)
	}

	message := []byte("this is a test message")

	// Encrypt
	encryptedMessage, err := EncryptWithAES(aesKey, message)
	if err != nil {
		t.Fatalf("failed to encrypt message: %v", err)
	}

	// Decrypt
	decryptedMessage, err := DecryptWithAES(aesKey, encryptedMessage)
	if err != nil {
		t.Fatalf("failed to decrypt message: %v", err)
	}

	if !bytes.Equal(message, decryptedMessage) {
		t.Fatalf("expected %s, got %s", message, decryptedMessage)
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	keyAlice, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate Alice's private key: %v", err)
	}
	keyBob, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate Bob's private key: %v", err)
	}
	keyEve, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate Eve's private key: %v", err)
	}

	message := []byte("this is a test message")

	// Encrypt
	encryptedMessage, err := Encrypt(keyAlice, keyBob.PublicKey(), message)
	if err != nil {
		t.Fatalf("failed to encrypt message: %v", err)
	}

	// Decrypt with wrong key
	_, err = Decrypt(keyAlice.PublicKey(), keyEve, encryptedMessage)
	if err == nil {
		t.Fatal("expected decryption to fail, but it succeeded")
	}
}

func TestDecryptTamperedMessage(t *testing.T) {
	keyAlice, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate Alice's private key: %v", err)
	}
	keyBob, err := GeneratePrivateKeyP521()
	if err != nil {
		t.Fatalf("failed to generate Bob's private key: %v", err)
	}

	message := []byte("this is a test message")

	// Encrypt
	encryptedMessage, err := Encrypt(keyAlice, keyBob.PublicKey(), message)
	if err != nil {
		t.Fatalf("failed to encrypt message: %v", err)
	}

	// Tamper with the message
	encryptedMessage[len(encryptedMessage)-1] ^= 0xff

	// Decrypt tampered message
	_, err = Decrypt(keyAlice.PublicKey(), keyBob, encryptedMessage)
	if err == nil {
		t.Fatal("expected decryption to fail, but it succeeded")
	}
}
