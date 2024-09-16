package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/manager"
	"github.com/pterm/pterm"
	"os/exec"
)

// Function to install necessary strfry package dependencies
func AptInstallDependencies() {
	spinner, _ := pterm.DefaultSpinner.Start("Installing strfry dependencies...")

	packages := []string{"git", "build-essential", "libyaml-perl", "libtemplate-perl", "libregexp-grammars-perl", "libssl-dev", "zlib1g-dev", "liblmdb-dev", "libflatbuffers-dev", "libsecp256k1-dev", "libzstd-dev"}

	// Check if package is installed, install if not
	for _, p := range packages {
		if manager.IsPackageInstalled(p) {
			spinner.UpdateText(fmt.Sprintf("%s is already installed.", p))
		} else {
			spinner.UpdateText(fmt.Sprintf("Installing %s...", p))
			exec.Command("apt", "install", "-y", "-qq", p).Run()
		}
	}

	spinner.Success("strfry dependencies installed successfully.")
}
