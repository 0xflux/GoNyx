package main

import (
	clientcore "GoNyx/pkg/client_core"
	"GoNyx/pkg/global"
	"fmt"
)

func main() {
	fmt.Println("Nyx Client starting, version: ", global.Version)

	clientcore.StartClient() // start the client
}
