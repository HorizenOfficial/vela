package main

import (
	"log"

	"github.com/horizen-pes/pkg/crypto"
)

func main() {
	log.Println("Generating random key pairs..")
	keyAES, _ := crypto.GenerateAESKey()
	strAES := crypto.ExportKeyAESToHex(&keyAES)
	log.Println("AES: 0x" + strAES)

	key25519, _ := crypto.GeneratePrivateKey25519()
	str25519 := crypto.ExportPrivateKey25519ToHex(key25519)
	log.Println("25519 (private): 0x" + str25519)
	str25519pub := crypto.ExportPublicKey25519ToHex(key25519.PublicKey())
	log.Println("25519 (public): 0x" + str25519pub)

	keySecp256k1, _ := crypto.GeneratePrivateKeySecp256k1()
	strKeySecp256k1 := crypto.ExportPrivateKeySecp256k1ToHex(keySecp256k1)
	log.Println("Secp256k1 (private): 0x" + strKeySecp256k1)
	strKeySecp256k1Pub := crypto.ExportPublicKeySecp256k1ToHex(keySecp256k1.PublicKey())
	log.Println("Secp256k1 (public): 0x" + strKeySecp256k1Pub)
	log.Println("Secp256k1 (Ethereum address): " + keySecp256k1.PublicKey().Address())

	keyP521, _ := crypto.GeneratePrivateKeyP521()
	keyP521Str := crypto.ExportPrivateKeyP521ToHex(keyP521)
	log.Println("P521 (private): 0x" + keyP521Str)
	keyP521StrPub := crypto.ExportPublicKeyP521ToHex(keyP521.PublicKey())
	log.Println("P521 (public): 0x" + keyP521StrPub)

}
