package network

const SSHJailFilePath = "/etc/fail2ban/jail.d/sshd.local"
const SSHJailFileTemplate = `[sshd]
enabled = true
port = 22
findtime = 5m
bantime = 2h
maxentry = 3
ignoreip = 127.0.0.1/8 ::1
backend = systemd
`
const CertificateDirPath = "/etc/letsencrypt/live"
const FullchainFile = "fullchain.pem"
const PrivkeyFile = "privkey.pem"
const ChainFile = "chain.pem"
const CertificateArchiveDirPath = "/etc/letsencrypt/archive"
const FullchainArchiveFile = "fullchain1.pem"
const PrivkeyArchiveFile = "privkey1.pem"
const ChainArchiveFile = "chain1.pem"
const CertArchiveFile = "cert1.pem"
const NginxConfDirPath = "/etc/nginx/conf.d"
const WWWDirPath = "/var/www"
const AcmeChallengeDirPath = ".well-known/acme-challenge"
