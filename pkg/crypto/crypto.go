package cryptolocal

import (
	connectionHandlers "GoNyx/pkg/connection_handlers"
	"GoNyx/pkg/global"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
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

func EncryptCommunication(secret []byte, data []byte) ([]byte, error) {

	key, err := hashSecretForAESKey(secret)
	if err != nil {
		return nil, errors.New("cannot hash secret in crypto.go")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherText := aesgcm.Seal(nil, nonce, data, nil)
	payload := append(nonce, cipherText...)

	connectionHandlers.SendEncryptedConnectionToRelay(payload, global.ListenIP, global.RelayPort) // todo move this

	return payload, nil
}

// hashSecretForAESKey generates a key size for AES-256 from the ecdh secret
func hashSecretForAESKey(secret []byte) ([]byte, error) {
	// might want to salt this in the future depending on ttl?
	salt := []byte(nil)
	hkdfReader := hkdf.New(sha256.New, secret, salt, nil)

	key := make([]byte, 32) // 32 bytes for AES-256
	if _, err := io.ReadFull(hkdfReader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// DecryptCommunication decrypts a message encrypted with ecdh + aes-256
func DecryptCommunication(cipherText []byte, secret []byte) ([]byte, error) {
	key, err := hashSecretForAESKey(secret)
	if err != nil {
		return nil, errors.New("cannot hash secret")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	// extract the nonce
	nonce := cipherText[:gcm.NonceSize()]
	data := cipherText[gcm.NonceSize():]

	// decrypt
	plainText, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
