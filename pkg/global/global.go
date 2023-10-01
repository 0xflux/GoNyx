package global

import "fmt"

/*
Global constants that will be used through the project, these may be
configuration settings, hardcoded paths, etc.
*/

const (
	ClientPort = 34888 // may want to change to var in future so people can change it on the fly?
	RelayPort  = 34889 // may want to change to var in future so people can change it on the fly?
	Version    = "0.0.1"
	ListenIP   = "127.0.0.1" // IP to listen on locally
)

// anything we cannot init as constant, or do not require as constants.
var (
	ClientListenAddr = fmt.Sprintf("localhost:%d", ClientPort)
)
