package configuration

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/nodetec/rwz/pkg/utils/network"
	"github.com/pterm/pterm"
)

type EnvFileParams struct {
	Domain       string
	PortNumber   string
	HTTPSEnabled bool
	PrivKey      string
	PubKey       string
	RelayContact string
}

func CreateEnvFile(currentUsername, envFilePath, envTemplate string, envFileParams *EnvFileParams) {
	if currentUsername == relays.RootUser {
		envFile, err := os.Create(envFilePath)
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to create environment file: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to create environment file: %v", err)
			os.Exit(1)
		}
		defer envFile.Close()

		envTmpl, err := template.New("env").Parse(envTemplate)
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to parse environment template: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to parse environment template: %v", err)
			os.Exit(1)
		}

		WSScheme := network.WSEnabled(envFileParams.HTTPSEnabled)

		err = envTmpl.Execute(envFile, struct{ Domain, PortNumber, WSScheme, PrivKey, PubKey, RelayContact string }{Domain: envFileParams.Domain, PortNumber: envFileParams.PortNumber, WSScheme: WSScheme, PrivKey: envFileParams.PrivKey, PubKey: envFileParams.PubKey, RelayContact: envFileParams.RelayContact})
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to execute environment template: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to execute environment template: %v", err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "touch", envFilePath).Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to create environment file: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to create environment file: %v", err)
			os.Exit(1)
		}

		files.SetPermissionsUsingLinux(currentUsername, envFilePath, "0666")

		envFile, err := os.OpenFile(envFilePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to open environment file: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to open environment file: %v", err)
			os.Exit(1)
		}
		defer envFile.Close()

		envTmpl, err := template.New("env").Parse(envTemplate)
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to parse environment template: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to parse environment template: %v", err)
			os.Exit(1)
		}

		WSScheme := network.WSEnabled(envFileParams.HTTPSEnabled)

		err = envTmpl.Execute(envFile, struct{ Domain, PortNumber, WSScheme, PrivKey, PubKey, RelayContact string }{Domain: envFileParams.Domain, PortNumber: envFileParams.PortNumber, WSScheme: WSScheme, PrivKey: envFileParams.PrivKey, PubKey: envFileParams.PubKey, RelayContact: envFileParams.RelayContact})
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to execute environment template: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to execute environment template: %v", err)
			os.Exit(1)
		}
	}
}
