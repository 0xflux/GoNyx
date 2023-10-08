package cryptolocal

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

// NewECDHKeyPair generate a key pair for use in Diffie-Hellman
func NewECDHKeyPair() (*ecdh.PrivateKey, *ecdh.PublicKey) {

	privateKey, err := ecdh.P521().GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Error making privateKey, %s", err)
	}

	publicKey := privateKey.PublicKey()

	return privateKey, publicKey
}

// Sha256Fingerprint generates a hash of a public key. Do not use this for hashing private keys.
func Sha256Fingerprint(publicKey *ecdh.PublicKey) string {
	// hash and return
	hash := sha256.Sum256(publicKey.Bytes())
	return hex.EncodeToString(hash[:])
}
