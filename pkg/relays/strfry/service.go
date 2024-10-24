package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to set up the relay service
func SetupRelayService(domain, relayContact string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay service...")

	// Ensure the data directory exists and set ownership
	spinner.UpdateText("Creating data directory...")
	directories.CreateDirectory(DataDirPath, 0755)

	// Use chown command to set ownership of the data directory to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, DataDirPath)

	// Ensure the config directory exists and set ownership
	spinner.UpdateText("Creating config directory...")
	directories.CreateDirectory(ConfigDirPath, 0755)

	// Use chown command to set ownership of the config directory to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, ConfigDirPath)

	// Check for and remove existing config file
	files.RemoveFile(ConfigFilePath)

	// Check if the service file exists and remove it if it does
	files.RemoveFile(ServiceFilePath)

	// Construct the sed command to change the db path
	files.InPlaceEdit(fmt.Sprintf(`s|db = ".*"|db = "%s"|`, DataDirPath), TmpConfigFilePath)

	// Construct the sed command to change the nofiles limit
	// TODO
	// Determine system hard limit
	// Determine preferred nofiles value
	files.InPlaceEdit(`s|nofiles = .*|nofiles = 0|`, TmpConfigFilePath)

	// Construct the sed command to change the contact
	files.InPlaceEdit(fmt.Sprintf(`s|contact = ".*"|contact = "%s"|`, relayContact), TmpConfigFilePath)

	// Copy config file to config directory
	files.CopyFile(TmpConfigFilePath, ConfigDirPath)

	// Use chown command to set ownership of the config file to the nostr user
	files.SetOwnerAndGroup(relays.User, relays.User, ConfigFilePath)

	// Create the systemd service file
	spinner.UpdateText("Creating service file...")
	serviceFileParams := systemd.ServiceFileParams{BinaryFilePath: BinaryFilePath, ConfigFilePath: ConfigFilePath}
	systemd.CreateServiceFile(ServiceFilePath, ServiceFileTemplate, &serviceFileParams)

	// Reload systemd to apply the new service
	spinner.UpdateText("Reloading systemd daemon...")
	systemd.Reload()

	// Enable and start the Nostr relay service
	spinner.UpdateText("Enabling and starting service...")
	systemd.EnableService(ServiceName)
	systemd.StartService(ServiceName)

	spinner.Success("Nostr relay service configured")
}
