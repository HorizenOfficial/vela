package crypto

import (
	"bytes"
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

func TestGeneratePrivateKey25519(t *testing.T) {
	key, err := GeneratePrivateKey25519()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	if key == nil {
		t.Fatal("key is nil")
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
