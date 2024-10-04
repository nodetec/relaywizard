package templates

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"text/template"
)

type IndexFileParams struct {
	Domain       string
	HTTPSEnabled bool
	PubKey       string
}

func CreateIndexFile(indexFilePath, indexTemplate string, indexFileParams *IndexFileParams) {
	indexFile, err := os.Create(indexFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to create index.html file: %v", err))
		os.Exit(1)
	}
	defer indexFile.Close()

	indexTmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to parse index.html template: %v", err))
		os.Exit(1)
	}

	var HTTPProtocol string
	if indexFileParams.HTTPSEnabled {
		HTTPProtocol = "https"
	} else {
		HTTPProtocol = "http"
	}

	err = indexTmpl.Execute(indexFile, struct{ Domain, HTTPProtocol, PubKey string }{Domain: indexFileParams.Domain, HTTPProtocol: HTTPProtocol, PubKey: indexFileParams.PubKey})
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to execute index.html template: %v", err))
		os.Exit(1)
	}
}
