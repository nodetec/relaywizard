package khatru29

import (
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string, httpsEnabled bool) {
	const dataDir = "/var/lib/khatru29"
	const envFile = "/etc/systemd/system/khatru29.env"
	const serviceFile = "/etc/systemd/system/khatru29.service"
	const service = "khatru29"
	const relayName = "Khatru29"
	const githubLink = "https://github.com/fiatjaf/relay29/tree/master"

	successMsgParams := messages.SuccessMsgParams{Domain: domain, HTTPSEnabled: httpsEnabled, DataDir: dataDir, EnvFile: envFile, ServiceFile: serviceFile, Service: service, RelayName: relayName, GitHubLink: githubLink}
	messages.SuccessMessages(&successMsgParams)
}
