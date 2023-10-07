package global

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

/*
Global constants that will be used through the project, these may be
configuration settings, hardcoded paths, etc.
*/

type ProgramLocations struct {
	Logs     string
	Settings string
}

const (
	ClientPort            = 34888 // may want to change to var in future so people can change it on the fly?
	RelayPort             = 34889 // may want to change to var in future so people can change it on the fly?
	Version               = "0.0.1"
	ListenIP              = "127.0.0.1" // IP to listen on locally
	RelaySettingsFileName = "relay_settings.json"
)

// anything we cannot init as constant, or do not require as constants.
var (
	ClientListenAddr = fmt.Sprintf("localhost:%d", ClientPort)
)

func GetFileLocations() (*ProgramLocations, error) {
	switch runtime.GOOS {
	case "windows":
		// for windows put all data into AppData
		appdata, exists := os.LookupEnv("APPDATA")
		if !exists {
			return &ProgramLocations{}, errors.New("cannot find appdata")
		}
		target := filepath.Join(appdata, "gonyx")
		if err := os.MkdirAll(target, 0755); err != nil {
			return &ProgramLocations{}, errors.New("cannot create gonyx folder")
		}
		return &ProgramLocations{Logs: target, Settings: target}, nil

	case "linux":
		locs := &ProgramLocations{Logs: "/var/log/gonyx/", Settings: "/etc/gonyx/"}
		return locs, nil

	default:
		return &ProgramLocations{}, errors.New(fmt.Sprintf("running on unknown os, %v", runtime.GOOS))
	}
}
