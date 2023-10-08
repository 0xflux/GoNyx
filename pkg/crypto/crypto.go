package cryptolocal

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

// DEPRECATED FOR NOW, REPLACED WITH NewECDHKeyPair
// NewDHKeyPair  generate a key pair for use in Diffie-Hellman exchanges
//func NewDHKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
//	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
//	if err != nil {
//		log.Fatal("Fatal error generating private key. ", err)
//	}
//
//	return privateKey, &privateKey.PublicKey
//}

// NewECDHKeyPair generate a key pair for use in Diffie-Hellman exchanges
func NewECDHKeyPair() (*ecdh.PrivateKey, *ecdh.PublicKey) {

	privateKey, err := ecdh.P521().GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Error making privateKey, %s", err)
	}
	fmt.Printf("Private key generated is: %v\n", privateKey)

	publicKey := privateKey.PublicKey()

	fmt.Printf("Public Key generated is: %v\n", publicKey)

	return privateKey, publicKey
}

// Sha256Fingerprint generates a hash of a public key. Do not use this for hashing private keys.
func Sha256Fingerprint(publicKey *ecdh.PublicKey) string {
	// hash and return
	hash := sha256.Sum256(publicKey.Bytes())
	return hex.EncodeToString(hash[:])
}

// power calculates (base^exp) mod modValue for the D-H secret
//func power(base, exp, mod *big.Int) *big.Int {
//	return new(big.Int).Exp(base, exp, mod)
//}
//
//func ComputeSharedSecret(externPubKey *ecdh.PublicKey, privateKey relay_core.JSONPrivKey) *big.Int {
//	p := ecdh.PrivateKey{}
//	secret, err := p.ECDH(externPubKey)
//	if err != nil {
//		log.Fatalf("Error getting ecdh, %s", err)
//	}
//	fmt.Println("secret: ", secret)
//
//	return 39784632846837643423434
//}
//
//// ParseJsonToKey parsing a json object to a key
//func ParseJsonToKey(jsonPrivate string, jsonPublic string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
//	var jKey relay_core.JSONPrivKey
//	err := json.Unmarshal([]byte(jsonPrivate), &jKey)
//	if err != nil {
//		log.Fatalf("error unmarshalling json to private key %v", err)
//	}
//
//	// Construct the private key
//	privateKey := &ecdsa.PrivateKey{
//		PublicKey: ecdsa.PublicKey{
//			Curve: elliptic.P521(),
//			X:     jKey.X,
//			Y:     jKey.Y,
//		},
//		D: jKey.D,
//	}
//
//	err = json.Unmarshal([]byte(jsonPublic), &jKey)
//	if err != nil {
//		log.Fatalf("error unmarshalling json to public key %v", err)
//	}
//
//	// Construct the private key
//	publicKey := &ecdsa.PublicKey{
//		Curve: elliptic.P521(),
//		X:     jKey.X,
//		Y:     jKey.Y,
//	}
//
//	return privateKey, publicKey
//}
