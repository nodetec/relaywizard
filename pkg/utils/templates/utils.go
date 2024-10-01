package templates

import (
	"log"
	"os"
	"text/template"
)

type IndexFileParams struct {
	Domain string
	PubKey string
}

func CreateIndexFile(indexFilePath, indexTemplate string, indexFileParams *IndexFileParams) {
	indexFile, err := os.Create(indexFilePath)
	if err != nil {
		log.Fatalf("Error creating index.html file: %v", err)
	}
	defer indexFile.Close()

	indexTmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		log.Fatalf("Error parsing index.html template: %v", err)
	}

	err = indexTmpl.Execute(indexFile, struct{ Domain, PubKey string }{Domain: indexFileParams.Domain, PubKey: indexFileParams.PubKey})
	if err != nil {
		log.Fatalf("Error executing index.html template: %v", err)
	}
}
