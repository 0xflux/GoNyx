package clientcore

import (
	"GoNyx/pkg/crypto"
	"crypto/ecdsa"
)

type Relay struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func New() *Relay {
	private, pub := crypto.NewKeyPair() // gen new key pair for comms encryption

	return &Relay{
		PrivateKey: private,
		PublicKey:  pub,
	}
}
