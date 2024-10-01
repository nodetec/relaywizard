package users

import (
	"log"
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
			log.Fatalf("Error creating user: %v", err)
		}
	}
}
