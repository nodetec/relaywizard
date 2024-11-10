package wot_relay

const GitRepoBranch = "v0.1.12"
const GitRepoURL = "https://github.com/bitvora/wot-relay.git"
const GitRepoTmpDirPath = "/tmp/wot-relay"
const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.4.0/wot-relay-0.1.12-x86_64-linux-gnu.tar.gz"
const BinaryName = "wot-relay"
const BinaryFilePath = "/usr/local/bin/wot-relay"
const NginxConfigFilePath = "/etc/nginx/conf.d/wot_relay.conf"
const DataDirPath = "/var/lib/wot-relay"
const ConfigDirPath = "/etc/wot-relay"
const TemplatesDirPath = "/etc/wot-relay/templates"
const IndexFilePath = "/etc/wot-relay/templates/index.html"
const TmpIndexFilePath = "/tmp/wot-relay/templates/index.html"
const StaticDirPath = "/etc/wot-relay/templates/static"
const TmpStaticDirPath = "/tmp/wot-relay/templates/static"
const ServiceName = "wot-relay"
const EnvFilePath = "/etc/wot-relay/wot-relay.env"
const EnvFileTemplate = `RELAY_NAME="WoT Relay"
RELAY_PUBKEY="{{.PubKey}}"
RELAY_DESCRIPTION="Stores only notes in your WoT"
RELAY_URL="{{.WSScheme}}://{{.Domain}}"
RELAY_ICON="https://pfp.nostr.build/56306a93a88d4c657d8a3dfa57b55a4ed65b709eee927b5dafaab4d5330db21f.png"
RELAY_CONTACT="{{.RelayContact}}"
DB_PATH="/var/lib/wot-relay/db"
INDEX_PATH="/etc/wot-relay/templates/index.html"
STATIC_PATH="/etc/wot-relay/templates/static"
REFRESH_INTERVAL_HOURS=3
MINIMUM_FOLLOWERS=1
ARCHIVAL_SYNC="FALSE"
ARCHIVE_REACTIONS="FALSE"
MAX_AGE_DAYS=0
`
const ServiceFilePath = "/etc/systemd/system/wot-relay.service"
const ServiceFileTemplate = `[Unit]
Description=WoT Relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
WorkingDirectory=/home/nostr
EnvironmentFile={{.EnvFilePath}}
ExecStart={{.BinaryFilePath}}
Restart=on-failure
MemoryMax=2G

[Install]
WantedBy=multi-user.target
`
const RelayName = "WoT Relay"
const GithubLink = "https://github.com/bitvora/wot-relay"
