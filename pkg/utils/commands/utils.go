package commands

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

func PipeTwoCommands(commandOne, commandTwo *exec.Cmd, errMsg string) {
	r, w, err := os.Pipe()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to create pipe: %v", err))
		os.Exit(1)
	}
	defer r.Close()
	commandOne.Stdout = w
	err = commandOne.Start()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("%s %v", errMsg, err))
		os.Exit(1)
	}
	defer commandOne.Wait()
	w.Close()
	commandTwo.Stdin = r
	commandTwo.Stdout = os.Stdout
	commandTwo.Run()
}