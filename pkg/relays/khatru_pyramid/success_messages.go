package khatru_pyramid

import (
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string) {
	const dataDir = "/var/lib/khatru-pyramid"
	const envFile = "/etc/systemd/system/khatru-pyramid.env"
	const serviceFile = "/etc/systemd/system/khatru-pyramid.service"
	const service = "khatru-pyramid"
	const relayName = "Khatru Pyramid"
	const githubLink = "https://github.com/github-tijlxyz/khatru-pyramid"

	successMsgParams := messages.SuccessMsgParams{Domain: domain, DataDir: dataDir, EnvFile: envFile, ServiceFile: serviceFile, Service: service, RelayName: relayName, GitHubLink: githubLink}
	messages.SuccessMessages(&successMsgParams)
}
