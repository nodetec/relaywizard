package khatru29

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to configure nginx for HTTP
func ConfigureNginxHttp(domainName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring nginx for HTTP...")

	const configFilePath = "/etc/nginx/conf.d/khatru29.conf"

	var configContent string

	directories.CreateDirectory(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", domainName), 0755)

	files.RemoveFile(configFilePath)

	configContent = fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream websocket_khatru29 {
    server 0.0.0.0:5577;
}

# %s
server {
    listen 80;
    listen [::]:80;
    server_name %s;

    location /.well-known/acme-challenge/ {
        root /var/www/%s;
        allow all;
    }

    location / {
        proxy_pass http://websocket_khatru29;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $remote_addr;
    }
}
`, domainName, domainName, domainName)

	files.WriteFile(configFilePath, configContent, 0644)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTP")
}
