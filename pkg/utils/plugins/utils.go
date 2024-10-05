package plugins

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"text/template"
)

type PluginFileParams struct {
	Domain         string
	RelaySecretKey string
	ConfigFilePath string
	BinaryFilePath string
}

func CreatePluginFile(pluginFilePath, pluginTemplate string, pluginFileParams *PluginFileParams) {
	pluginFile, err := os.Create(pluginFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to create plugin file: %v", err))
		os.Exit(1)
	}
	defer pluginFile.Close()

	pluginTmpl, err := template.New("plugin").Parse(pluginTemplate)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to parse plugin template: %v", err))
		os.Exit(1)
	}

	err = pluginTmpl.Execute(pluginFile, struct{ Domain, RelaySecretKey, ConfigFilePath, BinaryFilePath string }{Domain: pluginFileParams.Domain, RelaySecretKey: pluginFileParams.RelaySecretKey, ConfigFilePath: pluginFileParams.ConfigFilePath, BinaryFilePath: pluginFileParams.BinaryFilePath})
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to execute plugin template: %v", err))
		os.Exit(1)
	}
}
