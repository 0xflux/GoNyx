package relay_core

import (
	"GoNyx/pkg/crypto"
	"GoNyx/pkg/global"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

/*
An 'OOP' style approach to managing relays self-contained settings.
*/

type Relay struct {
	PrivateKeyFingerprint *ecdsa.PrivateKey `json:"privateKeyFingerprint"` // not sure if we need to save this private key
	PublicKeyFingerprint  *ecdsa.PublicKey  `json:"publicKeyFingerprint"`
	PublicKeyHash         string            `json:"PublicKeyHash"`
}

func NewRelay() *Relay {
	private, pub := crypto.NewDHKeyPair() // gen new key pair for long term fingerprinting
	pubHash := crypto.Sha256Fingerprint(pub)

	relay := &Relay{
		PrivateKeyFingerprint: private,
		PublicKeyFingerprint:  pub,
		PublicKeyHash:         pubHash,
	}

	// check if server is already initialised in previous session. if not, gen new public key & register server.
	locations, err := global.GetFileLocations()
	if err != nil {
		log.Fatalf("Error getting filepath, %s", err)
	}
	_, err = os.ReadFile(filepath.Join(locations.Settings, global.RelaySettingsFileName))
	if err != nil {
		if os.IsNotExist(err) {
			// relay not set up previously, so create as new relay
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
	}

	return relay
}
