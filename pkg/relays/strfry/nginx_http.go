package strfry

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

	const configFile = "strfry.conf"

	err = os.Remove(fmt.Sprintf("/etc/nginx/conf.d/%s", configFile))
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing nginx configuration: %v", err)
	}

	var configContent string

	configContent = fmt.Sprintf(`# %s
server {
    listen 80;
    listen [::]:80;
		server_name %s;

		location /.well-known/acme-challenge/ {
        root /var/www/%s;
        allow all;
    }

    location / {
        proxy_pass http://127.0.0.1:7777;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
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
