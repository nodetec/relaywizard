package network

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
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

func CreateSSHDConfigDConfigFile(currentUsername, sshdConfigDConfigFilePath, sshdConfigDConfigFileTemplate string, sshdConfigDConfigFileParams *SSHDConfigDConfigFileParams) {
	if currentUsername == relays.RootUser {
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
	} else {
		_, err := exec.Command("sudo", "touch", sshdConfigDConfigFilePath).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create sshd_config.d config file: %v", err)
			os.Exit(1)
		}

		_, err = exec.Command("sudo", "chmod", "0666", sshdConfigDConfigFilePath).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set permissions for sshd_config.d config file: %v", err)
			os.Exit(1)
		}

		sshdConfigDConfigFile, err := os.OpenFile(sshdConfigDConfigFilePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to open sshd_config.d config file: %v", err)
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

		_, err = exec.Command("sudo", "chmod", "0644", sshdConfigDConfigFilePath).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set permissions for sshd_config.d config file: %v", err)
			os.Exit(1)
		}
	}
}

// Function to create jail files for the intrusion detection system
func CreateJailFile(currentUsername, jailFilePath, jailTemplate string) {
	if currentUsername == relays.RootUser {
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
	} else {
		_, err := exec.Command("sudo", "touch", jailFilePath).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create jail file: %v", err)
			os.Exit(1)
		}

		_, err = exec.Command("sudo", "chmod", "0666", jailFilePath).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set permissions for jail file: %v", err)
			os.Exit(1)
		}

		jailFile, err := os.OpenFile(jailFilePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to open jail file: %v", err)
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

		_, err = exec.Command("sudo", "chmod", "0644", jailFilePath).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set permissions for jail file: %v", err)
			os.Exit(1)
		}
	}
}

// TODO
// May have to add check for support for both IPv4 and IPv6 addresses
// May also have to check if one type of IP address overrides the other

// Function to list the network socket file(s) using a provided IP version, protocol, and port number
func ListNetworkSocketFilesUsingIPVersionProtocolAndPortNumber(ipVersion, protocol, portNumber, currentUsername string) string {
	networkSocketFiles := fmt.Sprintf("-i%s%s:%s", ipVersion, protocol, portNumber)

	if currentUsername == relays.RootUser {
		out, err := exec.Command("lsof", "-Q", "-nP", networkSocketFiles).CombinedOutput()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				errorCode := exitError.ExitCode()
				// Network socket file(s) not found
				if errorCode == 1 {
					return ""
				} else {
					pterm.Println()
					pterm.Error.Printfln("Failed to list the network socket file(s) for %s: %v", networkSocketFiles, err)
					os.Exit(1)
				}
			}
		}
		lsofOutput := string(out)

		return lsofOutput
	} else {
		out, err := exec.Command("sudo", "lsof", "-Q", "-nP", networkSocketFiles).CombinedOutput()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				errorCode := exitError.ExitCode()
				// Network socket file(s) not found
				if errorCode == 1 {
					return ""
				} else {
					pterm.Println()
					pterm.Error.Printfln("Failed to list the network socket file(s) for %s: %v", networkSocketFiles, err)
					os.Exit(1)
				}
			}
		}
		lsofOutput := string(out)

		return lsofOutput
	}
}

// Function to list the network socket file(s) using a provided IP version, protocol, IP address, and port number
func ListNetworkSocketFilesUsingIPVersionIPAddressProtocolAndPortNumber(ipVersion, protocol, ipAddress, portNumber, currentUsername string) string {
	networkSocketFiles := fmt.Sprintf("-i%s%s@%s:%s", ipVersion, protocol, ipAddress, portNumber)

	if currentUsername == relays.RootUser {
		out, err := exec.Command("lsof", "-Q", "-nP", networkSocketFiles).CombinedOutput()

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				errorCode := exitError.ExitCode()
				// Network socket file(s) not found
				if errorCode == 1 {
					return ""
				} else {
					pterm.Println()
					pterm.Error.Printfln("Failed to list the network socket file(s) for %s: %v", networkSocketFiles, err)
					os.Exit(1)
				}
			}
		}

		lsofOutput := string(out)

		return lsofOutput
	} else {
		out, err := exec.Command("sudo", "lsof", "-Q", "-nP", networkSocketFiles).CombinedOutput()

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				errorCode := exitError.ExitCode()
				// Network socket file(s) not found
				if errorCode == 1 {
					return ""
				} else {
					pterm.Println()
					pterm.Error.Printfln("Failed to list the network socket file(s) for %s: %v", networkSocketFiles, err)
					os.Exit(1)
				}
			}
		}

		lsofOutput := string(out)

		return lsofOutput

	}
}
