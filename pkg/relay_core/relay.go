package relay_core

import (
	cryptolocal "GoNyx/pkg/crypto"
	"GoNyx/pkg/global"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
)

/*
An 'OOP' style approach to managing relays self-contained settings.
*/

type Relay struct {
	PrivateKey    []byte `json:"privateKey"`
	PublicKey     []byte `json:"publicKey"`
	PublicKeyHash string `json:"publicKeyHash"`
}

type RelayOld struct {
	PrivateKey        JSONPrivKey `json:"privatekey"`
	PublicKey         JSONPubKey  `json:"publickey"`
	PublicKeyHash     string      `json:"publickeyhash"`
	DigitalPublicKey  *ecdsa.PublicKey
	DigitalPrivateKey *ecdsa.PrivateKey
}

type JSONPrivKey struct {
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
	D *big.Int `json:"d"`
}
type JSONPubKey struct {
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
}

func NewRelay() *Relay {
	// check if server is already initialised in previous session. if not, gen new public key & register server.
	locations, err := global.GetFileLocations()
	if err != nil {
		log.Fatalf("Error getting filepath, %s", err)
	}

	settingsData, err := os.ReadFile(filepath.Join(locations.Settings, global.RelaySettingsFileName))
	var relay *Relay

	if err != nil {
		if os.IsNotExist(err) {
			// relay not set up previously, so create as new relay
			// generate keys
			private, pub, err := cryptolocal.NewECDHKeyPair() // gen new key pair for long term fingerprinting
			if err != nil {
				log.Fatal("Error generating long term keys. Quitting. ", err)
			}
			pubHash := cryptolocal.Sha256Fingerprint(pub)

			relay = &Relay{
				PrivateKey:    private.Bytes(), // encode to json as bytes
				PublicKey:     pub.Bytes(),     // encode to json as bytes
				PublicKeyHash: pubHash,
			}

			jsonData, err := json.MarshalIndent(relay, "", "	")
			if err != nil {
				log.Fatalf("Error marshalling json, %s", err)
			}

			err = os.WriteFile(filepath.Join(locations.Settings, global.RelaySettingsFileName), jsonData, 0644)
			if err != nil {
				log.Fatalf("Error writing settings json, %s", err)
			}
			// TODO check in the json file with the directory server
			fmt.Println("New server settings created.")
		} else {
			log.Fatalf("Error reading file, %s", err)
		}
	} else {
		fmt.Println("Server already registered, starting.")
		relay = &Relay{}
		if err = json.Unmarshal(settingsData, relay); err != nil {
			log.Fatalf("Error reading settings, %s", err)
		}
	}

	fmt.Println("My public key is: ", relay.PublicKey)

	return relay
}
