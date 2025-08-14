package network

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"text/template"
)

type SSHDConfigDConfigFileParams struct {
	Port                                string
	AllowOnlyPubkeyAuthenticationMethod string
	PasswordAuthentication              string
}

// Function to determine http scheme being used
func HTTPEnabled(httpsEnabled bool) string {
	if httpsEnabled {
		return "https"
	}
	return "http"
}

// Function to determine ws scheme being used
func WSEnabled(httpsEnabled bool) string {
	if httpsEnabled {
		return "wss"
	}
	return "ws"
}

func CreateSSHDConfigDConfigFile(sshdConfigDConfigFilePath, sshdConfigDConfigFileTemplate string, sshdConfigDConfigFileParams *SSHDConfigDConfigFileParams) {
	sshdConfigDConfigFile, err := os.Create(sshdConfigDConfigFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to create sshd_config.d config file: %v", err)
		os.Exit(1)
	}
	defer sshdConfigDConfigFile.Close()

	sshdConfigDConfigFileTmpl, err := template.New("config").Parse(sshdConfigDConfigFileTemplate)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to parse sshd_config.d config file template: %v", err)
		os.Exit(1)
	}

	err = sshdConfigDConfigFileTmpl.Execute(sshdConfigDConfigFile, struct{ Port, AllowOnlyPubkeyAuthenticationMethod, PasswordAuthentication string }{Port: sshdConfigDConfigFileParams.Port, AllowOnlyPubkeyAuthenticationMethod: sshdConfigDConfigFileParams.AllowOnlyPubkeyAuthenticationMethod, PasswordAuthentication: sshdConfigDConfigFileParams.PasswordAuthentication})
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to execute sshd_config.d config file template: %v", err)
		os.Exit(1)
	}
}

// Function to create jail files for the intrusion detection system
func CreateJailFile(jailFilePath, jailTemplate string) {
	jailFile, err := os.Create(jailFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to create jail file: %v", err)
		os.Exit(1)
	}
	defer jailFile.Close()

	jailTmpl, err := template.New("jail").Parse(jailTemplate)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to parse jail template: %v", err)
		os.Exit(1)
	}

	err = jailTmpl.Execute(jailFile, struct{}{})
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to execute jail template: %v", err)
		os.Exit(1)
	}
}

// TODO
// May have to add check for support for both IPv4 and IPv6 addresses
// May also have to check if one type of IP address overrides the other

// Function to list the network socket file(s) using a provided IP version, protocol, and port number
func ListNetworkSocketFilesUsingIPVersionProtocolAndPortNumber(ipVersion, protocol, portNumber string) string {
	networkSocketFiles := fmt.Sprintf("-i%s%s:%s", ipVersion, protocol, portNumber)

	out, err := exec.Command("lsof", "-Q", "-nP", networkSocketFiles).CombinedOutput()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to list the network socket file(s) for %s: %v", networkSocketFiles, err)
		os.Exit(1)
	}

	lsofOutput := string(out)

	return lsofOutput
}

// Function to list the network socket file(s) using a provided IP version, protocol, IP address, and port number
func ListNetworkSocketFilesUsingIPVersionIPAddressProtocolAndPortNumber(ipVersion, protocol, ipAddress, portNumber string) string {
	networkSocketFiles := fmt.Sprintf("-i%s%s@%s:%s", ipVersion, protocol, ipAddress, portNumber)

	out, err := exec.Command("lsof", "-Q", "-nP", networkSocketFiles).CombinedOutput()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to list the network socket file(s) for %s: %v", networkSocketFiles, err)
		os.Exit(1)
	}

	lsofOutput := string(out)

	return lsofOutput
}
