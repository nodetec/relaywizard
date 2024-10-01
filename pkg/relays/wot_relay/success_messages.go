package wot_relay

import (
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string) {
	const dataDir = "/var/lib/wot-relay"
	const indexFile = "/etc/wot-relay/templates/index.html"
	const staticDir = "/etc/wot-relay/templates/static"
	const envFile = "/etc/systemd/system/wot-relay.env"
	const serviceFile = "/etc/systemd/system/wot-relay.service"
	const service = "wot-relay"
	const relayName = "WoT Relay"
	const githubLink = "https://github.com/bitvora/wot-relay"

	successMsgParams := messages.SuccessMsgParams{Domain: domain, DataDir: dataDir, IndexFile: indexFile, StaticDir: staticDir, EnvFile: envFile, ServiceFile: serviceFile, Service: service, RelayName: relayName, GitHubLink: githubLink}
	messages.SuccessMessages(&successMsgParams)
}
