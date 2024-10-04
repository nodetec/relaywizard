package users

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

func UserExists(username string) bool {
	cmd := exec.Command("id", "-u", username)
	err := cmd.Run()
	return err == nil
}

func CreateUser(username string, disableLogin bool) {
	if disableLogin {
		err := exec.Command("adduser", "--disabled-login", "--gecos", "", username).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to create user: %v", err))
			os.Exit(1)
		}
	}
}
