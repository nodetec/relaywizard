package main

import (
	"github.com/nodetec/rwz/cmd"
	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/ui"
	"github.com/nodetec/rwz/pkg/utils/users"
)

func init() {
	ui.Greet()

	// Check current username
	currentUsername := users.CheckCurrentUsername()

	// Set up sudo session
	users.SetUpSudoSession(currentUsername)

	// Set up rwz log file
	logs.SetUpRWZLogFile(currentUsername)
}

func main() {
	// Check current username
	currentUsername := users.CheckCurrentUsername()

	cmd.Execute(currentUsername)
}
