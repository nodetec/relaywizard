package khatru_pyramid

const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.4.0/khatru-pyramid-0.1.0-x86_64-linux-gnu.tar.gz"
const BinaryName = "khatru-pyramid"
const BinaryFilePath = "/usr/local/bin/khatru-pyramid"
const NginxConfigFilePath = "/etc/nginx/conf.d/khatru_pyramid.conf"
const DataDirPath = "/var/lib/khatru-pyramid"
const ConfigDirPath = "/etc/khatru-pyramid"
const UsersFilePath = "/var/lib/khatru-pyramid/users.json"
const ServiceName = "khatru-pyramid"
const EnvFilePath = "/etc/khatru-pyramid/khatru-pyramid.env"
const EnvFileTemplate = `DOMAIN="{{.Domain}}"
PORT="3335"
DATABASE_PATH="/var/lib/khatru-pyramid/db"
USERDATA_PATH="/var/lib/khatru-pyramid/users.json"
MAX_INVITES_PER_PERSON="3"
RELAY_NAME="Khatru Pyramid"
RELAY_PUBKEY="{{.PubKey}}"
RELAY_DESCRIPTION="Khatru Pyramid Relay"
RELAY_CONTACT="{{.RelayContact}}"
`
const ServiceFilePath = "/etc/systemd/system/khatru-pyramid.service"
const ServiceFileTemplate = `[Unit]
Description=Khatru Pyramid Relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
WorkingDirectory=/home/nostr
EnvironmentFile={{.EnvFilePath}}
ExecStart={{.BinaryFilePath}}
Restart=on-failure

[Install]
WantedBy=multi-user.target
`
const RelayName = "Khatru Pyramid"
const GithubLink = "https://github.com/github-tijlxyz/khatru-pyramid"
