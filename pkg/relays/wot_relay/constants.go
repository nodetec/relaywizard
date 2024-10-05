package wot_relay

const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.3.0/wot-relay-0.1.12-x86_64-linux-gnu.tar.gz"
const BinaryName = "wot-relay"
const BinaryFilePath = "/usr/local/bin/wot-relay"
const NginxConfigFilePath = "/etc/nginx/conf.d/wot_relay.conf"
const TemplatesDirPath = "/etc/wot-relay/templates"
const IndexFilePath = "/etc/wot-relay/templates/index.html"
const IndexFileTemplate = `<!doctype html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			<title>WoT Relay</title>
			<meta name="description" content="WoT Relay" />
			<link href="{{.HTTPProtocol}}://{{.Domain}}" rel="canonical" />
	</head>
	<body>
		<main>
			<div>
				<div>
					<span>WoT Relay</span>
				</div>
				<div>
					<span>Domain: {{.Domain}}</span>
				</div>
				<div>
					<span>Pubkey: {{.PubKey}}</span>
				</div>
			</div>
		</main>
	</body>
</html>
`
const StaticDirPath = "/etc/wot-relay/templates/static"
const DataDirPath = "/var/lib/wot-relay"
const ServiceName = "wot-relay"
const EnvFilePath = "/etc/systemd/system/wot-relay.env"
const EnvFileTemplate = `RELAY_NAME="WoT Relay"
RELAY_DESCRIPTION="WoT Nostr Relay"
RELAY_ICON="https://pfp.nostr.build/56306a93a88d4c657d8a3dfa57b55a4ed65b709eee927b5dafaab4d5330db21f.png"
RELAY_URL="{{.WSProtocol}}://{{.Domain}}"
RELAY_PUBKEY="{{.PubKey}}"
RELAY_CONTACT="{{.RelayContact}}"
INDEX_PATH="/etc/wot-relay/templates/index.html"
STATIC_PATH="/etc/wot-relay/templates/static"
DB_PATH="/var/lib/wot-relay/db"
REFRESH_INTERVAL_HOURS=24
MINIMUM_FOLLOWERS=3
ARCHIVAL_SYNC="FALSE"
ARCHIVE_REACTIONS="FALSE"
`
const ServiceFilePath = "/etc/systemd/system/wot-relay.service"
const ServiceFileTemplate = `[Unit]
Description=WoT Nostr Relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
WorkingDirectory=/home/nostr
EnvironmentFile={{.EnvFilePath}}
ExecStart={{.BinaryFilePath}}
Restart=on-failure
MemoryHigh=512M
MemoryMax=1G

[Install]
WantedBy=multi-user.target
`
const RelayName = "WoT Relay"
const GithubLink = "https://github.com/bitvora/wot-relay"
