package network

import (
	"fmt"
	"os"
	"strings"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/nodetec/rwz/pkg/utils/network"
	"github.com/nodetec/rwz/pkg/utils/programs"
	"github.com/pterm/pterm"
)

func determineFirstPIDFromLsofOutput(currentUsername, lsofOutput string) string {
	var firstPIDFromLsofOutput string

	if lsofOutput == "" {
		return firstPIDFromLsofOutput
	} else {
		lsofOutputSplitByNewLine := strings.Split(lsofOutput, "\n")
		lsofOutputSplitByNewLineLength := len(lsofOutputSplitByNewLine)

		if lsofOutputSplitByNewLineLength > 1 {
			firstNetworkSocketFileOutput := lsofOutputSplitByNewLine[1]
			firstNetworkSocketFileOutputSplitBySpace := strings.Split(firstNetworkSocketFileOutput, " ")
			firstNetworkSocketFileOutputSplitBySpaceLength := len(firstNetworkSocketFileOutputSplitBySpace)

			// Assuming there could potentially be multiple spaces between the lsof COMMAND and PID column outputs and that the second
			// lsof row output will start with the COMMAND value followed by at least one space before the PID value
			// Start at the second element if it's an empty string keep looping until the string isn't empty which will give the PID
			if firstNetworkSocketFileOutputSplitBySpaceLength > 1 {
				for i := 1; i < firstNetworkSocketFileOutputSplitBySpaceLength; i++ {
					if firstNetworkSocketFileOutputSplitBySpace[i] != "" {
						firstPIDFromLsofOutput = firstNetworkSocketFileOutputSplitBySpace[i]
						break
					}
				}
			} else {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, "Failed to determine the first PID from the lsof output")
				pterm.Println()
				pterm.Error.Println("Failed to determine the first PID from the lsof output")
				os.Exit(1)
			}
		} else if lsofOutputSplitByNewLineLength == 1 {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, "Failed to parse lsof output")
			pterm.Println()
			pterm.Error.Println("Failed to parse lsof output")
			os.Exit(1)
		}

		return firstPIDFromLsofOutput
	}
}

// Function to check if the selected relay's port is available to use
func CheckPort(currentUsername, selectedRelayOption string) {
	spinner, _ := pterm.DefaultSpinner.Start("Checking relay port availability...")

	var relayBinaryFilePath string
	protocol := "TCP"
	var portNumber string
	var firstPIDFromLsofOutput string
	var lsofOutput string

	// TODO
	// May have to check both IPv4 and IPv6 addresses for Khatru Pyramid, WoT Relay, and Khatru29
	if selectedRelayOption == relays.KhatruPyramidRelayName {
		relayBinaryFilePath = relays.KhatruPyramidBinaryFilePath
		portNumber = relays.KhatruPyramidPortNumber

		lsofOutput = network.ListNetworkSocketFilesUsingIPVersionProtocolAndPortNumber(currentUsername, "6", protocol, portNumber)
	} else if selectedRelayOption == relays.NostrRsRelayName {
		relayBinaryFilePath = relays.NostrRsRelayBinaryFilePath
		portNumber = relays.NostrRsRelayPortNumber

		lsofOutput = network.ListNetworkSocketFilesUsingIPVersionIPAddressProtocolAndPortNumber(currentUsername, "4", protocol, relays.NostrRsRelayIPv4Address, portNumber)
	} else if selectedRelayOption == relays.StrfryRelayName {
		relayBinaryFilePath = relays.StrfryBinaryFilePath
		portNumber = relays.StrfryPortNumber

		lsofOutput = network.ListNetworkSocketFilesUsingIPVersionIPAddressProtocolAndPortNumber(currentUsername, "4", protocol, relays.StrfryIPv4Address, portNumber)
	} else if selectedRelayOption == relays.WotRelayName {
		relayBinaryFilePath = relays.WotRelayBinaryFilePath
		portNumber = relays.WotRelayPortNumber

		lsofOutput = network.ListNetworkSocketFilesUsingIPVersionProtocolAndPortNumber(currentUsername, "6", protocol, portNumber)
	} else if selectedRelayOption == relays.Khatru29RelayName {
		relayBinaryFilePath = relays.Khatru29BinaryFilePath
		portNumber = relays.Khatru29PortNumber

		lsofOutput = network.ListNetworkSocketFilesUsingIPVersionProtocolAndPortNumber(currentUsername, "6", protocol, portNumber)
	} else if selectedRelayOption == relays.Strfry29RelayName {
		relayBinaryFilePath = relays.Strfry29BinaryFilePath
		portNumber = relays.Strfry29PortNumber

		lsofOutput = network.ListNetworkSocketFilesUsingIPVersionIPAddressProtocolAndPortNumber(currentUsername, "4", protocol, relays.Strfry29IPv4Address, portNumber)
	}

	firstPIDFromLsofOutput = determineFirstPIDFromLsofOutput(currentUsername, lsofOutput)

	// TODO
	// Look into clarifying the explanation to end user
	if files.FileExists(relayBinaryFilePath) {
		pidsOfRelayBinary := programs.DeterminePidsOfProgram(currentUsername, relayBinaryFilePath)

		if pidsOfRelayBinary != nil {
			spinner.UpdateText(fmt.Sprintf("Relay binary located at %s is currently running...", relayBinaryFilePath))

			if firstPIDFromLsofOutput != "" {
				pidsOfRelayBinaryLength := len(pidsOfRelayBinary)

				for i := range pidsOfRelayBinaryLength {
					if firstPIDFromLsofOutput == pidsOfRelayBinary[i] {
						spinner.UpdateText(fmt.Sprintf("Relay binary located at %s is currently running on the default port, overwriting relay installation...", relayBinaryFilePath))
						break
					}

					if i == pidsOfRelayBinaryLength-1 {
						logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Unable to bind to port number %s, another process has a connection open on the required TCP IP address and port combination", portNumber))
						pterm.Println()
						pterm.Error.Printfln("Unable to bind to port number %s, another process has a connection open on the required TCP IP address and port combination", portNumber)
						spinner.Fail("Relay port unavailable")
						os.Exit(1)
					}
				}
			} else {
				// TODO
				// Give the user the option to continue with installation or not if it's determined they're using a custom port
				spinner.UpdateText(fmt.Sprintf("Relay binary located at %s is currently running on a custom port, overwriting any related relay installation files...", relayBinaryFilePath))
			}
		} else {
			spinner.UpdateText(fmt.Sprintf("Unable to find process for relay binary located at %s...", relayBinaryFilePath))

			if firstPIDFromLsofOutput != "" {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Unable to bind to port number %s, another process has a connection open on the required TCP IP address and port combination", portNumber))
				pterm.Println()
				pterm.Error.Printfln("Unable to bind to port number %s, another process has a connection open on the required TCP IP address and port combination", portNumber)
				spinner.Fail("Relay port unavailable")
				os.Exit(1)
			}
		}
	} else {
		spinner.UpdateText(fmt.Sprintf("Relay binary located at %s not found...", relayBinaryFilePath))

		if firstPIDFromLsofOutput != "" {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Unable to bind to port number %s, another process has a connection open on the required TCP IP address and port combination", portNumber))
			pterm.Println()
			pterm.Error.Printfln("Unable to bind to port number %s, another process has a connection open on the required TCP IP address and port combination", portNumber)
			spinner.Fail("Relay port unavailable")
			os.Exit(1)
		}
	}

	spinner.Success("Relay port available")
}
