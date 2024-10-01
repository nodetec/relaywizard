package strfry

import (
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string) {
	const dataDir = "/var/lib/strfry"
	const configFile = "/etc/strfry.conf"
	const serviceFile = "/etc/systemd/system/strfry.service"
	const service = "strfry"
	const relayName = "strfry"
	const githubLink = "https://github.com/hoytech/strfry"

	successMsgParams := messages.SuccessMsgParams{Domain: domain, DataDir: dataDir, ConfigFile: configFile, ServiceFile: serviceFile, Service: service, RelayName: relayName, GitHubLink: githubLink}
	messages.SuccessMessages(&successMsgParams)
}
