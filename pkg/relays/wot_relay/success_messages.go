package wot_relay

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string, httpsEnabled bool) {
	successMsgParams := messages.SuccessMsgParams{Domain: domain, HTTPSEnabled: httpsEnabled, DataDirPath: DataDirPath, IndexFilePath: fmt.Sprintf("%s/%s/%s", network.WWWDirPath, domain, IndexFile), StaticDirPath: fmt.Sprintf("%s/%s/%s", network.WWWDirPath, domain, StaticDir), NginxConfigFilePath: NginxConfigFilePath, BinaryFilePath: BinaryFilePath, EnvFilePath: EnvFilePath, ServiceFilePath: ServiceFilePath, ServiceName: ServiceName, RelayName: RelayName, GitHubLink: GithubLink}
	messages.SuccessMessages(&successMsgParams)
}
