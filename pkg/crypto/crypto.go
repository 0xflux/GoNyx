package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

// NewDHKeyPair  generate a key pair for use in Diffie-Hellman exchanges
func NewDHKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatal("Fatal error generating private key. ", err)
	}

	pubKey := &privateKey.PublicKey

	return privateKey, pubKey
}

// Sha256Fingerprint generates a hash of a public key. Do not use this for hashing private keys.
func Sha256Fingerprint(publicKey *ecdsa.PublicKey) string {
	// combine the x and y points
	x := publicKey.X.Bytes()
	y := publicKey.Y.Bytes()
	combined := append(x, y...)

	// hash and return
	hash := sha256.Sum256(combined)
	return hex.EncodeToString(hash[:])
}
