package khatru29

const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.5.0/relay29-0.5.1-khatru29-x86_64-linux-gnu.tar.gz"
const BinaryName = "khatru29"
const DataDirPath = "/var/lib/khatru29"
const DatabaseFilePath = "/var/lib/khatru29/db/data.mdb"
const DatabaseBackupsDirPath = "/var/lib/khatru29/db/backups"
const BackupFileNameBase = "data.mdb-bak"
const ConfigDirPath = "/etc/khatru29"
const ServiceName = "khatru29"
const EnvFilePath = "/etc/khatru29/khatru29.env"
const EnvFileTemplate = `DOMAIN="{{.Domain}}"
PORT="{{.PortNumber}}"
RELAY_NAME="Khatru29"
RELAY_PRIVKEY="{{.PrivKey}}"
RELAY_DESCRIPTION="Khatru29 Relay"
RELAY_CONTACT="{{.RelayContact}}"
DATABASE_PATH="/var/lib/khatru29/db"
`
const ServiceFilePath = "/etc/systemd/system/khatru29.service"
const ServiceFileTemplate = `[Unit]
Description=Khatru29 Relay Service
After=network.target

[Service]
Type=simple
User={{.RelayUser}}
Group={{.RelayUser}}
WorkingDirectory=/home/{{.RelayUser}}
EnvironmentFile={{.EnvFilePath}}
ExecStart={{.BinaryFilePath}}
Restart=on-failure

[Install]
WantedBy=multi-user.target
`
const GithubLink = "https://github.com/fiatjaf/relay29"
