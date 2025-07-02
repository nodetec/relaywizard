package network

import (
	"github.com/pterm/pterm"
	"os"
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
