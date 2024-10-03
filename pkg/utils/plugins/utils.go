package plugins

import (
	"log"
	"os"
	"text/template"
)

type PluginFileParams struct {
	Domain         string
	RelaySecretKey string
}

func CreatePluginFile(pluginFilePath, pluginTemplate string, pluginFileParams *PluginFileParams) {
	pluginFile, err := os.Create(pluginFilePath)
	if err != nil {
		log.Fatalf("Error creating plugin file: %v", err)
	}
	defer pluginFile.Close()

	pluginTmpl, err := template.New("plugin").Parse(pluginTemplate)
	if err != nil {
		log.Fatalf("Error parsing plugin template: %v", err)
	}

	err = pluginTmpl.Execute(pluginFile, struct{ Domain, RelaySecretKey string }{Domain: pluginFileParams.Domain, RelaySecretKey: pluginFileParams.RelaySecretKey})
	if err != nil {
		log.Fatalf("Error executing plugin template: %v", err)
	}
}
