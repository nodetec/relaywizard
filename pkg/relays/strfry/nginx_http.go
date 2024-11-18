package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to configure Nginx for HTTP
func ConfigureNginxHttp(domainName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx for HTTP...")

	files.RemoveFile(NginxConfigFilePath)

	directories.CreateDirectory(fmt.Sprintf("%s/%s", network.WWWDirPath, domainName), 0755)
	directories.CreateDirectory(fmt.Sprintf("%s/%s/%s/", network.WWWDirPath, domainName, network.AcmeChallengeDirPath), 0755)
	directories.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, fmt.Sprintf("%s/%s", network.WWWDirPath, domainName))

	configContent := fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream strfry_websocket {
    server 127.0.0.1:7777;
}

server {
    listen 80;
    listen [::]:80;
    server_name %s;

    root %s/%s;

    location /%s/ {
        default_type "text/plain";
    }

    location / {
        proxy_pass http://strfry_websocket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
        # First attempt to serve request as file, then
        # as directory, then fall back to displaying 404.
        try_files $uri $uri/ =404;
    }

    # Only return Nginx in server header
    server_tokens off;

    #### Security Headers ####
    # Test configuration:
    # https://securityheaders.com/
    # https://observatory.mozilla.org/
		add_header X-Frame-Options DENY;

    # Avoid MIME type sniffing
    add_header X-Content-Type-Options "nosniff" always;

    add_header Referrer-Policy "no-referrer" always;

    add_header X-XSS-Protection "1; mode=block" always;

    add_header Permissions-Policy "geolocation=(), midi=(), sync-xhr=(), microphone=(), camera=(), magnetometer=(), gyroscope=(), fullscreen=(self), payment=()" always;

    #### Content-Security-Policy (CSP) ####
    add_header Content-Security-Policy "base-uri 'self'; object-src 'none'; frame-ancestors 'none';" always;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name %s;

    root %s/%s;

    location / {
        return 301 http://%s$request_uri;
    }
}
`, domainName, network.WWWDirPath, domainName, network.AcmeChallengeDirPath, domainName, network.WWWDirPath, domainName, domainName)

	files.WriteFile(NginxConfigFilePath, configContent, 0644)
	files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, NginxConfigFilePath)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTP")
}
