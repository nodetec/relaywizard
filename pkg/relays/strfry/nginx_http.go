package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to configure Nginx for HTTP
func ConfigureNginxHttp(domainName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx for HTTP...")

	files.RemoveFile(NginxConfigFilePath)

	directories.CreateDirectory(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", domainName), 0755)

	configContent := fmt.Sprintf(`# %s
server {
    listen 80;
    listen [::]:80;
    server_name %s;

    location /.well-known/acme-challenge/ {
        root /var/www/%s;
        allow all;
    }

    location / {
        # First attempt to serve request as file, then
        # as directory, then fall back to displaying 404.
        try_files $uri $uri/ =404;
        proxy_pass http://127.0.0.1:7777;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Only return Nginx in server header
    server_tokens off;

    #### Security Headers ####
    # Test configuration:
    # https://securityheaders.com/
    # https://observatory.mozilla.org/
		add_header X-Frame-Options DENY;

    # Avoid MIME type sniffing
		add_header X-Content-Type-Options nosniff always;

		add_header Referrer-Policy "no-referrer" always;

		add_header X-XSS-Protection 0 always;

		add_header Permissions-Policy "geolocation=(), midi=(), sync-xhr=(), microphone=(), camera=(), magnetometer=(), gyroscope=(), fullscreen=(self), payment=()" always;

    #### Content-Security-Policy (CSP) ####
		add_header Content-Security-Policy "base-uri 'self'; object-src 'none'; frame-ancestors 'none';" always;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name %s;

    location / {
        return 301 http://%s$request_uri;
    }
}
`, domainName, domainName, domainName, domainName, domainName)

	files.WriteFile(NginxConfigFilePath, configContent, 0644)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTP")
}
