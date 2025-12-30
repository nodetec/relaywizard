package users

import (
	"errors"
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
	_, err := user.Lookup(username)
	if err != nil {
		var notFoundErr user.UnknownUserError
		if errors.As(err, &notFoundErr) {
			return false
		}

		pterm.Println()
		pterm.Error.Printfln("Failed to check if user exists: %v", err)
		os.Exit(1)
	}

	return true
}

func CreateUser(currentUsername, username string) {
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

func SetUpSudoSession(currentUsername string) {
	if currentUsername != relays.RootUser {
		err := exec.Command("sudo", "-v").Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to get password to set up sudo session: %v", err)
			os.Exit(1)
		}
		// TODO
		// Double check this command
		// What happens if a user's sudo session expires before 30 seconds, i.e., before the session can be extended by this loop?
		err = exec.Command("/bin/sh", "-c", "while true; do sudo -v; sleep 30 kill -0 $$ 2>/dev/null || exit; done &").Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set up sudo session: %v", err)
			os.Exit(1)
		}
	}
}
