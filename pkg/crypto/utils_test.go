package crypto

import (
	"fmt"
	"testing"
)

func TestEncryptionOfAString(t *testing.T) {
	keyAlice := GeneratePrivateKeyP521()
	keyBob := GeneratePrivateKeyP521()
	message := "this is a test message"

	messageEnc := Encrypt(keyAlice, keyBob.PublicKey(), []byte(message))
	fmt.Printf("🔐 Encrypted message: %x\n", messageEnc)

	//dest decrypt
	plainMessage := Decrypt(keyAlice.PublicKey(), keyBob, messageEnc)
	plainText := string(plainMessage)
	fmt.Printf("🔐 Dencrypted message: %v\n", plainText)

	if message != plainText {
		t.Fatalf("Expected %v got %v", message, plainText)
	}

}
