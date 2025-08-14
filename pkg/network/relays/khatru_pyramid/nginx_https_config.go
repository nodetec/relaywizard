package khatru_pyramid

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
)

func NginxHttpsConfigContent(domainName, wwwDirPath, acmeChallengeDirPath, certificateDirPath, fullchainFile, privkeyFile, chainFile string) string {
	configContent := fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream khatru_pyramid_websocket {
    server %s:%s;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name %s;

    root %s/%s;

    location /%s/ {
        default_type "text/plain";
    }

    location / {
        proxy_pass http://khatru_pyramid_websocket;
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

    #### SSL Configuration ####
    # Test configuration:
    # https://www.ssllabs.com/ssltest/analyze.html
    # https://cryptcheck.fr/
    ssl_certificate %s/%s/%s;
    ssl_certificate_key %s/%s/%s;
    # Verify chain of trust of OCSP response using Root CA and Intermediate certs
    ssl_trusted_certificate %s/%s/%s;

    # TODO
    # Add support to generate the file in the script
    #ssl_dhparam /etc/ssl/certs/dhparam.pem;

    ssl_protocols TLSv1.3 TLSv1.2;

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
    add_header X-Content-Type-Options "nosniff" always;

    add_header Referrer-Policy "no-referrer" always;

    add_header X-XSS-Protection "1; mode=block" always;

    add_header Permissions-Policy "geolocation=(), midi=(), sync-xhr=(), microphone=(), camera=(), magnetometer=(), gyroscope=(), fullscreen=(self), payment=()" always;

    #### Content-Security-Policy (CSP) ####
    add_header Content-Security-Policy "base-uri 'self'; object-src 'none'; frame-ancestors 'none'; upgrade-insecure-requests;" always;
}

server {
    listen 80;
    listen [::]:80;
    server_name %s;

    root %s/%s;

    location / {
        return 301 https://%s$request_uri;
    }

    # Only return Nginx in server header
    server_tokens off;
}
`, relays.KhatruPyramidIPv4Address, relays.KhatruPyramidPortNumber, domainName, wwwDirPath, domainName, acmeChallengeDirPath, certificateDirPath, domainName, fullchainFile, certificateDirPath, domainName, privkeyFile, certificateDirPath, domainName, chainFile, domainName, wwwDirPath, domainName, domainName)

	return configContent
}
