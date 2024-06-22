package network

import (
	"fmt"
	"github.com/nodetec/relaywiz/pkg/utils"
	"log"
	"os"
	"os/exec"

	"github.com/pterm/pterm"
)

// Function to configure nginx for HTTPS
func ConfigureNginxHttps(domainName string) {
	dirName := utils.GetDirectoryName(domainName)

	spinner, _ := pterm.DefaultSpinner.Start("Configuring nginx for HTTPS...")
	err := os.Remove("/etc/nginx/conf.d/nostr_relay.conf")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing nginx configuration: %v", err)
	}

	configContent := fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream websocket {
    server 0.0.0.0:3334;
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;
    server_name %s;

    location / {
        proxy_pass http://websocket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $remote_addr;
    }

    #### SSL Configuration ####
    ssl_certificate /etc/letsencrypt/live/%s/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/%s/privkey.pem;

    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_protocols TLSv1.2 TLSv1.1 TLSv1;
    ssl_prefer_server_ciphers on;
    ssl_ciphers "EECDH+ECDSA+AESGCM EECDH+aRSA+AESGCM EECDH+ECDSA+SHA384 EECDH+ECDSA+SHA256 EECDH+aRSA+SHA384 EECDH+aRSA+SHA256 EECDH+aRSA+RC4 EECDH EDH+aRSA RC4 !aNULL !eNULL !LOW !3DES !MD5 !EXP !PSK !SRP !DSS";
    ssl_stapling on;
    ssl_stapling_verify on;
    ssl_ecdh_curve secp384r1;

    add_header Strict-Transport-Security "max-age=31536000; includeSubdomains";
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Referrer-Policy same-origin;
    add_header Feature-Policy "geolocation none;midi none;notifications none;push none;sync-xhr none;microphone none;camera none;magnetometer none;gyroscope none;speaker self;vibrate none;fullscreen self;payment none;";
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
`, domainName, domainName, domainName, domainName, dirName, domainName)

	err = os.WriteFile("/etc/nginx/conf.d/nostr_relay.conf", []byte(configContent), 0644)
	if err != nil {
		log.Fatalf("Error writing nginx configuration: %v", err)
	}

	err = exec.Command("systemctl", "reload", "nginx").Run()
	if err != nil {
		log.Fatalf("Error reloading nginx: %v", err)
	}

	spinner.Success("Nginx configured for HTTPS")
}
