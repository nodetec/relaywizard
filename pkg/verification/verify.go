package verification

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/commands"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Function to verify relay binaries
func VerifyRelayBinary(currentUsername, relayName, path string) {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Verifying %s binary...", relayName))

	// Import NODE-TEC PGP key
	commands.PipeTwoCommands(exec.Command("curl", NodeTecKeybasePGPKeyURL), exec.Command("gpg", "--import"), "Failed to import NODE-TEC PGP key:")

	spinner.UpdateText("Imported NODE-TEC PGP key")

	// Determine the file name from the URL
	relaysManifestSigFile := filepath.Base(RelaysManifestSigFileURL)

	// Temporary file path
	relaysManifestSigFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, relaysManifestSigFile)

	// Check if the relay manifest signature file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(relaysManifestSigFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, relaysManifestSigFilePath)
	}

	// Download and copy the file
	files.DownloadAndCopyFile(currentUsername, relaysManifestSigFilePath, RelaysManifestSigFileURL, 0644)

	// Determine the file name from the URL
	relaysManifestFile := filepath.Base(RelaysManifestFileURL)

	// Temporary file path
	relaysManifestFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, relaysManifestFile)

	// Check if the temporary file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(relaysManifestFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, relaysManifestFilePath)
	}

	// Download and copy the file
	files.DownloadAndCopyFile(currentUsername, relaysManifestFilePath, RelaysManifestFileURL, 0644)

	// Use GPG to verify the manifest signature file
	out, err := exec.Command("gpg", "--status-fd", "1", "--verify", relaysManifestSigFilePath).Output()

	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to run the gpg verify command on the %s file: %v", relaysManifestSigFilePath, err)
		os.Exit(1)
	}

	validSig := strings.Contains(string(out), fmt.Sprintf("[GNUPG:] VALIDSIG %s", NodeTecSigningSubkeyFingerprint))

	if validSig {
		spinner.UpdateText(fmt.Sprintf("Verified the signature of the %s file", relaysManifestFilePath))
	} else {
		pterm.Println()
		pterm.Error.Printfln("Failed to verify the signature of the %s file", relaysManifestFilePath)
		os.Exit(1)
	}

	// Compute the SHA512 hash of the compressed relay binary file
	out, err = exec.Command("sha512sum", path).Output()

	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to compute the SHA512 hash of the %s file: %v", path, err)
		os.Exit(1)
	}

	// Extract the SHA512 hash from the output
	sha512Hash, _, _ := strings.Cut(string(out), fmt.Sprintf("  %s", path))

	// Read the manifest file
	data, err := os.ReadFile(relaysManifestFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to read the %s file: %v", relaysManifestFilePath, err)
		os.Exit(1)
	}

	// Search the manifest file for the hash
	if strings.Contains(string(data), sha512Hash) {
		spinner.UpdateText(fmt.Sprintf("Verified the SHA512 hash of the %s file", path))
		spinner.Success(fmt.Sprintf("%s binary verified", relayName))
	} else {
		pterm.Println()
		pterm.Error.Printfln("Failed to verify the %s file, the SHA512 hash doesn't match the SHA512 hash in the %s file", path, relaysManifestFilePath)
		os.Exit(1)
	}
}
