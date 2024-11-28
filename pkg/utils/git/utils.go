package git

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
	"os/exec"
)

type FileMode = fs.FileMode

// Function to clone a repository
func Clone(branch, url, destDir string) {
	err := exec.Command("git", "clone", "-b", branch, url, destDir).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to download repository: %v", err))
		os.Exit(1)
	}
}

// TODO
// Just download the necessary files instead of the whole repo
func RemoveThenClone(repoDirPath, branch, URL string, permissions FileMode) {
	// Check for and remove existing git repository
	directories.RemoveDirectory(repoDirPath)

	// Download git repository
	Clone(branch, URL, repoDirPath)

	directories.SetPermissions(repoDirPath, permissions)
}
