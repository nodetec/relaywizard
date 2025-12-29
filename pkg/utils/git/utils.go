package git

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/pterm/pterm"
)

type FileMode = fs.FileMode

// Function to clone a repository
func Clone(currentUsername, branch, url, destDir string) {
	err := exec.Command("git", "clone", "-b", branch, url, destDir).Run()
	if err != nil {
		logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to download repository: %v", err))
		pterm.Println()
		pterm.Error.Printfln("Failed to download repository: %v", err)
		os.Exit(1)
	}
}

// TODO
// Just download the necessary files instead of the whole repo
func RemoveThenClone(currentUsername, repoDirPath, branch, URL, permissionsAsString string, permissions FileMode) {
	if currentUsername == relays.RootUser {
		// Check for and remove existing git repository
		directories.RemoveDirectory(repoDirPath)

		// Download git repository
		Clone(currentUsername, branch, URL, repoDirPath)

		directories.SetPermissions(repoDirPath, permissions)
	} else {
		// Check for and remove existing git repository
		directories.RemoveDirectoryUsingLinux(currentUsername, repoDirPath)

		// Download git repository
		Clone(currentUsername, branch, URL, repoDirPath)

		directories.SetPermissionsUsingLinux(currentUsername, repoDirPath, permissionsAsString)
	}
}
