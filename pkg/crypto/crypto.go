package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

// NewKeyPair generate a key pair
func NewKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal("Fatal error generating private key. ", err)
	}

	pubKey := &privateKey.PublicKey

	return privateKey, pubKey
}

//func deriveSharedSecret(selfPrivateKey *ecdsa.PrivateKey, remotePublicKey *ecdsa.PublicKey) []byte {
//	secret, _ := ecdh.Com
//}
