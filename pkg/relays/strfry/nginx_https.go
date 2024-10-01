package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to configure nginx for HTTPS
func ConfigureNginxHttps(domainName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring nginx for HTTPS...")

	const configFilePath = "/etc/nginx/conf.d/strfry.conf"

	var configContent string

	files.RemoveFile(configFilePath)

	directories.CreateDirectory(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", domainName), 0755)

	configContent = fmt.Sprintf(`server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name %s;

    root /var/www/%s;

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

    #### SSL Configuration ####
    # Test configuration:
    # https://www.ssllabs.com/ssltest/analyze.html
    # https://cryptcheck.fr/
    ssl_certificate /etc/letsencrypt/live/%s/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/%s/privkey.pem;
    # Verify chain of trust of OCSP response using Root CA and Intermediate certs
    ssl_trusted_certificate /etc/letsencrypt/live/%s/chain.pem;

    # TODO
    # Add support to generate the file in the script
    #ssl_dhparam /etc/ssl/certs/dhparam.pem;

    ssl_protocols TLSv1.2 TLSv1.3;

    # For more information on the security of different cipher suites, you can refer to the following link:
    # https://ciphersuite.info/
    # Compilation of the top cipher suites 2024:
    # https://ssl-config.mozilla.org/#server=nginx
    ssl_ciphers "ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305";

    # Perfect Forward Secrecy (PFS) is frequently compromised without this
    ssl_prefer_server_ciphers on;

    ssl_session_tickets off;

    # Enable SSL session caching for improved performance
    # Try setting ssl_session_timeout to 1d if performance is bad
    ssl_session_timeout 10m;
    ssl_session_cache shared:SSL:10m;

    # By default, the buffer size is 16k, which corresponds to minimal overhead when sending big responses.
    # To minimize Time To First Byte it may be beneficial to use smaller values
    ssl_buffer_size 8k;

    # OCSP stapling
    ssl_stapling on;
    ssl_stapling_verify on;

		#### Security Headers ####
    # Test configuration:
    # https://securityheaders.com/
    # https://observatory.mozilla.org/
    add_header Strict-Transport-Security "max-age=31536000; includeSubdomains; preload";

    add_header X-Frame-Options DENY;

    # Avoid MIME type sniffing
    add_header X-Content-Type-Options nosniff always;

    add_header Referrer-Policy "no-referrer" always;

    add_header X-XSS-Protection 0 always;

    add_header Permissions-Policy "geolocation=(), midi=(), sync-xhr=(), microphone=(), camera=(), magnetometer=(), gyroscope=(), fullscreen=(self), payment=()" always;

    #### Content-Security-Policy (CSP) ####
    add_header Content-Security-Policy "base-uri 'self'; object-src 'none'; frame-ancestors 'none'; upgrade-insecure-requests;" always;
}

server {
    listen 80;
    listen [::]:80;
    server_name %s;

    location /.well-known/acme-challenge/ {
        root /var/www/%s;
        allow all;
    }

    location / {
        return 301 https://%s$request_uri;
    }
}
`, domainName, domainName, domainName, domainName, domainName, domainName, domainName, domainName)

	files.WriteFile(configFilePath, configContent, 0644)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTPS")
}
