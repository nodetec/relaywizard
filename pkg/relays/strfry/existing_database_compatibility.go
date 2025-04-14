package strfry

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"strings"
)

// Warn the user if strfry binary and database version are incompatible
// TODO
// Build an earlier version to get an older database version to see how the output should be handled
func CheckBinaryAndDatabaseCompatibility() {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking %s binary and existing database compatibility...", BinaryName))

	out, err := exec.Command(BinaryName, "--config", ConfigFilePath, "info").CombinedOutput()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to check %s binary and existing database compatibility: %v", BinaryName, err))
		os.Exit(1)
	}

	databaseInfoOutput := string(out)

	databaseVersionSupported := strings.Contains(databaseInfoOutput, SupportedDatabaseVersionOutput)

	if databaseVersionSupported {
		spinner.Success("Existing database version supported")
	} else {
		spinner.Warning(fmt.Sprintf("Existing database version is incompatible with %s version %s.", BinaryName, BinaryVersion))
		pterm.Println()
		pterm.Println(pterm.Cyan(fmt.Sprintf("Upgrade your database to version %s.", SupportedDatabaseVersion)))
	}
}
