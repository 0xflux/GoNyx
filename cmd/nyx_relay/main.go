package main

import (
	"GoNyx/pkg/global"
	"GoNyx/pkg/relay_core"
	"fmt"
)

func main() {
	fmt.Println("Nyx Relay starting, version no ", global.Version)

	relay_core.StartRelay()
}
