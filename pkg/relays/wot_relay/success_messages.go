package wot_relay

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string, httpsEnabled bool) {
	successMsgParams := messages.SuccessMsgParams{Domain: domain, HTTPSEnabled: httpsEnabled, DataDirPath: DataDirPath, IndexFilePath: fmt.Sprintf("%s/%s/%s", network.WWWDirPath, domain, IndexFile), StaticDirPath: fmt.Sprintf("%s/%s/%s", network.WWWDirPath, domain, StaticDir), NginxConfigFilePath: relays.WotRelayNginxConfigFilePath, BinaryFilePath: relays.WotRelayBinaryFilePath, EnvFilePath: EnvFilePath, ServiceFilePath: ServiceFilePath, ServiceName: ServiceName, RelayName: relays.WotRelayName, GitHubLink: GithubLink}
	messages.SuccessMessages(&successMsgParams)
}
