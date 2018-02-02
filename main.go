package main

import (
	"fmt"
	"os"

	"cadvisor-cli/commands"
)

//
func main() {
	if err := commands.Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
