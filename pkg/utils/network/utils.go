package network

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

type SSHDConfigDConfigFileParams struct {
	Port                                string
	AllowOnlyPubkeyAuthenticationMethod string
	PasswordAuthentication              string
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
		err := exec.Command("sudo", "touch", sshdConfigDConfigFilePath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create sshd_config.d config file: %v", err)
			os.Exit(1)
		}

		files.SetPermissionsUsingLinux(currentUsername, sshdConfigDConfigFilePath, "0666")

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

		files.SetPermissionsUsingLinux(currentUsername, sshdConfigDConfigFilePath, "0644")
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
		err := exec.Command("sudo", "touch", jailFilePath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create jail file: %v", err)
			os.Exit(1)
		}

		files.SetPermissionsUsingLinux(currentUsername, jailFilePath, "0666")

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

		files.SetPermissionsUsingLinux(currentUsername, jailFilePath, "0644")
	}
}

// TODO
// May have to add check for support for both IPv4 and IPv6 addresses
// May also have to check if one type of IP address overrides the other

// Function to list the network socket file(s) using a provided IP version, protocol, and port number
func ListNetworkSocketFilesUsingIPVersionProtocolAndPortNumber(currentUsername, ipVersion, protocol, portNumber string) string {
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
func ListNetworkSocketFilesUsingIPVersionIPAddressProtocolAndPortNumber(currentUsername, ipVersion, protocol, ipAddress, portNumber string) string {
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
