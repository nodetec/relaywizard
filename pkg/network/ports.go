package network

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/programs"
	"github.com/pterm/pterm"
)

// Function to check if the selected relay's port(s) is/are available to use
func CheckPorts(selectedRelayOption string) {
	spinner, _ := pterm.DefaultSpinner.Start("Checking relay port availability...")

	var relayBinaryFilePath string

	if selectedRelayOption == relays.KhatruPyramidRelayName {
		relayBinaryFilePath = relays.KhatruPyramidBinaryFilePath
	} else if selectedRelayOption == relays.NostrRsRelayName {
		relayBinaryFilePath = relays.NostrRsRelayBinaryFilePath
	} else if selectedRelayOption == relays.StrfryRelayName {
		relayBinaryFilePath = relays.StrfryBinaryFilePath
	} else if selectedRelayOption == relays.WotRelayName {
		relayBinaryFilePath = relays.WotRelayBinaryFilePath
	} else if selectedRelayOption == relays.Khatru29RelayName {
		relayBinaryFilePath = relays.Khatru29BinaryFilePath
	} else if selectedRelayOption == relays.Strfry29RelayName {
		relayBinaryFilePath = relays.Strfry29BinaryFilePath
	}

	if files.FileExists(relayBinaryFilePath) {
		pidsOfRelayBinary := programs.DeterminePidsOfProgram(relayBinaryFilePath)

		if pidsOfRelayBinary != nil {
			spinner.UpdateText(fmt.Sprintf("Relay binary located at %s is currently running...", relayBinaryFilePath))
		} else {
			spinner.UpdateText(fmt.Sprintf("Relay binary located at %s is currently not running...", relayBinaryFilePath))
		}
	} else {
		spinner.UpdateText(fmt.Sprintf("Relay binary located at %s not found...", relayBinaryFilePath))
	}

	spinner.Success("Relay port(s) available")
}
