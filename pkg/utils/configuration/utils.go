package configuration

import (
	"github.com/nodetec/rwz/pkg/utils/network"
	"github.com/pterm/pterm"
	"os"
	"text/template"
)

type EnvFileParams struct {
	Domain       string
	PortNumber   string
	HTTPSEnabled bool
	PrivKey      string
	PubKey       string
	RelayContact string
}

func CreateEnvFile(envFilePath, envTemplate string, envFileParams *EnvFileParams) {
	envFile, err := os.Create(envFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to create environment file: %v", err)
		os.Exit(1)
	}
	defer envFile.Close()

	envTmpl, err := template.New("env").Parse(envTemplate)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to parse environment template: %v", err)
		os.Exit(1)
	}

	WSScheme := network.WSEnabled(envFileParams.HTTPSEnabled)

	err = envTmpl.Execute(envFile, struct{ Domain, PortNumber, WSScheme, PrivKey, PubKey, RelayContact string }{Domain: envFileParams.Domain, PortNumber: envFileParams.PortNumber, WSScheme: WSScheme, PrivKey: envFileParams.PrivKey, PubKey: envFileParams.PubKey, RelayContact: envFileParams.RelayContact})
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to execute environment template: %v", err)
		os.Exit(1)
	}
}
