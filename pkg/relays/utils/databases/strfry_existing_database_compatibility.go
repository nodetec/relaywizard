package databases

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/pterm/pterm"
)

// Warn the user if strfry binary and database version are incompatible
// TODO
// Build an earlier version to get an older database version to see how the output should be handled
func CheckStrfryBinaryAndDatabaseCompatibility(currentUsername, binaryName, configFilePath, supportedDatabaseVersionOutput, binaryVersion, supportedDatabaseVersion string) {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking %s binary and existing database compatibility...", binaryName))

	var strfryInfoOutput []byte
	if currentUsername == relays.RootUser {
		out, err := exec.Command(binaryName, "--config", configFilePath, "info").CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to check %s binary and existing database compatibility: %v", binaryName, err)
			os.Exit(1)
		}
		strfryInfoOutput = out
	} else {
		out, err := exec.Command("sudo", binaryName, "--config", configFilePath, "info").CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to check %s binary and existing database compatibility: %v", binaryName, err)
			os.Exit(1)
		}
		strfryInfoOutput = out
	}

	databaseInfoOutput := string(strfryInfoOutput)

	databaseVersionSupported := strings.Contains(databaseInfoOutput, supportedDatabaseVersionOutput)

	if databaseVersionSupported {
		spinner.Success("Existing database version supported")
	} else {
		spinner.Warning(fmt.Sprintf("Existing database version is incompatible with %s version %s.", binaryName, binaryVersion))
		pterm.Println()
		pterm.Printfln(pterm.LightCyan("Upgrade your database to version %s."), supportedDatabaseVersion)
	}
}
