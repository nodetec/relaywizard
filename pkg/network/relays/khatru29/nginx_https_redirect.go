package khatru29

import (
	"fmt"
)

func NginxHttpsRedirectConfigContent(domainName, wwwDirPath, acmeChallengeDirPath, certificateDirPath, fullchainFile, privkeyFile, chainFile string) string {
	configContent := fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
  	'' close;
}

upstream khatru29_websocket {
    server 0.0.0.0:5577;
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
        proxy_pass http://khatru29_websocket;
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
				if ($host = %s) {
        	return 301 http://%s$request_uri;
				}
    }

    # Only return Nginx in server header
    server_tokens off;

    #### SSL Configuration ####
    # Test configuration:
    # https://www.ssllabs.com/ssltest/analyze.html
    # https://cryptcheck.fr/
    ssl_certificate %s/%s/%s;
    ssl_certificate_key %s/%s/%s;
    # Verify chain of trust of OCSP response using Root CA and Intermediate certs
    ssl_trusted_certificate %s/%s/%s;
}
`, domainName, wwwDirPath, domainName, acmeChallengeDirPath, domainName, wwwDirPath, domainName, domainName, domainName, certificateDirPath, domainName, fullchainFile, certificateDirPath, domainName, privkeyFile, certificateDirPath, domainName, chainFile)

	return configContent
}
