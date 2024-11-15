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
func VerifyRelayBinary(path string) {
	spinner, _ := pterm.DefaultSpinner.Start("Verifying relay binary...")
	pterm.Println()

	// Import NODE-TEC PGP key
	commands.PipeTwoCommands(exec.Command("curl", NodeTecKeybasePGPKeyURL), exec.Command("gpg", "--import"), "Failed to import NODE-TEC PGP key:")

	spinner.UpdateText("Imported NODE-TEC PGP key")

	// Determine the file name from the URL
	relaysManifestSigFile := filepath.Base(RelaysManifestSigFileURL)

	// Temporary file path
	relaysManifestSigFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, relaysManifestSigFile)

	// Check if the relay manifest signature file exists and remove it if it does
	files.RemoveFile(relaysManifestSigFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(relaysManifestSigFilePath, RelaysManifestSigFileURL)

	// Determine the file name from the URL
	relaysManifestFile := filepath.Base(RelaysManifestFileURL)

	// Temporary file path
	relaysManifestFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, relaysManifestFile)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(relaysManifestFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(relaysManifestFilePath, RelaysManifestFileURL)

	// Use GPG to verify the manifest signature file and output the primary key and signature subkey fingerprints
	cmd := exec.Command("gpg", "--verify", "--with-fingerprint", "--with-subkey-fingerprints", relaysManifestSigFilePath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to run the gpg verify command on the %s file: %v", relaysManifestSigFilePath, err))
		os.Exit(1)
	}

	gpgVerifyOutput := string(out)

	goodSig := strings.Contains(gpgVerifyOutput, NodeTecGoodSigMsg)

	// Extract the formatted primary key and formatted signature subkey fingerprints from the output
	_, formattedPrimaryAndSubKeyFingerprints, _ := strings.Cut(gpgVerifyOutput, "Primary key fingerprint: ")

	formattedPrimaryKeyFingerprint, formattedSubkeyFingerprint, _ := strings.Cut(formattedPrimaryAndSubKeyFingerprints, "Subkey fingerprint: ")

	// Remove the spaces and new line characters from the formatted primary key and formatted signature subkey fingerprints
	formattedPrimaryKeyFingerprint = strings.ReplaceAll(formattedPrimaryKeyFingerprint, " ", "")
	formattedSubkeyFingerprint = strings.ReplaceAll(formattedSubkeyFingerprint, " ", "")

	primaryKeyFingerprint := strings.ReplaceAll(formattedPrimaryKeyFingerprint, "\n", "")
	subkeyFingerprint := strings.ReplaceAll(formattedSubkeyFingerprint, "\n", "")

	if goodSig && primaryKeyFingerprint == NodeTecPrimaryKeyFingerprint && subkeyFingerprint == NodeTecSigningSubkeyFingerprint {
		spinner.UpdateText(fmt.Sprintf("Verified the signature of the %s file and the fingerprints", relaysManifestFilePath))
	} else {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to verify the signature of the %s file and/or the fingerprints", relaysManifestFilePath))
		os.Exit(1)
	}

	// Compute the SHA512 hash of the compressed relay binary file
	out, err = exec.Command("sha512sum", path).Output()

	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to compute the SHA512 hash of the %s file: %v", path, err))
		os.Exit(1)
	}

	// Extract the SHA512 hash from the output
	sha512Hash, _, _ := strings.Cut(string(out), fmt.Sprintf("  %s", path))

	// Read the manifest file
	data, err := os.ReadFile(relaysManifestFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to read the %s file: %v", relaysManifestFilePath, err))
		os.Exit(1)
	}

	// Search the manifest file for the hash
	if strings.Contains(string(data), sha512Hash) {
		spinner.UpdateText(fmt.Sprintf("Verified the SHA512 hash of the %s file", path))
		spinner.Success("Relay binary verified")
	} else {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to verify the %s file, the SHA512 hash doesn't match the SHA512 hash in the %s file", path, relaysManifestFilePath))
		os.Exit(1)
	}
}
