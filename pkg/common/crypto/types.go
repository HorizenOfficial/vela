package crypto

import (
	"crypto/ecdh"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

// PublicKeyP521 is a public key with Elliptic Curve Diffie-Hellman over NIST P-521 curve, also known as secp521r1.
// (This can be used for encryption/decription)
type PublicKeyP521 struct {
	*ecdh.PublicKey
}

// Bytes returns the byte representation of the public key.
func (p *PublicKeyP521) Bytes() []byte {
	if p == nil || p.PublicKey == nil {
		return nil
	}
	return p.PublicKey.Bytes()
}

// NewPublicKeyP521 creates a PublicKeyP521 from a byte slice.
func NewPublicKeyP521(b []byte) (*PublicKeyP521, error) {
	curve := ecdh.P521()
	pub, err := curve.NewPublicKey(b)
	if err != nil {
		return nil, err
	}
	return &PublicKeyP521{pub}, nil
}

// PrivateKeyP521 is a private key with Elliptic Curve Diffie-Hellman over NIST P-521 curve, also known as secp521r1.
// (This can be used for encryption/decription)
type PrivateKeyP521 struct {
	*ecdh.PrivateKey
}

// PublicKey returns the public key associated with a private key.
func (p *PrivateKeyP521) PublicKey() *PublicKeyP521 {
	return &PublicKeyP521{p.PrivateKey.PublicKey()}
}

// PublicKey25519 is a public key with Elliptic Curve 25519.
type PublicKey25519 struct {
	*ecdh.PublicKey
}

// PrivateKey25519 is a private key with with Elliptic Curve 25519.
type PrivateKey25519 struct {
	*ecdh.PrivateKey
}

// PublicKey returns the public key associated with a private key.
func (p *PrivateKey25519) PublicKey() *PublicKey25519 {
	return &PublicKey25519{p.PrivateKey.PublicKey()}
}

// PublicKeySecp256k1 is a public key with Elliptic Curve secp256k1.
// (This is the curve used in Bitcoin and Ethereum)
type PublicKeySecp256k1 struct {
	*ecdsa.PublicKey
}

// Address returns the Ethereum address of the public key.
func (p *PublicKeySecp256k1) Address() string {
	return crypto.PubkeyToAddress(*p.PublicKey).Hex()
}

// PrivateKeySecp256k1 is a private key with Elliptic Curve secp256k1.
// (This is the curve used in Bitcoin and Ethereum)
type PrivateKeySecp256k1 struct {
	*ecdsa.PrivateKey
}

// PublicKey returns the public key associated with a private key.
func (p *PrivateKeySecp256k1) PublicKey() *PublicKeySecp256k1 {
	return &PublicKeySecp256k1{&p.PrivateKey.PublicKey}
}

// Sign calculates an Ethereum signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func (p *PrivateKeySecp256k1) Sign(digest []byte) ([]byte, error) {
	return crypto.Sign(digest, p.PrivateKey)
}

// Represents a AES-256-GCM key
type AES256Key [32]byte
