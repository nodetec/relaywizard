package strfry

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pterm/pterm"
)

// Function to download, build, and install the binary
func InstallRelayBinary() {
	// TODO
	// Create the binary on a different machine then download it instead of building it here
	// Use these variables and model it after Khatru Pyramid installation

	// Temporary directory for git repository
	const tempDir = "/tmp/strfry"

	// URL of the binary to download
	// const downloadURL = "https://..."

	// Name of the binary after downloading
	const binaryName = "nostr-relay-strfry"

	// Destination directory for the binary
	const destDir = "/usr/local/bin"

	// Data directory for the relay
	// const dataDir = "/var/lib/nostr-relay-strfry"

	spinner, _ := pterm.DefaultSpinner.Start("Installing strfry relay...")

	pterm.Println()
	pterm.Println(pterm.Magenta("Go get coffee, this may take a few minutes..."))
	pterm.Println()

	// Download
	// Check for and remove existing repository
	err := os.RemoveAll(fmt.Sprintf("%s", tempDir))
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing repository: %v", err)
	}

	// Download git repository
	err = exec.Command("git", "clone", "https://github.com/hoytech/strfry.git", fmt.Sprintf("%s", tempDir)).Run()
	if err != nil {
		log.Fatalf("Error downloading repository: %v", err)
	}

	// Build
	// TODO
	// Check for development environment variable instead of commenting and uncommenting these lines

	// Check for and remove existing binary
	// When developing comment to prevent unecessary builds
	// err = os.Remove(fmt.Sprintf("%s/%s", destDir, binaryName))
	// if err != nil && !os.IsNotExist(err) {
	// 	log.Fatalf("Error removing existing binary: %v", err)
	// }

	// Check if binary exists
	// When developing uncomment to prevent unecessary builds
	_, err = os.Stat(fmt.Sprintf("%s/%s", destDir, binaryName))
	if os.IsNotExist(err) {
		// Intialize and update git submodule
		err = exec.Command("git", "-C", fmt.Sprintf("%s", tempDir), "submodule", "update", "--init").Run()
		if err != nil {
			log.Fatalf("Error initializing and updating git submodule: %v", err)
		}

		// Make setup-golpe
		err = exec.Command("make", "-C", fmt.Sprintf("%s", tempDir), "setup-golpe").Run()
		if err != nil {
			log.Fatalf("Error making setup-golpe: %v", err)
		}

		// Make -j2
		err = exec.Command("make", "-C", fmt.Sprintf("%s", tempDir), "-j2").Run()
		if err != nil {
			log.Fatalf("Error making -j2: %v", err)
		}

		// Install
		err = exec.Command("mv", fmt.Sprintf("%s/strfry", tempDir), fmt.Sprintf("%s/%s", destDir, binaryName)).Run()
		if err != nil {
			log.Fatalf("Error installing binary: %v", err)
		}
	}

	spinner.Success("strfry relay installed successfully.")
}
