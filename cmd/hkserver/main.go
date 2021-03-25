package main

import (
	"fmt"
	"os"

	"github.com/masterhung0112/hk_server/v5/cmd/hkserver/commands"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
