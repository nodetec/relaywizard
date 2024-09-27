package khatru_pyramid

import (
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"os"
	"os/exec"
)

// Function to configure nginx for HTTP
func ConfigureNginxHttp(domainName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring nginx for HTTP...")

	err := os.MkdirAll(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", domainName), 0755)
	if err != nil {
		log.Fatalf("Error creating directories: %v", err)
	}

	const configFile = "khatru_pyramid.conf"

	err = os.Remove(fmt.Sprintf("/etc/nginx/conf.d/%s", configFile))
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing nginx configuration: %v", err)
	}

	var configContent string

	configContent = fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream websocket_khatru_pyramid {
    server 0.0.0.0:3334;
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
        proxy_pass http://websocket_khatru_pyramid;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $remote_addr;
    }
}
`, domainName, domainName, domainName)

	err = os.WriteFile(fmt.Sprintf("/etc/nginx/conf.d/%s", configFile), []byte(configContent), 0644)
	if err != nil {
		log.Fatalf("Error writing nginx configuration: %v", err)
	}

	err = exec.Command("systemctl", "restart", "nginx").Run()
	if err != nil {
		log.Fatalf("Error reloading nginx: %v", err)
	}

	spinner.Success("Nginx configured for HTTP")
}
