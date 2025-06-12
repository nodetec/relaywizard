package databases

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
func CheckStrfryBinaryAndDatabaseCompatibility(binaryName, configFilePath, supportedDatabaseVersionOutput, binaryVersion, supportedDatabaseVersion string) {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking %s binary and existing database compatibility...", binaryName))

	out, err := exec.Command(binaryName, "--config", configFilePath, "info").CombinedOutput()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to check %s binary and existing database compatibility: %v", binaryName, err)
		os.Exit(1)
	}

	databaseInfoOutput := string(out)

	databaseVersionSupported := strings.Contains(databaseInfoOutput, supportedDatabaseVersionOutput)

	if databaseVersionSupported {
		spinner.Success("Existing database version supported")
	} else {
		spinner.Warning(fmt.Sprintf("Existing database version is incompatible with %s version %s.", binaryName, binaryVersion))
		pterm.Println()
		pterm.Printfln(pterm.LightCyan("Upgrade your database to version %s.", supportedDatabaseVersion))
	}
}
