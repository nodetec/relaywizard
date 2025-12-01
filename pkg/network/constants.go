package network

const SSHDirPath = "/etc/ssh"
const SSHDConfigFilePath = "/etc/ssh/sshd_config"
const SSHDConfigFileIncludeAllSSHDConfigDConfFilesLinePattern = "Include /etc/ssh/sshd_config.d/\\*.conf"
const SSHDConfigDDirPath = "/etc/ssh/sshd_config.d"
const SSHDConfigDRWZConfigFilePath = "/etc/ssh/sshd_config.d/99-rwz.conf"
const RootHiddenSSHDirPath = "/root/.ssh"
const RootHiddenSSHAuthorizedKeysFilePath = "/root/.ssh/authorized_keys"
const DefaultSSHPort = "22"
const AllowOnlyPubkeyAuthenticationMethod = "AuthenticationMethods publickey"
const SSHDConfigDRWZConfigFileTemplate = `Port {{.Port}}
AddressFamily any
ListenAddress 0.0.0.0
ListenAddress ::

HostKey /etc/ssh/ssh_host_ed25519_key
HostKey /etc/ssh/ssh_host_rsa_key
HostKey /etc/ssh/ssh_host_ecdsa_key

# Ciphers and keying
RekeyLimit default none

KexAlgorithms sntrup761x25519-sha512@openssh.com,curve25519-sha256,curve25519-sha256@libssh.org,ecdh-sha2-nistp521,ecdh-sha2-nistp384,ecdh-sha2-nistp256,diffie-hellman-group-exchange-sha256

Ciphers chacha20-poly1305@openssh.com,aes256-gcm@openssh.com,aes128-gcm@openssh.com,aes256-ctr,aes192-ctr,aes128-ctr

MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com,umac-128-etm@openssh.com,hmac-sha2-512,hmac-sha2-256,umac-128@openssh.com

# Logging
SyslogFacility AUTH
LogLevel VERBOSE

# Authentication:

LoginGraceTime 2m
PermitRootLogin yes
StrictModes yes
MaxAuthTries 6
MaxSessions 10

{{.AllowOnlyPubkeyAuthenticationMethod}}

PubkeyAuthentication yes

AuthorizedKeysFile	.ssh/authorized_keys

AuthorizedPrincipalsFile none

AuthorizedKeysCommand none
AuthorizedKeysCommandUser nobody

# For this to work you will also need host keys in /etc/ssh/ssh_known_hosts
HostbasedAuthentication no
# Change to yes if you don't trust ~/.ssh/known_hosts for
# HostbasedAuthentication
IgnoreUserKnownHosts no
# Don't read the user's ~/.rhosts and ~/.shosts files
IgnoreRhosts yes

# To disable tunneled clear text passwords, change to no here!
PasswordAuthentication {{.PasswordAuthentication}}
PermitEmptyPasswords no

# Change to yes to enable challenge-response passwords (beware issues with
# some PAM modules and threads)
KbdInteractiveAuthentication no

# Kerberos options
KerberosAuthentication no
KerberosOrLocalPasswd yes
KerberosTicketCleanup yes

# GSSAPI options
GSSAPIAuthentication no
GSSAPICleanupCredentials yes
GSSAPIStrictAcceptorCheck yes
GSSAPIKeyExchange no

# Set this to 'yes' to enable PAM authentication, account processing,
# and session processing. If this is enabled, PAM authentication will
# be allowed through the KbdInteractiveAuthentication and
# PasswordAuthentication.  Depending on your PAM configuration,
# PAM authentication via KbdInteractiveAuthentication may bypass
# the setting of "PermitRootLogin prohibit-password".
# If you just want the PAM account and session checks to run without
# PAM authentication, then enable this but set PasswordAuthentication
# and KbdInteractiveAuthentication to 'no'.
UsePAM no

AllowAgentForwarding yes
AllowTcpForwarding yes
GatewayPorts no
X11Forwarding no
X11DisplayOffset 10
X11UseLocalhost yes
PermitTTY yes
PrintMotd yes
PrintLastLog yes
TCPKeepAlive yes
PermitUserEnvironment no
Compression delayed
ClientAliveInterval 300
ClientAliveCountMax 3
UseDNS no
PidFile /run/sshd.pid
MaxStartups 10:30:100
PermitTunnel no
ChrootDirectory none
VersionAddendum none

# no default banner path
Banner none

# Allow client to pass locale environment variables
AcceptEnv LANG LC_*
`

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
const WellKnownDir = ".well-known"
const AcmeChallengeDir = "acme-challenge"
const AcmeChallengeDirPath = ".well-known/acme-challenge"
