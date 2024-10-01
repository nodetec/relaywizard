package wot_relay

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

	const configFilePath = "/etc/nginx/conf.d/wot_relay.conf"

	var configContent string

	directories.CreateDirectory(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", domainName), 0755)

	files.RemoveFile(configFilePath)

	configContent = fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
		'' close;
}

upstream websocket_wot_relay {
    server localhost:3334;
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
		    proxy_pass http://websocket_wot_relay;
        proxy_set_header Host $host;
		    proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
		    proxy_set_header Connection $connection_upgrade;
    }
}
`, domainName, domainName, domainName)

	files.WriteFile(configFilePath, configContent, 0644)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTP")
}
