package relay_core

import (
	"GoNyx/pkg/crypto"
	"GoNyx/pkg/global"
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
	PrivateKey    JSONPrivKey `json:"privatekey"`
	PublicKey     JSONPubKey  `json:"publickey"`
	PublicKeyHash string      `json:"publickeyhash"`
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
			private, pub := crypto.NewDHKeyPair() // gen new key pair for long term fingerprinting
			pubHash := crypto.Sha256Fingerprint(pub)

			relay = &Relay{
				PrivateKey:    JSONPrivKey{X: private.X, Y: private.Y, D: private.D},
				PublicKey:     JSONPubKey{X: pub.X, Y: pub.Y},
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

	return relay
}
