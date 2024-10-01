package git

import (
	"io/fs"
	"log"
	"os/exec"
)

type FileMode = fs.FileMode

// Function to clone a repository
func Clone(branch, url, destDir string) {
	err := exec.Command("git", "clone", "-b", branch, url, destDir).Run()
	if err != nil {
		log.Fatalf("Error downloading repository: %v", err)
	}
}
