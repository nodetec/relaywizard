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
	ThemeDefault := pterm.ThemeDefault

	prompt := pterm.InteractiveContinuePrinter{
		DefaultValueIndex: 0,
		DefaultText:       "Do you want to continue with the installation?",
		TextStyle:         &ThemeDefault.PrimaryStyle,
		Options:           []string{"no", "yes"},
		OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
		SuffixStyle:       &ThemeDefault.SecondaryStyle,
		Delimiter:         ": ",
	}

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

	cmd := exec.Command("gpg", "--verify", relaysManifestSigFilePath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to run the gpg verify command on the %s file: %v", relaysManifestSigFilePath, err))
		os.Exit(1)
	}

	gpgVerifyOutput := string(out)

	goodSig := strings.Contains(gpgVerifyOutput, NodeTecGoodSigMsg)

	// Extract the formatted primary key and subkey fingerprints from the output
	_, formattedPrimaryAndSubKeyFingerprints, foundPrimaryKeyText := strings.Cut(gpgVerifyOutput, "Primary key fingerprint: ")

	formattedPrimaryKeyFingerprint, formattedSubkeyFingerprint, foundSubkeyText := strings.Cut(formattedPrimaryAndSubKeyFingerprints, "Subkey fingerprint: ")

	if foundPrimaryKeyText && foundSubkeyText {
		formattedPrimaryKeyFingerprint = strings.ReplaceAll(formattedPrimaryKeyFingerprint, " ", "")
		formattedSubkeyFingerprint = strings.ReplaceAll(formattedSubkeyFingerprint, " ", "")

		primaryKeyFingerprint := strings.ReplaceAll(formattedPrimaryKeyFingerprint, "\n", "")
		subkeyFingerprint := strings.ReplaceAll(formattedSubkeyFingerprint, "\n", "")

		if goodSig && primaryKeyFingerprint == NodeTecPrimaryKeyFingerprint && subkeyFingerprint == NodeTecSigningSubkeyFingerprint {
			spinner.UpdateText(fmt.Sprintf("Verified the signature of the %s file and the fingerprints", relaysManifestFilePath))
		} else {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to verify the signature of the %s file", relaysManifestFilePath))
			os.Exit(1)
		}
	} else {
		if goodSig {
			spinner.UpdateText(fmt.Sprintf("Verified the signature of the %s file", relaysManifestFilePath))
		} else {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to verify the signature of the %s file", relaysManifestFilePath))
			os.Exit(1)
		}
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
		pterm.Println()

		// Prompt user if they want to continue with installation without verifying fingerprints
		if !foundPrimaryKeyText || !foundSubkeyText {
			pterm.Println()
			spinner.Warning(fmt.Sprintf("Warning: The signature of the %s file was valid but the fingerprints were not checked.", relaysManifestFilePath))

			pterm.Println()

			result, _ := prompt.Show()

			if result == "no" {
				os.Exit(1)
			}
		}

		pterm.Println()
		spinner.Success("Relay binary verified")
	} else {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to verify the %s file, the SHA512 hash doesn't match the SHA512 hash in the %s file", path, relaysManifestFilePath))
		os.Exit(1)
	}
}
