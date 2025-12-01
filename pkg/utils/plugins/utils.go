package plugins

import (
	"os"
	"os/exec"
	"text/template"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

type PluginFileParams struct {
	Domain         string
	RelaySecretKey string
	ConfigFilePath string
	BinaryFilePath string
}

func CreatePluginFile(currentUsername, pluginFilePath, pluginTemplate string, pluginFileParams *PluginFileParams) {
	if currentUsername == relays.RootUser {
		pluginFile, err := os.Create(pluginFilePath)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create plugin file: %v", err)
			os.Exit(1)
		}
		defer pluginFile.Close()

		pluginTmpl, err := template.New("plugin").Parse(pluginTemplate)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to parse plugin template: %v", err)
			os.Exit(1)
		}

		err = pluginTmpl.Execute(pluginFile, struct{ Domain, RelaySecretKey, ConfigFilePath, BinaryFilePath string }{Domain: pluginFileParams.Domain, RelaySecretKey: pluginFileParams.RelaySecretKey, ConfigFilePath: pluginFileParams.ConfigFilePath, BinaryFilePath: pluginFileParams.BinaryFilePath})
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to execute plugin template: %v", err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "touch", pluginFilePath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create plugin file: %v", err)
			os.Exit(1)
		}

		files.SetPermissionsUsingLinux(currentUsername, pluginFilePath, "0666")

		pluginFile, err := os.OpenFile(pluginFilePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to open plugin file: %v", err)
			os.Exit(1)
		}
		defer pluginFile.Close()

		pluginTmpl, err := template.New("plugin").Parse(pluginTemplate)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to parse plugin template: %v", err)
			os.Exit(1)
		}

		err = pluginTmpl.Execute(pluginFile, struct{ Domain, RelaySecretKey, ConfigFilePath, BinaryFilePath string }{Domain: pluginFileParams.Domain, RelaySecretKey: pluginFileParams.RelaySecretKey, ConfigFilePath: pluginFileParams.ConfigFilePath, BinaryFilePath: pluginFileParams.BinaryFilePath})
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to execute plugin template: %v", err)
			os.Exit(1)
		}
	}
}
