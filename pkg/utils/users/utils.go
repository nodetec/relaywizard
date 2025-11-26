package users

import (
	"os"
	"os/exec"
	"os/user"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/pterm/pterm"
)

func CheckCurrentUsername() string {
	currentUser, err := user.Current()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to get current user: %v", err)
		os.Exit(1)
	}

	return currentUser.Username
}

func UserExists(username string) bool {
	cmd := exec.Command("id", "-u", username)
	err := cmd.Run()
	return err == nil
}

func CreateUser(currentUsername, username string, disableLogin bool) {
	if disableLogin {
		if currentUsername == relays.RootUser {
			err := exec.Command("adduser", "--disabled-login", "--gecos", "", username).Run()
			if err != nil {
				pterm.Println()
				pterm.Error.Printfln("Failed to create user: %v", err)
				os.Exit(1)
			}
		} else {
			err := exec.Command("sudo", "adduser", "--disabled-login", "--gecos", "", username).Run()
			if err != nil {
				pterm.Println()
				pterm.Error.Printfln("Failed to create user: %v", err)
				os.Exit(1)
			}
		}
	}
}
