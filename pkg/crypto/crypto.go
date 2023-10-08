package cryptolocal

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

// NewECDHKeyPair generate a key pair for use in Diffie-Hellman
func NewECDHKeyPair() (*ecdh.PrivateKey, *ecdh.PublicKey, error) {

	privateKey, err := ecdh.P521().GenerateKey(rand.Reader)
	if err != nil {
		log.Printf("Error making privateKey, %s\n", err)
		return nil, nil, err
	}

	publicKey := privateKey.PublicKey()

	return privateKey, publicKey, nil
}

// Sha256Fingerprint generates a hash of a public key. Do not use this for hashing private keys.
func Sha256Fingerprint(publicKey *ecdh.PublicKey) string {
	// hash and return
	hash := sha256.Sum256(publicKey.Bytes())
	return hex.EncodeToString(hash[:])
}

// ComputeSharedSecret computes the shared secret in a Diffie-Hellman exchange and returns the secret
func ComputeSharedSecret(externPublicKey *ecdh.PublicKey, privateKey *ecdh.PrivateKey) ([]byte, error) {
	secret, err := privateKey.ECDH(externPublicKey)
	if err != nil {
		fmt.Println("Error generating shared secret. ", err)
		return nil, err
	}

	fmt.Println("Secret is: ", secret)

	return secret, nil
}
