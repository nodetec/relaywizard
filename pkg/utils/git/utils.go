package git

import (
	"fmt"
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
