package wot_relay

import (
	"fmt"
)

func NginxHttpConfigContent(domainName, wwwDirPath, acmeChallengeDirPath string) string {
	configContent := fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream wot_relay_websocket {
    server localhost:3334;
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
        proxy_pass http://wot_relay_websocket;
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
`, domainName, wwwDirPath, domainName, acmeChallengeDirPath)

	return configContent
}
