package git

import (
	"io/fs"
	"os"
	"os/exec"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/pterm/pterm"
)

type FileMode = fs.FileMode

// Function to clone a repository
func Clone(branch, url, destDir string) {
	err := exec.Command("git", "clone", "-b", branch, url, destDir).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to download repository: %v", err)
		os.Exit(1)
	}
}

// TODO
// Just download the necessary files instead of the whole repo
func RemoveThenClone(currentUsername, repoDirPath, branch, URL string, permissions FileMode) {
	if currentUsername == relays.RootUser {
		// Check for and remove existing git repository
		directories.RemoveDirectory(repoDirPath)

		// Download git repository
		Clone(branch, URL, repoDirPath)

		directories.SetPermissions(repoDirPath, permissions)
	} else {
		// Check for and remove existing git repository
		directories.RemoveDirectoryUsingLinux(currentUsername, repoDirPath)

		// Download git repository
		Clone(branch, URL, repoDirPath)

		directories.SetPermissionsUsingLinux(currentUsername, repoDirPath, "0755")
	}
}
